package execution

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/services/codegen"
	"github.com/prabalesh/loco/backend/internal/services/piston"
	"golang.org/x/sync/errgroup"
)

type ExecutionService struct {
	pistonClient       *piston.PistonClient
	languageMapper     *piston.LanguageMapper
	boilerplateService *codegen.BoilerplateService
	codegenService     *codegen.CodeGenService
	problemRepo        domain.ProblemRepository
}

type ExecutionRequest struct {
	ProblemID  int
	LanguageID int
	UserCode   string
	TestCases  []domain.TestCase
}

type ExecutionResult struct {
	Status       domain.SubmissionStatus `json:"status"` // Accepted, Wrong Answer, etc.
	TestResults  []domain.TestCaseResult `json:"test_results"`
	TotalTests   int                     `json:"total_tests"`
	PassedTests  int                     `json:"passed_tests"`
	Runtime      int                     `json:"runtime"` // milliseconds
	Memory       int                     `json:"memory"`  // kilobytes
	ErrorMessage string                  `json:"error_message,omitempty"`
}

func NewExecutionService(pistonURL string, boilerplateService *codegen.BoilerplateService, codegenService *codegen.CodeGenService, problemRepo domain.ProblemRepository, executionRepo domain.PistonExecutionRepository) *ExecutionService {
	return &ExecutionService{
		pistonClient:       piston.NewPistonClient(pistonURL, executionRepo),
		languageMapper:     piston.NewLanguageMapper(),
		boilerplateService: boilerplateService,
		codegenService:     codegenService,
		problemRepo:        problemRepo,
	}
}

// ExecuteSubmission is a wrapper for backward compatibility
func (s *ExecutionService) ExecuteSubmission(req ExecutionRequest, languageSlug string) (*ExecutionResult, error) {
	return s.ExecuteBatchSubmission(context.Background(), req, languageSlug)
}

// ExecuteBatchSubmission executes user code against test cases in parallel batches
func (s *ExecutionService) ExecuteBatchSubmission(ctx context.Context, req ExecutionRequest, languageSlug string) (*ExecutionResult, error) {
	// Validate
	if req.UserCode == "" {
		return nil, errors.New("user code is required")
	}
	if len(req.TestCases) == 0 {
		return nil, errors.New("at least one test case is required")
	}

	// 1. Get Problem Info
	problem, err := s.problemRepo.GetByID(req.ProblemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get problem: %w", err)
	}

	var params []domain.SchemaParameter
	if problem.Parameters != nil {
		json.Unmarshal(*problem.Parameters, &params)
	}

	schema := domain.ProblemSchema{
		FunctionName: func() string {
			if problem.FunctionName != nil {
				return *problem.FunctionName
			}
			return ""
		}(),
		ReturnType: domain.GenericType(func() string {
			if problem.ReturnType != nil {
				return *problem.ReturnType
			}
			return "int"
		}()),
		Parameters: params,
	}

	// 2. Generate Universal Harness (same code for all batches)
	fullCode, err := s.codegenService.GenerateTestHarness(schema, req.UserCode, languageSlug, []domain.TestCase{}, problem.ValidationType)
	if err != nil {
		return nil, fmt.Errorf("failed to generate harness: %w", err)
	}

	runtime, err := s.languageMapper.GetPistonRuntime(languageSlug)
	if err != nil {
		return nil, err
	}

	// 3. Batching
	batchSize := 8
	var batches [][]domain.TestCase
	for i := 0; i < len(req.TestCases); i += batchSize {
		end := i + batchSize
		if end > len(req.TestCases) {
			end = len(req.TestCases)
		}
		batches = append(batches, req.TestCases[i:end])
	}

	// 4. Parallel Execution with ErrGroup
	g, gCtx := errgroup.WithContext(ctx)
	results := make([]*ExecutionResult, len(batches))

	for i, batch := range batches {
		i, batch := i, batch
		g.Go(func() error {
			// Check if cancelled
			select {
			case <-gCtx.Done():
				return gCtx.Err()
			default:
			}

			// Prepare STDIN for batch
			stdinBytes, _ := json.Marshal(batch)
			stdin := string(stdinBytes)

			pistonReq := piston.ExecuteRequest{
				Language: runtime.Language,
				Version:  runtime.Version,
				Files: []piston.File{
					{
						Name:    runtime.FileName,
						Content: fullCode,
					},
				},
				Stdin:          stdin,
				CompileTimeout: 10000,
				RunTimeout:     5000,
				ProblemID:      req.ProblemID,
				SubmissionID:   nil, // Test runs don't have submission ID here
			}

			fmt.Printf("\n--- Piston Execution Request ---\n")
			fmt.Printf("Stdin (batch size %d): %s\n", len(batch), stdin)
			fmt.Printf("Time Limit (RunTimeout): %d ms\n", pistonReq.RunTimeout)
			fmt.Printf("Compile Timeout: %d ms\n", pistonReq.CompileTimeout)
			fmt.Printf("Problem Memory Limit: %d MB\n", problem.MemoryLimit)
			fmt.Printf("-------------------------------\n\n")

			pistonResp, err := s.pistonClient.Execute(pistonReq)
			if err != nil {
				return err
			}

			// Check for compilation errors
			if pistonResp.Compile != nil && pistonResp.Compile.Code != 0 {
				results[i] = &ExecutionResult{
					Status:       domain.SubmissionStatusCompilationError,
					ErrorMessage: pistonResp.Compile.Stderr,
				}
				return fmt.Errorf("compilation error") // Stop other batches
			}

			// Parse results for this batch
			batchRes, err := s.validateOutput(pistonResp.Run.Stdout, batch)
			if err != nil {
				return err
			}

			results[i] = batchRes

			// Short circuit WA/TLE/RE in batches
			if batchRes.Status != domain.SubmissionStatusAccepted {
				return fmt.Errorf("test failure") // Stops other batches via gCtx
			}

			return nil
		})
	}

	// Wait for all batches (or first error)
	_ = g.Wait()

	// 5. Aggregate Results
	finalResult := &ExecutionResult{
		Status:      domain.SubmissionStatusAccepted,
		TotalTests:  len(req.TestCases),
		PassedTests: 0,
		Runtime:     0,
		Memory:      0,
		TestResults: make([]domain.TestCaseResult, 0, len(req.TestCases)),
	}

	for _, res := range results {
		if res == nil {
			continue // Batch might not have started or was cancelled
		}
		finalResult.PassedTests += res.PassedTests
		finalResult.TestResults = append(finalResult.TestResults, res.TestResults...)
		if res.Runtime > finalResult.Runtime {
			finalResult.Runtime = res.Runtime
		}
		if res.Memory > finalResult.Memory {
			finalResult.Memory = res.Memory
		}

		if res.Status != domain.SubmissionStatusAccepted && finalResult.Status == domain.SubmissionStatusAccepted {
			finalResult.Status = res.Status
			finalResult.ErrorMessage = res.ErrorMessage
		}
	}

	// Re-check status if we stopped early
	if finalResult.PassedTests < finalResult.TotalTests && finalResult.Status == domain.SubmissionStatusAccepted {
		// If we haven't failed but haven't passed all, something went wrong (cancelled)
		// But in our case, if one failed, we'll have a non-accepted status.
	}

	return finalResult, nil
}

type harnessVerdict struct {
	Verdict     string `json:"verdict"`
	Runtime     int    `json:"runtime"`
	Memory      int    `json:"memory"`
	TestResults []struct {
		Passed bool   `json:"passed"`
		Input  string `json:"input"`
		Actual string `json:"actual"`
		Error  string `json:"error"`
	} `json:"test_results"`
}

// validateOutput parses stdout from the harness and aggregates results
func (s *ExecutionService) validateOutput(stdout string, testCases []domain.TestCase) (*ExecutionResult, error) {
	var verdict harnessVerdict
	if err := json.Unmarshal([]byte(stdout), &verdict); err != nil {
		return &ExecutionResult{
			Status:       domain.SubmissionStatusRuntimeError,
			ErrorMessage: fmt.Sprintf("Failed to parse output: %v\nRaw output: %s", err, stdout),
			TestResults:  []domain.TestCaseResult{},
		}, nil
	}

	overallStatus := domain.SubmissionStatusAccepted
	switch verdict.Verdict {
	case "ACCEPTED":
		overallStatus = domain.SubmissionStatusAccepted
	case "WRONG_ANSWER":
		overallStatus = domain.SubmissionStatusWrongAnswer
	case "TLE":
		overallStatus = domain.SubmissionStatusTimeLimitExceeded
	case "RUNTIME_ERROR":
		overallStatus = domain.SubmissionStatusRuntimeError
	default:
		overallStatus = domain.SubmissionStatusRuntimeError
	}

	results := make([]domain.TestCaseResult, len(verdict.TestResults))
	passedCount := 0
	for i, tr := range verdict.TestResults {
		status := "Passed"
		if !tr.Passed {
			status = "Failed"
		} else {
			passedCount++
		}

		results[i] = domain.TestCaseResult{
			TestID:       i + 1,
			Status:       status,
			Input:        tr.Input,
			ActualOutput: tr.Actual,
			Error:        tr.Error,
		}

		if i < len(testCases) {
			results[i].IsSample = testCases[i].IsSample
			results[i].ExpectedOutput = testCases[i].ExpectedOutput
		}
	}

	return &ExecutionResult{
		Status:      overallStatus,
		TestResults: results,
		TotalTests:  len(results),
		PassedTests: passedCount,
		Runtime:     verdict.Runtime,
		Memory:      verdict.Memory,
	}, nil
}
