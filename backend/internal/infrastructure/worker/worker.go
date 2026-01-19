package worker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"sync"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/infrastructure/piston"
	"github.com/prabalesh/loco/backend/internal/infrastructure/queue"
	"github.com/prabalesh/loco/backend/pkg/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Worker struct {
	queue                queue.JobQueue
	submissionRepo       domain.SubmissionRepository
	problemRepo          domain.ProblemRepository
	testCaseRepo         domain.TestCaseRepository
	languageRepo         domain.LanguageRepository
	problemLanguageRepo  domain.ProblemLanguageRepository
	pistonService        piston.PistonService
	userProblemStatsRepo domain.UserProblemStatsRepository
	logger               *zap.Logger
	stopChan             chan struct{}
	redisClient          *redis.Client
	workerID             string
	config               *config.Config
}

func NewWorker(
	queue queue.JobQueue,
	submissionRepo domain.SubmissionRepository,
	problemRepo domain.ProblemRepository,
	testCaseRepo domain.TestCaseRepository,
	languageRepo domain.LanguageRepository,
	problemLanguageRepo domain.ProblemLanguageRepository,
	pistonService piston.PistonService,
	userProblemStatsRepo domain.UserProblemStatsRepository,
	logger *zap.Logger,
	redisClient *redis.Client,
	cfg *config.Config,
) *Worker {
	return &Worker{
		queue:                queue,
		submissionRepo:       submissionRepo,
		problemRepo:          problemRepo,
		testCaseRepo:         testCaseRepo,
		languageRepo:         languageRepo,
		problemLanguageRepo:  problemLanguageRepo,
		pistonService:        pistonService,
		userProblemStatsRepo: userProblemStatsRepo,
		logger:               logger,
		stopChan:             make(chan struct{}),
		redisClient:          redisClient,
		workerID:             generateWorkerID(),
		config:               cfg,
	}
}

// Start begins processing jobs from the queue
func (w *Worker) Start(ctx context.Context) {
	w.logger.Info("Worker started, waiting for jobs...",
		zap.Int("max_concurrent_submissions", w.config.Worker.MaxConcurrentSubmissions),
	)
	// Start heartbeat goroutine
	go w.startHeartbeat(ctx)

	// Semaphore to limit concurrent submissions
	sem := make(chan struct{}, w.config.Worker.MaxConcurrentSubmissions)

	for {
		select {
		case <-w.stopChan:
			w.logger.Info("Worker stopped")
			w.stopHeartbeat()
			return
		case <-ctx.Done():
			w.logger.Info("Worker context cancelled")
			w.stopHeartbeat()
			return
		default:
			// Dequeue a job (blocking call)
			job, err := w.queue.DequeueSubmission(ctx)
			if err != nil {
				w.logger.Error("Failed to dequeue job", zap.Error(err))
				time.Sleep(1 * time.Second) // Back off on error
				continue
			}

			if job == nil {
				// No job available, continue to next iteration
				continue
			}

			// Wait for a slot in the semaphore
			sem <- struct{}{}

			// Process the job in a new goroutine
			go func(submissionID int) {
				defer func() { <-sem }()
				w.processSubmission(ctx, submissionID)
			}(job.SubmissionID)
		}
	}
}

// Stop gracefully stops the worker
func (w *Worker) Stop() {
	close(w.stopChan)
}

// heartbeat management
var HeartbeatInterval = 10 * time.Second

func (w *Worker) startHeartbeat(ctx context.Context) {
	ticker := time.NewTicker(HeartbeatInterval)
	defer ticker.Stop()
	key := "worker:" + w.workerID + ":heartbeat"
	for {
		select {
		case <-ticker.C:
			// Set key with TTL twice the interval
			if err := w.redisClient.Set(ctx, key, "alive", 2*HeartbeatInterval).Err(); err != nil {
				w.logger.Error("Failed to set heartbeat", zap.Error(err))
			}
		case <-w.stopChan:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (w *Worker) stopHeartbeat() {
	key := "worker:" + w.workerID + ":heartbeat"
	_ = w.redisClient.Del(context.Background(), key)
}

func generateWorkerID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// processSubmission processes a single submission
func (w *Worker) processSubmission(ctx context.Context, submissionID int) {
	w.logger.Info("Processing submission",
		zap.Int("submission_id", submissionID),
	)

	// Fetch submission
	submission, err := w.submissionRepo.GetByID(submissionID)
	if err != nil {
		w.logger.Error("Failed to fetch submission",
			zap.Error(err),
			zap.Int("submission_id", submissionID),
		)
		return
	}

	// Check if already processed
	if submission.Status != domain.SubmissionStatusPending {
		w.logger.Warn("Submission already processed",
			zap.Int("submission_id", submissionID),
			zap.String("status", string(submission.Status)),
		)
		return
	}

	// Fetch problem and language
	problem, err := w.problemRepo.GetByID(submission.ProblemID)
	if err != nil {
		w.logger.Error("Failed to fetch problem",
			zap.Error(err),
			zap.Int("submission_id", submissionID),
		)
		w.updateSubmissionError(submission, domain.SubmissionStatusInternalError, "Problem not found")
		return
	}

	language, err := w.languageRepo.GetByID(submission.LanguageID)
	if err != nil {
		w.logger.Error("Failed to fetch language",
			zap.Error(err),
			zap.Int("submission_id", submissionID),
		)
		w.updateSubmissionError(submission, domain.SubmissionStatusInternalError, "Language not found")
		return
	}

	// Evaluate submission (always run all test cases for submissions)
	w.evaluateSubmission(submission, problem, language, false)
}

// evaluateSubmission executes the submission against test cases
func (w *Worker) evaluateSubmission(submission *domain.Submission, problem *domain.Problem, language *domain.Language, runOnlyPublicTests bool) {
	var testCases []domain.TestCase
	var err error

	if runOnlyPublicTests {
		testCases, err = w.testCaseRepo.GetSamples(submission.ProblemID)
	} else {
		testCases, err = w.testCaseRepo.GetByProblemID(submission.ProblemID)
	}
	if err != nil {
		w.logger.Error("Failed to fetch test cases",
			zap.Error(err),
			zap.Int("submission_id", submission.ID),
		)
		w.updateSubmissionError(submission, domain.SubmissionStatusInternalError, "Failed to fetch test cases")
		return
	}

	finalCode := submission.Code
	finalStatus := domain.SubmissionStatusAccepted
	errorMessage := ""
	passCount := 0
	totalCount := len(testCases)

	submission.TotalTestCases = totalCount

	// Results channel to collect test case results
	type tcResult struct {
		index  int
		result *piston.ExecutionResult
		err    error
	}
	resultsChan := make(chan tcResult, totalCount)

	// Semaphore to limit concurrent test cases within one submission
	tcSem := make(chan struct{}, w.config.Worker.MaxConcurrentTestCases)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i, tc := range testCases {
		wg.Add(1)
		go func(idx int, testCase domain.TestCase) {
			defer wg.Done()

			// Acquire semaphore slot
			select {
			case tcSem <- struct{}{}:
				defer func() { <-tcSem }()
			case <-ctx.Done():
				return
			}

			res, err := w.pistonService.Execute(language.Slug, language.Version, finalCode, testCase.Input)
			resultsChan <- tcResult{index: idx, result: res, err: err}

			// If any test case fails, we can stop others (optional optimization)
			if err != nil || (res != nil && (res.ExitCode != 0 || strings.TrimSpace(res.Output) != strings.TrimSpace(testCase.ExpectedOutput))) {
				// We don't want to cancel immediately if we want to run all tests,
				// but typically for "Accepted" we need all to pass.
				// For competitive programming, we often stop at first fail.
				// However, if we want full results, we shouldn't cancel.
				// Let's keep it running for all if we want total pass count.
			}
		}(i, tc)
	}

	// Close results channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Process results as they come in
	// Note: since we need total passCount, we must wait for all or handle them as they arrive.
	// To maintain short-circuiting logic but still get counts, it's tricky.
	// Let's collect all results for now to keep it simple and accurate for passCount.
	allResults := make([]tcResult, totalCount)
	for res := range resultsChan {
		allResults[res.index] = res
	}

	// Now analyze results in order to find the first failure (if any)
	for i, r := range allResults {
		if r.err != nil {
			w.logger.Error("Piston execution failed",
				zap.Error(r.err),
				zap.Int("submission_id", submission.ID),
				zap.Int("test_case_index", i),
			)
			if finalStatus == domain.SubmissionStatusAccepted {
				finalStatus = domain.SubmissionStatusInternalError
				errorMessage = "Execution system error"
			}
			continue
		}

		if r.result.ExitCode != 0 {
			if finalStatus == domain.SubmissionStatusAccepted {
				finalStatus = domain.SubmissionStatusRuntimeError
				errorMessage = r.result.Error
			}
			continue
		}

		actual := strings.TrimSpace(r.result.Output)
		expected := strings.TrimSpace(testCases[i].ExpectedOutput)

		if actual != expected {
			if finalStatus == domain.SubmissionStatusAccepted {
				finalStatus = domain.SubmissionStatusWrongAnswer
				errorMessage = fmt.Sprintf("Failed on input: %s\nExpected: %s\nActual: %s", testCases[i].Input, expected, actual)
			}
			continue
		}
		passCount++
	}

	submission.PassedTestCases = passCount

	// Aggregating metrics (using max for runtime and memory)
	maxRuntime := 0
	maxMemory := 0
	for _, r := range allResults {
		if r.result != nil {
			if r.result.Runtime > maxRuntime {
				maxRuntime = r.result.Runtime
			}
			if r.result.Memory > maxMemory {
				maxMemory = r.result.Memory
			}
		}
	}
	submission.Runtime = maxRuntime
	submission.Memory = maxMemory

	// Update submission status
	w.updateSubmissionResult(submission, finalStatus, errorMessage)

	// If it's a validation submission, update ProblemLanguage status
	if submission.IsValidationSubmission {
		now := time.Now()
		pl, err := w.problemLanguageRepo.GetByProblemAndLanguage(submission.ProblemID, submission.LanguageID)
		if err == nil {
			pl.LastValidationStatus = string(finalStatus)
			pl.LastValidationError = errorMessage
			pl.LastPassCount = passCount
			pl.LastTotalCount = totalCount

			if finalStatus == domain.SubmissionStatusAccepted {
				pl.IsValidated = true
				pl.ValidatedAt = &now
			} else {
				pl.IsValidated = false
				pl.ValidatedAt = &now
			}
			w.problemLanguageRepo.Update(pl)
		}
	} else {
		// Update Problem and User stats for regular submissions
		w.updateProblemAndUserStats(submission, finalStatus)
	}

	w.logger.Info("Submission processed",
		zap.Int("submission_id", submission.ID),
		zap.String("status", string(finalStatus)),
		zap.Int("passed", passCount),
		zap.Int("total", totalCount),
	)
}

func (w *Worker) updateProblemAndUserStats(submission *domain.Submission, finalStatus domain.SubmissionStatus) {
	isAccepted := finalStatus == domain.SubmissionStatusAccepted

	// 1. Update Problem Stats (Global)
	if err := w.problemRepo.IncrementStats(submission.ProblemID, isAccepted); err != nil {
		w.logger.Error("Failed to increment problem stats",
			zap.Error(err),
			zap.Int("problem_id", submission.ProblemID),
		)
	}

	// 2. Update User Stats (Per User-Problem)
	stats, err := w.userProblemStatsRepo.Get(submission.UserID, submission.ProblemID)
	if err != nil {
		w.logger.Error("Failed to get user problem stats",
			zap.Error(err),
			zap.Int("user_id", submission.UserID),
			zap.Int("problem_id", submission.ProblemID),
		)
		return
	}

	now := time.Now()
	if stats == nil {
		stats = &domain.UserProblemStats{
			UserID:    submission.UserID,
			ProblemID: submission.ProblemID,
			Attempts:  1,
			Status:    "attempted",
		}
		if isAccepted {
			stats.Status = "solved"
			stats.FirstSolvedAt = &now
			stats.BestSubmissionID = &submission.ID
		}
		if err := w.userProblemStatsRepo.Create(stats); err != nil {
			w.logger.Error("Failed to create user problem stats", zap.Error(err))
		}
	} else {
		stats.Attempts++
		if isAccepted {
			if stats.Status != "solved" {
				stats.Status = "solved"
				stats.FirstSolvedAt = &now
				stats.BestSubmissionID = &submission.ID
			} else {
				// Already solved, check if this is "better"?
				// For now just keep first solved ID or update if better?
				// Let's just keep it simple.
			}
		} else if stats.Status != "solved" {
			stats.Status = "attempted"
		}
		if err := w.userProblemStatsRepo.Update(stats); err != nil {
			w.logger.Error("Failed to update user problem stats", zap.Error(err))
		}
	}

	// 3. Trigger achievement evaluation
	if err := w.queue.EnqueueAchievement(context.Background(), submission.ID); err != nil {
		w.logger.Error("Failed to enqueue achievement job",
			zap.Error(err),
			zap.Int("submission_id", submission.ID),
		)
	}
}

func (w *Worker) updateSubmissionResult(submission *domain.Submission, status domain.SubmissionStatus, errorMsg string) {
	submission.Status = status
	submission.ErrorMessage = errorMsg
	if err := w.submissionRepo.Update(submission); err != nil {
		w.logger.Error("Failed to update submission status",
			zap.Error(err),
			zap.Int("submission_id", submission.ID),
		)
	}
}

func (w *Worker) updateSubmissionError(submission *domain.Submission, status domain.SubmissionStatus, errorMsg string) {
	w.updateSubmissionResult(submission, status, errorMsg)
}
