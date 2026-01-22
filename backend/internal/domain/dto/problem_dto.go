package dto

import "gorm.io/datatypes"

// CreateProblemRequest defines the payload for creating a problem
type CreateProblemRequest struct {
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
	TimeLimit   int    `json:"time_limit"`
	MemoryLimit int    `json:"memory_limit"`
	Status      string `json:"status"`
	Visibility  string `json:"visibility"`
	IsActive    bool   `json:"is_active"`
	TagIDs      []int  `json:"tag_ids"`
	CategoryIDs []int  `json:"category_ids"`

	// Core V2 Fields
	FunctionName            string          `json:"function_name"`
	ReturnType              string          `json:"return_type"`
	Parameters              *datatypes.JSON `json:"parameters"` // Using datatypes.JSON for flexibility
	ValidationType          string          `json:"validation_type"`
	ExpectedTimeComplexity  string          `json:"expected_time_complexity"`
	ExpectedSpaceComplexity string          `json:"expected_space_complexity"`

	TestCases []TestCaseInput `json:"test_cases"`
}

// UpdateProblemRequest defines the payload for updating a problem
type UpdateProblemRequest struct {
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
	TimeLimit   int    `json:"time_limit"`
	MemoryLimit int    `json:"memory_limit"`
	Status      string `json:"status"`
	Visibility  string `json:"visibility"`
	IsActive    *bool  `json:"is_active"`
	TagIDs      []int  `json:"tag_ids"`
	CategoryIDs []int  `json:"category_ids"`

	FunctionName            *string         `json:"function_name"`
	ReturnType              *string         `json:"return_type"`
	Parameters              *datatypes.JSON `json:"parameters"`
	ValidationType          *string         `json:"validation_type"`
	ExpectedTimeComplexity  *string         `json:"expected_time_complexity"`
	ExpectedSpaceComplexity *string         `json:"expected_space_complexity"`
}

type TestCaseInput struct {
	Input          interface{} `json:"input"`
	ExpectedOutput interface{} `json:"expected_output"`
	IsSample       bool        `json:"is_sample"`
	InputSize      *int        `json:"input_size"`
	TimeLimitMs    *int        `json:"time_limit_ms"`
	MemoryLimitMb  *int        `json:"memory_limit_mb"`
}

type ListProblemsRequest struct {
	Page       int      `json:"page"`
	Limit      int      `json:"limit"`
	Difficulty string   `json:"difficulty"`
	Status     string   `json:"status"`
	Visibility string   `json:"visibility"`
	Search     string   `json:"search"`
	Tags       []string `json:"tags"`
	Categories []string `json:"categories"`
}
