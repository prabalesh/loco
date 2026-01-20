package dto

import (
	"time"

	"github.com/prabalesh/loco/backend/internal/domain"
)

type CreateSubmissionRequest struct {
	ProblemID  int    `json:"problem_id" validate:"required"`
	LanguageID int    `json:"language_id" validate:"required"`
	Code       string `json:"code" validate:"required"`
}

type SubmissionResponse struct {
	ID              int                     `json:"id"`
	UserID          int                     `json:"user_id"`
	ProblemID       int                     `json:"problem_id"`
	LanguageID      int                     `json:"language_id"`
	FunctionCode    string                  `json:"function_code"`
	Status          domain.SubmissionStatus `json:"status"`
	ErrorMessage    string                  `json:"error_message,omitempty"`
	Runtime         int                     `json:"runtime"`
	Memory          int                     `json:"memory"`
	PassedTestCases int                     `json:"passed_test_cases"`
	TotalTestCases  int                     `json:"total_test_cases"`
	CreatedAt       time.Time               `json:"created_at"`
	TestCaseResults domain.TestCaseResults  `json:"test_case_results,omitempty"`
	User            *UserResponse           `json:"user,omitempty"`
	Problem         *ProblemResponse        `json:"problem,omitempty"`
	Language        *LanguageResponse       `json:"language,omitempty"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type ProblemResponse struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

type LanguageResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type SubmissionsListResponse struct {
	Data  []SubmissionResponse `json:"data"`
	Total int64                `json:"total"`
	Page  int                  `json:"page"`
	Limit int                  `json:"limit"`
}

func ToSubmissionResponse(s *domain.Submission) SubmissionResponse {
	resp := SubmissionResponse{
		ID:              s.ID,
		UserID:          s.UserID,
		ProblemID:       s.ProblemID,
		LanguageID:      s.LanguageID,
		FunctionCode:    s.FunctionCode,
		Status:          s.Status,
		ErrorMessage:    s.ErrorMessage,
		Runtime:         s.Runtime,
		Memory:          s.Memory,
		PassedTestCases: s.PassedTestCases,
		TotalTestCases:  s.TotalTestCases,
		CreatedAt:       s.CreatedAt,
		TestCaseResults: s.TestCaseResults,
	}

	if s.User != nil {
		resp.User = &UserResponse{
			ID:       s.User.ID,
			Username: s.User.Username,
		}
	}

	if s.Problem != nil {
		resp.Problem = &ProblemResponse{
			ID:    s.Problem.ID,
			Title: s.Problem.Title,
			Slug:  s.Problem.Slug,
		}
	}

	if s.Language != nil {
		resp.Language = &LanguageResponse{
			ID:   s.Language.ID,
			Name: s.Language.Name,
			Slug: s.Language.Slug,
		}
	}

	return resp
}
