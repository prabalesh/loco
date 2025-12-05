package domain

// CreateTestCaseRequest for creating new test cases
type CreateTestCaseRequest struct {
	ProblemID        int              `json:"problem_id"`
	Input            string           `json:"input"`
	ExpectedOutput   string           `json:"expected_output"`
	IsSample         bool             `json:"is_sample"`
	ValidationConfig ValidationConfig `json:"validation_config"`
	OrderIndex       int              `json:"order_index"`
}

// UpdateTestCaseRequest for updating test cases
type UpdateTestCaseRequest struct {
	Input            string           `json:"input"`
	ExpectedOutput   string           `json:"expected_output"`
	IsSample         *bool            `json:"is_sample"`
	ValidationConfig ValidationConfig `json:"validation_config"`
	OrderIndex       *int             `json:"order_index"`
}

// ReorderTestCasesRequest for bulk reordering
type ReorderTestCasesRequest struct {
	ProblemID int             `json:"problem_id"`
	TestCases []TestCaseOrder `json:"test_cases"`
}

type TestCaseOrder struct {
	ID         int `json:"id"`
	OrderIndex int `json:"order_index"`
}

// ListTestCasesRequest for listing test cases
type ListTestCasesRequest struct {
	ProblemID int   `json:"problem_id"`
	IsSample  *bool `json:"is_sample"`
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
}
