package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/infrastructure/piston"
	"github.com/prabalesh/loco/backend/internal/infrastructure/queue"
	"github.com/prabalesh/loco/backend/internal/services/codegen"
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
	boilerplateService   *codegen.BoilerplateService
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
	boilerplateService *codegen.BoilerplateService,
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
		boilerplateService:   boilerplateService,
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
		w.logger.Error("Failed to fetch test cases", zap.Error(err), zap.Int("submission_id", submission.ID))
		w.updateSubmissionError(submission, domain.SubmissionStatusInternalError, "Failed to fetch test cases")
		return
	}

	submission.TotalTestCases = len(testCases)

	// 1. Get test harness template
	harnessTemplate, err := w.boilerplateService.GetTestHarnessTemplate(submission.ProblemID, submission.LanguageID)
	if err != nil {
		w.logger.Error("Failed to get harness template", zap.Error(err), zap.Int("submission_id", submission.ID))
		w.updateSubmissionError(submission, domain.SubmissionStatusInternalError, "Harness template not found")
		return
	}

	// 2. Inject user code
	fullCode := w.boilerplateService.InjectUserCodeIntoHarness(harnessTemplate, submission.Code)

	// 3. Prepare test cases input (JSON)
	type TestInput struct {
		Input    interface{} `json:"input"`
		Expected interface{} `json:"expected"`
	}
	inputs := make([]TestInput, len(testCases))
	for i, tc := range testCases {
		var inputData interface{}
		_ = json.Unmarshal([]byte(tc.Input), &inputData)
		var expectedData interface{}
		_ = json.Unmarshal([]byte(tc.ExpectedOutput), &expectedData)
		inputs[i] = TestInput{Input: inputData, Expected: expectedData}
	}
	testInputJSON, _ := json.Marshal(inputs)

	// 4. Execute on Piston
	res, err := w.pistonService.Execute(language.Slug, language.Version, fullCode, string(testInputJSON))
	if err != nil {
		w.logger.Error("Piston execution failed", zap.Error(err), zap.Int("submission_id", submission.ID))
		w.updateSubmissionError(submission, domain.SubmissionStatusInternalError, "Execution system error")
		return
	}

	// 5. Check for compilation error
	if res.ExitCode != 0 && res.Error != "" && !strings.Contains(res.Output, "[{\"test_id\"") {
		w.updateSubmissionError(submission, domain.SubmissionStatusCompilationError, res.Error)
		return
	}

	// 6. Parse detailed results
	var testResults []domain.TestCaseResult
	if err := json.Unmarshal([]byte(res.Output), &testResults); err != nil {
		// If parsing fails, it might be a runtime error of the harness itself
		w.updateSubmissionError(submission, domain.SubmissionStatusRuntimeError, res.Error+"\n"+res.Output)
		return
	}

	// 7. Process results and determine final status
	finalStatus := domain.SubmissionStatusAccepted
	passCount := 0
	maxTime := 0
	maxMemory := 0
	errorMessage := ""

	for i := range testResults {
		tr := &testResults[i]
		if i < len(testCases) {
			tr.IsSample = testCases[i].IsSample
			// If not provided by harness, fill from DB for completeness
			if tr.Input == "" {
				tr.Input = testCases[i].Input
			}
			if tr.ExpectedOutput == "" {
				tr.ExpectedOutput = testCases[i].ExpectedOutput
			}
		}

		if tr.Status == "passed" {
			passCount++
		} else if finalStatus == domain.SubmissionStatusAccepted {
			switch tr.Status {
			case "timeout":
				finalStatus = domain.SubmissionStatusTimeLimitExceeded
			case "runtime_error":
				finalStatus = domain.SubmissionStatusRuntimeError
				errorMessage = tr.Error
			default:
				finalStatus = domain.SubmissionStatusWrongAnswer
				errorMessage = fmt.Sprintf("Failed on test %d", tr.TestID)
			}
		}

		if tr.TimeMS > maxTime {
			maxTime = tr.TimeMS
		}
		if tr.MemoryKB > maxMemory {
			maxMemory = tr.MemoryKB
		}
	}

	submission.PassedTestCases = passCount
	submission.Runtime = maxTime
	submission.Memory = maxMemory
	submission.TestCaseResults = testResults
	submission.ExecutionMetadata, _ = json.Marshal(testResults)

	// Update submission status
	w.updateSubmissionResult(submission, finalStatus, errorMessage)

	// Update ProblemLanguage status if validation
	if submission.IsValidationSubmission {
		w.updateValidationStatus(submission, finalStatus, errorMessage, passCount, len(testCases))
	} else {
		w.updateProblemAndUserStats(submission, finalStatus)
	}

	w.logger.Info("Submission processed",
		zap.Int("submission_id", submission.ID),
		zap.String("status", string(finalStatus)),
		zap.Int("passed", passCount),
		zap.Int("total", len(testCases)),
	)
}

func (w *Worker) updateValidationStatus(submission *domain.Submission, status domain.SubmissionStatus, errorMsg string, passCount, totalCount int) {
	now := time.Now()
	pl, err := w.problemLanguageRepo.GetByProblemAndLanguage(submission.ProblemID, submission.LanguageID)
	if err != nil {
		return
	}
	pl.LastValidationStatus = string(status)
	pl.LastValidationError = errorMsg
	pl.LastPassCount = passCount
	pl.LastTotalCount = totalCount
	pl.IsValidated = (status == domain.SubmissionStatusAccepted)
	pl.ValidatedAt = &now
	w.problemLanguageRepo.Update(pl)
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
