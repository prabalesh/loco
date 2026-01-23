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

func NewExecutionService(pistonURL string, boilerplateService *codegen.BoilerplateService, codegenService *codegen.CodeGenService, problemRepo domain.ProblemRepository) *ExecutionService {
	return &ExecutionService{
		pistonClient:       piston.NewPistonClient(pistonURL),
		languageMapper:     piston.NewLanguageMapper(),
		boilerplateService: boilerplateService,
		codegenService:     codegenService,
		problemRepo:        problemRepo,
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

	// Get problem schema for harness generation
	problem, err := s.problemRepo.GetByID(req.ProblemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get problem: %w", err)
	}

	// Parse parameters from JSON
	var params []domain.SchemaParameter
	if problem.Parameters != nil {
		if err := json.Unmarshal(*problem.Parameters, &params); err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
	}

	// Build problem schema
	functionName := ""
	if problem.FunctionName != nil {
		functionName = *problem.FunctionName
	}
	returnType := domain.TypeInteger
	if problem.ReturnType != nil {
		returnType = domain.GenericType(*problem.ReturnType)
	}

	schema := domain.ProblemSchema{
		FunctionName: functionName,
		ReturnType:   returnType,
		Parameters:   params,
	}

	// Generate fresh harness with user code (this applies all our fixes)
	fullCode, err := s.codegenService.GenerateTestHarness(schema, req.UserCode, languageSlug, req.TestCases, problem.ValidationType)
	if err != nil {
		return nil, fmt.Errorf("failed to generate harness: %w", err)
	}

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
			Status:       domain.SubmissionStatusInternalError,
			ErrorMessage: fmt.Sprintf("Failed to parse output: %v\nRaw output: %s", err, stdout),
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
