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

	// Prepare test cases input (JSON format)
	testInput, err := s.prepareTestInput(req.TestCases)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare test input: %w", err)
	}

	// Get Piston runtime info
	runtime, err := s.languageMapper.GetPistonRuntime(languageSlug)
	if err != nil {
		return nil, err
	}

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
		Stdin:          testInput,
		Args:           []string{},
		CompileTimeout: 10000,
		RunTimeout:     5000, // 5 seconds max
	}

	pistonResp, err := s.pistonClient.Execute(pistonReq)
	if err != nil {
		return nil, fmt.Errorf("execution failed: %w", err)
	}

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

// prepareTestInput converts test cases to JSON for stdin
func (s *ExecutionService) prepareTestInput(testCases []domain.TestCase) (string, error) {
	type TestInput struct {
		Input interface{} `json:"input"`
	}

	inputs := []TestInput{}
	for _, tc := range testCases {
		var input interface{}
		if err := json.Unmarshal([]byte(tc.Input), &input); err != nil {
			return "", fmt.Errorf("invalid test case input: %w", err)
		}
		inputs = append(inputs, TestInput{Input: input})
	}

	data, err := json.Marshal(inputs)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// validateOutput parses stdout and compares with expected outputs
func (s *ExecutionService) validateOutput(stdout string, testCases []domain.TestCase) (*ExecutionResult, error) {
	// Parse stdout as JSON array
	type OutputResult struct {
		Output interface{} `json:"output"`
		Error  *string     `json:"error"`
	}

	var outputs []OutputResult
	if err := json.Unmarshal([]byte(stdout), &outputs); err != nil {
		return &ExecutionResult{
			Status:       domain.SubmissionStatusInternalError,
			ErrorMessage: fmt.Sprintf("Failed to parse output: %v\nRaw output: %s", err, stdout),
		}, nil
	}

	// Validate count matches
	if len(outputs) != len(testCases) {
		return &ExecutionResult{
			Status:       domain.SubmissionStatusInternalError,
			ErrorMessage: fmt.Sprintf("Output count mismatch: got %d, expected %d", len(outputs), len(testCases)),
		}, nil
	}

	// Compare each test case
	testResults := []domain.TestCaseResult{}
	passedCount := 0

	for i, tc := range testCases {
		output := outputs[i]

		// Check if test had runtime error
		if output.Error != nil {
			testResults = append(testResults, domain.TestCaseResult{
				Input:          tc.Input,
				ExpectedOutput: tc.ExpectedOutput,
				ActualOutput:   "",
				Status:         "Failed",
				IsSample:       tc.IsSample,
				// ErrorMessage:   *output.Error, // TestCaseResult doesn't have ErrorMessage in domain
			})
			continue
		}

		// Compare output
		actualJSON, _ := json.Marshal(output.Output)
		passed := s.compareOutputs(string(actualJSON), tc.ExpectedOutput)

		status := "Failed"
		if passed {
			status = "Passed"
			passedCount++
		}

		testResults = append(testResults, domain.TestCaseResult{
			Input:          tc.Input,
			ExpectedOutput: tc.ExpectedOutput,
			ActualOutput:   string(actualJSON),
			Status:         status,
			IsSample:       tc.IsSample,
		})
	}

	// Determine overall status
	overallStatus := domain.SubmissionStatusWrongAnswer
	if passedCount == len(testCases) {
		overallStatus = domain.SubmissionStatusAccepted
	}

	return &ExecutionResult{
		Status:      overallStatus,
		TestResults: testResults,
		TotalTests:  len(testCases),
		PassedTests: passedCount,
	}, nil
}

// compareOutputs compares actual vs expected (handles exact match for now)
func (s *ExecutionService) compareOutputs(actual, expected string) bool {
	// Normalize JSON (remove whitespace differences)
	var actualData, expectedData interface{}

	if err := json.Unmarshal([]byte(actual), &actualData); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(expected), &expectedData); err != nil {
		return false
	}

	actualNorm, _ := json.Marshal(actualData)
	expectedNorm, _ := json.Marshal(expectedData)

	return string(actualNorm) == string(expectedNorm)
}
