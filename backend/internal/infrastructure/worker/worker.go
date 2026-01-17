package worker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/infrastructure/piston"
	"github.com/prabalesh/loco/backend/internal/infrastructure/queue"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Worker struct {
	queue               queue.JobQueue
	submissionRepo      domain.SubmissionRepository
	problemRepo         domain.ProblemRepository
	testCaseRepo        domain.TestCaseRepository
	languageRepo        domain.LanguageRepository
	problemLanguageRepo domain.ProblemLanguageRepository
	pistonService       piston.PistonService
	logger              *zap.Logger
	stopChan            chan struct{}
	redisClient         *redis.Client
	workerID            string
}

func NewWorker(
	queue queue.JobQueue,
	submissionRepo domain.SubmissionRepository,
	problemRepo domain.ProblemRepository,
	testCaseRepo domain.TestCaseRepository,
	languageRepo domain.LanguageRepository,
	problemLanguageRepo domain.ProblemLanguageRepository,
	pistonService piston.PistonService,
	logger *zap.Logger,
	redisClient *redis.Client,
) *Worker {
	return &Worker{
		queue:               queue,
		submissionRepo:      submissionRepo,
		problemRepo:         problemRepo,
		testCaseRepo:        testCaseRepo,
		languageRepo:        languageRepo,
		problemLanguageRepo: problemLanguageRepo,
		pistonService:       pistonService,
		logger:              logger,
		stopChan:            make(chan struct{}),
		redisClient:         redisClient,
		workerID:            generateWorkerID(),
	}
}

// Start begins processing jobs from the queue
func (w *Worker) Start(ctx context.Context) {
	w.logger.Info("Worker started, waiting for jobs...")
	// Start heartbeat goroutine
	go w.startHeartbeat(ctx)

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

			// Process the job
			w.processSubmission(ctx, job.SubmissionID)
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

	for _, tc := range testCases {
		result, err := w.pistonService.Execute(language.Slug, language.Version, finalCode, tc.Input)
		if err != nil {
			w.logger.Error("Piston execution failed",
				zap.Error(err),
				zap.Int("submission_id", submission.ID),
			)
			finalStatus = domain.SubmissionStatusInternalError
			errorMessage = "Execution system error"
			break
		}

		if result.ExitCode != 0 {
			finalStatus = domain.SubmissionStatusRuntimeError
			errorMessage = result.Error
			break
		}

		// Normalize output (trim whitespace)
		actual := strings.TrimSpace(result.Output)
		expected := strings.TrimSpace(tc.ExpectedOutput)

		if actual != expected {
			finalStatus = domain.SubmissionStatusWrongAnswer
			errorMessage = fmt.Sprintf("Failed on input: %s\nExpected: %s\nActual: %s", tc.Input, expected, actual)
			break
		}
		passCount++
	}

	submission.PassedTestCases = passCount

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
	}

	w.logger.Info("Submission processed",
		zap.Int("submission_id", submission.ID),
		zap.String("status", string(finalStatus)),
		zap.Int("passed", passCount),
		zap.Int("total", totalCount),
	)
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
