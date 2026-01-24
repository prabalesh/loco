package validator

import (
	"fmt"

	"github.com/prabalesh/loco/backend/internal/domain/dto"
)

// ValidateCreateTestCaseRequest validates create test case request
func ValidateCreateTestCaseRequest(req *dto.CreateTestCaseRequest) map[string]string {
	errors := make(map[string]string)

	if req.ProblemID <= 0 {
		errors["problem_id"] = "problem_id is required"
	}
	if req.Input == "" {
		errors["input"] = "input is required"
	}
	if len(req.Input) > 10000 {
		errors["input"] = "input too long (max 10000 chars)"
	}
	if req.ExpectedOutput == "" {
		errors["expected_output"] = "expected_output is required"
	}
	if len(req.ExpectedOutput) > 10000 {
		errors["expected_output"] = "expected_output too long (max 10000 chars)"
	}
	if req.OrderIndex < 0 {
		errors["order_index"] = "order_index must be non-negative"
	}

	return errors
}

// ValidateUpdateTestCaseRequest validates update test case request
func ValidateUpdateTestCaseRequest(req *dto.UpdateTestCaseRequest) map[string]string {
	errors := make(map[string]string)

	if req.Input != "" && len(req.Input) > 10000 {
		errors["input"] = "input too long (max 10000 chars)"
	}
	if req.ExpectedOutput != "" && len(req.ExpectedOutput) > 10000 {
		errors["expected_output"] = "expected_output too long (max 10000 chars)"
	}
	if req.OrderIndex != nil && *req.OrderIndex < 0 {
		errors["order_index"] = "order_index must be non-negative"
	}

	return errors
}

// ValidateReorderTestCasesRequest validates reorder request
func ValidateReorderTestCasesRequest(req *dto.ReorderTestCasesRequest) map[string]string {
	errors := make(map[string]string)

	if req.ProblemID <= 0 {
		errors["problem_id"] = "problem_id is required"
	}
	if len(req.TestCases) == 0 {
		errors["test_cases"] = "test_cases is required"
	}
	for i, tc := range req.TestCases {
		if tc.ID <= 0 {
			errors[fmt.Sprintf("test_cases[%d].id", i)] = "id is required"
		}
		if tc.OrderIndex < 0 {
			errors[fmt.Sprintf("test_cases[%d].order_index", i)] = "order_index must be non-negative"
		}
	}

	return errors
}
