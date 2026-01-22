package execution

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/services/codegen"
	"github.com/prabalesh/loco/backend/internal/services/piston"
)

type ExecutionService struct {
	pistonClient       *piston.PistonClient
	languageMapper     *piston.LanguageMapper
	boilerplateService *codegen.BoilerplateService
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

func NewExecutionService(pistonURL string, boilerplateService *codegen.BoilerplateService) *ExecutionService {
	return &ExecutionService{
		pistonClient:       piston.NewPistonClient(pistonURL),
		languageMapper:     piston.NewLanguageMapper(),
		boilerplateService: boilerplateService,
	}
}

// ExecuteSubmission executes user code against test cases
func (s *ExecutionService) ExecuteSubmission(req ExecutionRequest, languageSlug string) (*ExecutionResult, error) {
	// Validate
	if req.UserCode == "" {
		return nil, errors.New("user code is required")
	}
	if len(req.TestCases) == 0 {
		return nil, errors.New("at least one test case is required")
	}

	// Get test harness template
	harnessTemplate, err := s.boilerplateService.GetTestHarnessTemplate(req.ProblemID, req.LanguageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get harness template: %w", err)
	}

	// Inject user code into harness
	fullCode := s.boilerplateService.InjectUserCodeIntoHarness(harnessTemplate, req.UserCode)

	// Get Piston runtime info
	runtime, err := s.languageMapper.GetPistonRuntime(languageSlug)
	if err != nil {
		return nil, err
	}

	fmt.Println("Runtime: ", runtime)
	fmt.Println("Full Code: ", fullCode)
	fmt.Println("Test Input: ", req.TestCases)

	// Execute on Piston
	pistonReq := piston.ExecuteRequest{
		Language: runtime.Language,
		Version:  runtime.Version,
		Files: []piston.File{
			{
				Name:    runtime.FileName,
				Content: fullCode,
			},
		},
		Stdin:          "", // Test cases are embedded in the harness
		Args:           []string{},
		CompileTimeout: 10000,
		RunTimeout:     5000, // 5 seconds max
	}

	pistonResp, err := s.pistonClient.Execute(pistonReq)
	if err != nil {
		return nil, fmt.Errorf("execution failed: %w", err)
	}

	fmt.Println("Piston Stdout: ", pistonResp.Run.Stdout)

	// Check for compilation errors
	if pistonResp.Compile != nil && pistonResp.Compile.Code != 0 {
		return &ExecutionResult{
			Status:       domain.SubmissionStatusCompilationError,
			ErrorMessage: pistonResp.Compile.Stderr,
		}, nil
	}

	// Check for runtime errors
	if pistonResp.Run.Code != 0 {
		return &ExecutionResult{
			Status:       domain.SubmissionStatusRuntimeError,
			ErrorMessage: pistonResp.Run.Stderr,
		}, nil
	}

	// Parse output and validate
	return s.validateOutput(pistonResp.Run.Stdout, req.TestCases)
}

// validateOutput parses stdout from the harness and aggregates results
func (s *ExecutionService) validateOutput(stdout string, testCases []domain.TestCase) (*ExecutionResult, error) {
	var results []domain.TestCaseResult
	if err := json.Unmarshal([]byte(stdout), &results); err != nil {
		return &ExecutionResult{
			Status:       domain.SubmissionStatusInternalError,
			ErrorMessage: fmt.Sprintf("Failed to parse output: %v\nRaw output: %s", err, stdout),
		}, nil
	}

	passedCount := 0
	maxTime := 0
	maxMemory := 0
	overallStatus := domain.SubmissionStatusAccepted

	for i := range results {
		res := &results[i]
		if i < len(testCases) {
			res.IsSample = testCases[i].IsSample
			if res.Input == "" {
				res.Input = testCases[i].Input
			}
			if res.ExpectedOutput == "" {
				res.ExpectedOutput = testCases[i].ExpectedOutput
			}
		}

		if res.Status == "passed" {
			passedCount++
			res.Status = "Passed"
		} else if res.Status == "failed" {
			res.Status = "Wrong Answer"
			if overallStatus == domain.SubmissionStatusAccepted {
				overallStatus = domain.SubmissionStatusWrongAnswer
			}
		} else if res.Status == "timeout" {
			res.Status = "Time Limit Exceeded"
			if overallStatus == domain.SubmissionStatusAccepted || overallStatus == domain.SubmissionStatusWrongAnswer {
				overallStatus = domain.SubmissionStatusTimeLimitExceeded
			}
		} else if res.Status == "runtime_error" {
			res.Status = "Runtime Error"
			if overallStatus == domain.SubmissionStatusAccepted || overallStatus == domain.SubmissionStatusWrongAnswer || overallStatus == domain.SubmissionStatusTimeLimitExceeded {
				overallStatus = domain.SubmissionStatusRuntimeError
			}
		}

		if res.TimeMS > maxTime {
			maxTime = res.TimeMS
		}
		if res.MemoryKB > maxMemory {
			maxMemory = res.MemoryKB
		}
	}

	return &ExecutionResult{
		Status:      overallStatus,
		TestResults: results,
		TotalTests:  len(results),
		PassedTests: passedCount,
		Runtime:     maxTime,
		Memory:      maxMemory,
	}, nil
}
