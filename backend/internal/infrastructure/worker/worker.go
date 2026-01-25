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
	"golang.org/x/sync/errgroup"
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

// evaluateSubmission executes the submission against test cases in parallel batches of 8
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

	// 3. Chunk test cases into batches
	batchSize := w.config.Worker.BatchSize
	if batchSize <= 0 {
		batchSize = 4 // Safety fallback
	}
	var batches [][]domain.TestCase
	for i := 0; i < len(testCases); i += batchSize {
		end := i + batchSize
		if end > len(testCases) {
			end = len(testCases)
		}
		batches = append(batches, testCases[i:end])
	}

	// 4. Parallel execution with errgroup
	g, gCtx := errgroup.WithContext(context.Background())
	batchResults := make([]struct {
		TestResults []domain.TestCaseResult
		Verdict     string
		Memory      int
		Runtime     int
		Error       string
		ExitCode    int
		Output      string
	}, len(batches))

	for i, batch := range batches {
		i, batch := i, batch
		g.Go(func() error {
			// Early exit if another batch already failed
			select {
			case <-gCtx.Done():
				return nil
			default:
			}

			// Prepare batch JSON
			type TestInput struct {
				Input    interface{} `json:"input"`
				Expected interface{} `json:"expected"`
			}
			inputs := make([]TestInput, len(batch))
			for j, tc := range batch {
				var inputData interface{}
				_ = json.Unmarshal([]byte(tc.Input), &inputData)
				var expectedData interface{}
				_ = json.Unmarshal([]byte(tc.ExpectedOutput), &expectedData)
				inputs[j] = TestInput{Input: inputData, Expected: expectedData}
			}
			testInputJSON, _ := json.Marshal(inputs)

			// Execute batch on Piston
			fmt.Printf("\n--- Piston Execution Request (Worker Batch %d) ---\n", i)
			fmt.Printf("Stdin (batch size %d): %s\n", len(batch), string(testInputJSON))
			fmt.Printf("Language: %s, Version: %s\n", language.Slug, language.Version)
			fmt.Printf("Memory Limit: %d MB\n", problem.MemoryLimit)
			fmt.Printf("--------------------------------------------------\n\n")

			res, err := w.pistonService.Execute(problem.ID, &submission.ID, language.Slug, language.Version, fullCode, string(testInputJSON))

			if err != nil {
				return fmt.Errorf("batch %d failed: %w", i, err)
			}

			// Store raw output for fallback/compilation check
			batchResults[i].Output = res.Output
			batchResults[i].Error = res.Error
			batchResults[i].ExitCode = res.ExitCode

			// Parse detailed results
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

			if err := json.Unmarshal([]byte(res.Output), &resultObj); err == nil && len(resultObj.TestResults) > 0 {
				batchResults[i].TestResults = make([]domain.TestCaseResult, len(resultObj.TestResults))
				for j, tr := range resultObj.TestResults {
					row := tr.TestCaseResult
					if row.Status == "" && tr.Passed != nil {
						if *tr.Passed {
							row.Status = "passed"
						} else {
							row.Status = "failed"
						}
					}
					if row.ActualOutput == "" && tr.Actual != "" {
						row.ActualOutput = tr.Actual
					}
					batchResults[i].TestResults[j] = row
				}
				batchResults[i].Verdict = resultObj.Verdict
				batchResults[i].Memory = resultObj.Memory
				batchResults[i].Runtime = resultObj.Runtime

				// Short-circuit: if this batch had a failure, signal to stop other batches
				if batchResults[i].Verdict != "ACCEPTED" && batchResults[i].Verdict != "" {
					return fmt.Errorf("short-circuit: batch %d failed with %s", i, batchResults[i].Verdict)
				}
			} else {
				// Compilation or terminal runtime error
				if res.ExitCode != 0 && res.Error != "" && !strings.Contains(res.Output, "verdict") {
					return fmt.Errorf("compilation error in batch %d", i)
				}
				return fmt.Errorf("unexpected output in batch %d", i)
			}

			return nil
		})
	}

	// Wait for all batches or first failure
	if err := g.Wait(); err != nil {
		w.logger.Warn("Batch execution finished with error (short-circuit or Piston error)", zap.Error(err))
	}

	// 5. Aggregate results from all batches
	var finalTestResults []domain.TestCaseResult
	finalStatus := domain.SubmissionStatusAccepted
	passCount := 0
	errorMessage := ""
	maxMemory := 0
	maxRuntime := 0

	for i, res := range batchResults {
		// If a batch didn't run or returned a critical error (compilation)
		if res.ExitCode != 0 && res.Error != "" && !strings.Contains(res.Output, "verdict") {
			w.updateSubmissionError(submission, domain.SubmissionStatusCompilationError, res.Error)
			return
		}

		if len(res.TestResults) == 0 && res.Output != "" {
			// Terminal error in batch execution (crashed before emitting JSON)
			if finalStatus == domain.SubmissionStatusAccepted {
				finalStatus = domain.SubmissionStatusRuntimeError
				errorMessage = res.Error + "\n" + res.Output
			}
		}

		for j := range res.TestResults {
			tr := res.TestResults[j]
			// Restore context from global testCases list
			globalIdx := i*batchSize + j
			if globalIdx < len(testCases) {
				tr.TestID = globalIdx + 1
				tr.IsSample = testCases[globalIdx].IsSample
				tr.Input = testCases[globalIdx].Input
				tr.ExpectedOutput = testCases[globalIdx].ExpectedOutput
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
			finalTestResults = append(finalTestResults, tr)
		}

		if res.Memory > maxMemory {
			maxMemory = res.Memory
		}
		if res.Runtime > maxRuntime {
			maxRuntime = res.Runtime
		}
	}

	// 6. Handle cases where not all batches executed (short-circuit)
	if len(finalTestResults) < len(testCases) && finalStatus == domain.SubmissionStatusAccepted {
		// This should not happen unless g.Wait returned early without setting a failure status
		finalStatus = domain.SubmissionStatusWrongAnswer
		errorMessage = "Execution short-circuited unexpectedly"
	}

	// Metrics fallback
	if maxRuntime == 0 && finalStatus == domain.SubmissionStatusAccepted {
		maxRuntime = 1
	}

	submission.Memory = maxMemory
	submission.Runtime = maxRuntime
	submission.PassedTestCases = passCount
	submission.TestCaseResults = finalTestResults
	submission.ExecutionMetadata, _ = json.Marshal(finalTestResults)

	// 7. Update database and stats
	w.updateSubmissionResult(submission, finalStatus, errorMessage)

	if submission.IsValidationSubmission {
		w.updateValidationStatus(submission, finalStatus, errorMessage, passCount, len(testCases))
	} else {
		w.updateProblemAndUserStats(submission, finalStatus)
	}

	w.logger.Info("Submission processed via parallel batches",
		zap.Int("submission_id", submission.ID),
		zap.Int("batches", len(batches)),
		zap.Int("passed", passCount),
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
	now := time.Now()
	stats := &domain.UserProblemStats{
		UserID:    submission.UserID,
		ProblemID: submission.ProblemID,
		Attempts:  1,
		Status:    "attempted",
		UpdatedAt: now,
	}

	if isAccepted {
		stats.Status = "solved"
		stats.FirstSolvedAt = &now
		stats.BestSubmissionID = &submission.ID
	}

	// Upsert handles both creation and increments/updates
	if err := w.userProblemStatsRepo.Upsert(stats); err != nil {
		w.logger.Error("Failed to upsert user problem stats",
			zap.Error(err),
			zap.Int("user_id", submission.UserID),
			zap.Int("problem_id", submission.ProblemID),
		)
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
