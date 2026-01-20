package domain

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type SubmissionStatus string

const (
	SubmissionStatusPending             SubmissionStatus = "Pending"
	SubmissionStatusProcessing          SubmissionStatus = "Processing"
	SubmissionStatusAccepted            SubmissionStatus = "Accepted"
	SubmissionStatusWrongAnswer         SubmissionStatus = "Wrong Answer"
	SubmissionStatusTimeLimitExceeded   SubmissionStatus = "Time Limit Exceeded"
	SubmissionStatusMemoryLimitExceeded SubmissionStatus = "Memory Limit Exceeded"
	SubmissionStatusRuntimeError        SubmissionStatus = "Runtime Error"
	SubmissionStatusCompilationError    SubmissionStatus = "Compilation Error"
	SubmissionStatusInternalError       SubmissionStatus = "Internal Error"
)

type TestCaseResult struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	ActualOutput   string `json:"actual_output"`
	Status         string `json:"status"` // "Passed" or "Failed"
	IsSample       bool   `json:"is_sample"`
}

// Slice of TestCaseResult
type TestCaseResults []TestCaseResult

func (t TestCaseResults) Value() (driver.Value, error) {
	if len(t) == 0 {
		return nil, nil
	}
	return json.Marshal(t)
}

func (t *TestCaseResults) Scan(value interface{}) error {
	if value == nil {
		*t = make(TestCaseResults, 0)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, t)
}

type Submission struct {
	ID              int              `json:"id" gorm:"primaryKey"`
	UserID          int              `json:"user_id" gorm:"not null"`
	ProblemID       int              `json:"problem_id" gorm:"not null"`
	LanguageID      int              `json:"language_id" gorm:"not null"`
	Code            string           `json:"code,omitempty" gorm:"type:text;not null"`
	FunctionCode    string           `json:"function_code" gorm:"type:text;default:''"`
	Status          SubmissionStatus `json:"status" gorm:"type:varchar(50);default:'Pending'"`
	ErrorMessage    string           `json:"error_message,omitempty" gorm:"type:text"` // For compile/runtime errors
	Runtime         int              `json:"runtime" gorm:"default:0"`                 // in milliseconds
	Memory          int              `json:"memory" gorm:"default:0"`                  // in kilobytes
	PassedTestCases int              `json:"passed_test_cases" gorm:"default:0"`
	TotalTestCases  int              `json:"total_test_cases" gorm:"default:0"`
	CreatedAt       time.Time        `json:"created_at" gorm:"autoCreateTime"`

	// Detailed results
	TestCaseResults TestCaseResults `json:"test_case_results" gorm:"type:jsonb"`

	// Queue metadata
	QueuedAt    *time.Time `json:"queued_at,omitempty" gorm:"index"`    // When job was enqueued
	ProcessedAt *time.Time `json:"processed_at,omitempty" gorm:"index"` // When job was processed

	// Admin context
	IsAdminSubmission      bool `json:"is_admin_submission" gorm:"default:false"`      // Distinguishes admin test submissions
	IsValidationSubmission bool `json:"is_validation_submission" gorm:"default:false"` // Specifically for validating problem-language solution
	SubmittedBy            *int `json:"submitted_by,omitempty" gorm:"index"`           // Admin user ID if admin submission

	// Associations
	User     *User     `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Admin    *User     `json:"admin,omitempty" gorm:"foreignKey:SubmittedBy;references:ID;constraint:OnDelete:SET NULL"`
	Problem  *Problem  `json:"problem,omitempty" gorm:"foreignKey:ProblemID;references:ID"`
	Language *Language `json:"language,omitempty" gorm:"foreignKey:LanguageID;references:ID"`
}

func (s *Submission) Sanitize() {
	for i := range s.TestCaseResults {
		if !s.TestCaseResults[i].IsSample {
			s.TestCaseResults[i].Input = ""
			s.TestCaseResults[i].ExpectedOutput = ""
			s.TestCaseResults[i].ActualOutput = ""
		}
	}
}

type CreateSubmissionRequest struct {
	ProblemID  int    `json:"problem_id" validate:"required"`
	LanguageID int    `json:"language_id" validate:"required"`
	Code       string `json:"code" validate:"required"`
}

type SubmissionStats struct {
	TotalSubmissions int `json:"total_submissions"`
	Accepted         int `json:"accepted"`
	WrongAnswer      int `json:"wrong_answer"`
	RuntimeError     int `json:"runtime_error"`
}

type DailySubmissionStat struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// RunCodeResult represents the result of running code without creating a submission
type RunCodeResult struct {
	Status          SubmissionStatus `json:"status"`
	ErrorMessage    string           `json:"error_message,omitempty"`
	PassedTestCases int              `json:"passed_test_cases"`
	TotalTestCases  int              `json:"total_test_cases"`
	Results         []TestCaseResult `json:"results"`
}

func (r *RunCodeResult) Sanitize() {
	for i := range r.Results {
		if !r.Results[i].IsSample {
			r.Results[i].Input = ""
			r.Results[i].ExpectedOutput = ""
			r.Results[i].ActualOutput = ""
		}
	}
}

type RunCodeRequest struct {
	LanguageID int    `json:"language_id" validate:"required"`
	Code       string `json:"code" validate:"required"`
}
