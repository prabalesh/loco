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
	queue                 queue.JobQueue
	submissionRepo        domain.SubmissionRepository
	problemRepo           domain.ProblemRepository
	testCaseRepo          domain.TestCaseRepository
	languageRepo          domain.LanguageRepository
	problemLanguageRepo   domain.ProblemLanguageRepository
	referenceSolutionRepo domain.ReferenceSolutionRepository
	pistonService         piston.PistonService
	boilerplateService    *codegen.BoilerplateService
	userProblemStatsRepo  domain.UserProblemStatsRepository
	logger                *zap.Logger
	stopChan              chan struct{}
	redisClient           *redis.Client
	workerID              string
	config                *config.Config
}

func NewWorker(
	queue queue.JobQueue,
	submissionRepo domain.SubmissionRepository,
	problemRepo domain.ProblemRepository,
	testCaseRepo domain.TestCaseRepository,
	languageRepo domain.LanguageRepository,
	problemLanguageRepo domain.ProblemLanguageRepository,
	referenceSolutionRepo domain.ReferenceSolutionRepository,
	pistonService piston.PistonService,
	boilerplateService *codegen.BoilerplateService,
	userProblemStatsRepo domain.UserProblemStatsRepository,
	logger *zap.Logger,
	redisClient *redis.Client,
	cfg *config.Config,
) *Worker {
	return &Worker{
		queue:                 queue,
		submissionRepo:        submissionRepo,
		problemRepo:           problemRepo,
		testCaseRepo:          testCaseRepo,
		languageRepo:          languageRepo,
		problemLanguageRepo:   problemLanguageRepo,
		referenceSolutionRepo: referenceSolutionRepo,
		pistonService:         pistonService,
		boilerplateService:    boilerplateService,
		userProblemStatsRepo:  userProblemStatsRepo,
		logger:                logger,
		stopChan:              make(chan struct{}),
		redisClient:           redisClient,
		workerID:              generateWorkerID(),
		config:                cfg,
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

	if runOnlyPublicTests || submission.IsRunOnly {
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
	if res.ExitCode != 0 && res.Error != "" && !strings.Contains(res.Output, "verdict") {
		w.updateSubmissionError(submission, domain.SubmissionStatusCompilationError, res.Error)
		return
	}

	// 6. Parse detailed results
	var testResults []domain.TestCaseResult

	type HarnessTestResult struct {
		domain.TestCaseResult
		Passed *bool  `json:"passed"`
		Actual string `json:"actual"`
	}

	var resultObj struct {
		TestResults []HarnessTestResult `json:"test_results"`
		Verdict     string              `json:"verdict"`
		Memory      int                 `json:"memory"`
		Runtime     int                 `json:"runtime"`
	}

	// Try unmarshaling the JSON output from the C++ Harness
	if err := json.Unmarshal([]byte(res.Output), &resultObj); err == nil && len(resultObj.TestResults) > 0 {
		testResults = make([]domain.TestCaseResult, len(resultObj.TestResults))
		for i, tr := range resultObj.TestResults {
			row := tr.TestCaseResult
			// Map 'passed' bool to status string
			if row.Status == "" && tr.Passed != nil {
				if *tr.Passed {
					row.Status = "passed"
				} else {
					row.Status = "failed"
				}
			}
			// Map 'actual' to 'actual_output'
			if row.ActualOutput == "" && tr.Actual != "" {
				row.ActualOutput = tr.Actual
			}
			testResults[i] = row
		}
	} else {
		// Fallback for real runtime errors or crashes
		w.updateSubmissionError(submission, domain.SubmissionStatusRuntimeError, res.Error+"\n"+res.Output)
		return
	}

	// 7. Process results and determine final status
	finalStatus := domain.SubmissionStatusAccepted
	passCount := 0
	errorMessage := ""

	// FIX: Assign Aggregate Metrics from Harness directly
	submission.Memory = resultObj.Memory
	submission.Runtime = resultObj.Runtime

	for i := range testResults {
		tr := &testResults[i]
		if tr.TestID == 0 {
			tr.TestID = i + 1
		}

		if i < len(testCases) {
			tr.IsSample = testCases[i].IsSample
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
	}

	// Final Fallback and Safety Checks
	if submission.Memory == 0 && res.Memory > 0 {
		submission.Memory = res.Memory
	}
	// If the algorithm ran so fast it returned 0ms, force 1ms for better UI experience
	if submission.Runtime == 0 && finalStatus == domain.SubmissionStatusAccepted {
		submission.Runtime = 1
	}

	submission.PassedTestCases = passCount
	submission.TestCaseResults = testResults
	submission.ExecutionMetadata, _ = json.Marshal(testResults)

	// 8. Update database record
	w.updateSubmissionResult(submission, finalStatus, errorMessage)

	// 9. Update stats/validation status
	if submission.IsValidationSubmission {
		w.updateValidationStatus(submission, finalStatus, errorMessage, passCount, len(testCases))
	} else {
		w.updateProblemAndUserStats(submission, finalStatus)
	}

	w.logger.Info("Submission processed successfully",
		zap.Int("submission_id", submission.ID),
		zap.Int("memory", submission.Memory),
		zap.Int("runtime", submission.Runtime),
		zap.String("status", string(finalStatus)),
	)
}

func (w *Worker) updateValidationStatus(submission *domain.Submission, status domain.SubmissionStatus, errorMsg string, passCount, totalCount int) {
	// Update ReferenceSolution
	refSol, err := w.referenceSolutionRepo.GetByProblemAndLanguage(submission.ProblemID, submission.LanguageID)
	if err == nil {
		refSol.IsValidated = (status == domain.SubmissionStatusAccepted)

		validationResult := map[string]interface{}{
			"is_valid":      refSol.IsValidated,
			"passed_tests":  passCount,
			"total_tests":   totalCount,
			"error_message": errorMsg,
			"test_results":  submission.TestCaseResults,
		}

		resultJSON, _ := json.Marshal(validationResult)
		refSol.ValidationResults = resultJSON

		if err := w.referenceSolutionRepo.Update(refSol); err != nil {
			w.logger.Error("Failed to update reference solution", zap.Error(err))
		}
	} else {
		w.logger.Error("Reference solution not found for validation update", zap.Error(err))
	}

	// Update Problem Validation Status if valid
	if status == domain.SubmissionStatusAccepted {
		problem, err := w.problemRepo.GetByID(submission.ProblemID)
		if err == nil {
			problem.ValidationStatus = "validated"
			problem.HasReferenceSolution = true
			w.problemRepo.Update(problem)
		}
	}

	// Legacy support: Update ProblemLanguage as well (remvoed as per user request to stop using it)
	// The table might still exist but we don't want to fail or log errors if records are missing,
	// and we are moving away from it.

}

func (w *Worker) updateProblemAndUserStats(submission *domain.Submission, finalStatus domain.SubmissionStatus) {
	if submission.IsRunOnly {
		return
	}

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

	// Debug: Print submission data before saving
	fmt.Println("========================================")
	fmt.Println("=== SUBMISSION RESULT ===")
	fmt.Printf("Submission ID: %d\n", submission.ID)
	fmt.Printf("User ID: %d\n", submission.UserID)
	fmt.Printf("Problem ID: %d\n", submission.ProblemID)
	fmt.Printf("Status: %s\n", status)
	fmt.Printf("Error Message: %s\n", errorMsg)
	fmt.Printf("Passed Test Cases: %d\n", submission.PassedTestCases)
	fmt.Printf("Total Test Cases: %d\n", submission.TotalTestCases)
	fmt.Printf("Runtime (ms): %d\n", submission.Runtime)
	fmt.Printf("Memory (kb): %d\n", submission.Memory)
	fmt.Printf("Number of Results: %d\n", len(submission.TestCaseResults))
	fmt.Println("========================================")

	// Print each test case result
	for i, tcResult := range submission.TestCaseResults {
		fmt.Printf("Test Case %d:\n", i)
		fmt.Printf("  Status: %s\n", tcResult.Status)
		fmt.Printf("  Time (ms): %d\n", tcResult.TimeMS)
		fmt.Printf("  Memory (kb): %d\n", tcResult.MemoryKB)
		fmt.Printf("  Is Sample: %v\n", tcResult.IsSample)
		fmt.Println("---")
	}
	fmt.Println("========================================")

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
