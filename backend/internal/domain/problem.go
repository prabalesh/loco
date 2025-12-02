package domain

import "time"

// Problem entity
type Problem struct {
	ID              int       `json:"id" db:"id"`
	Title           string    `json:"title" db:"title"`
	Slug            string    `json:"slug" db:"slug"`
	Description     string    `json:"description" db:"description"`
	Difficulty      string    `json:"difficulty" db:"difficulty"`
	TimeLimit       int       `json:"time_limit" db:"time_limit"`
	MemoryLimit     int       `json:"memory_limit" db:"memory_limit"`
	ValidatorType   string    `json:"validator_type" db:"validator_type"`
	InputFormat     string    `json:"input_format" db:"input_format"`
	OutputFormat    string    `json:"output_format" db:"output_format"`
	Constraints     string    `json:"constraints" db:"constraints"`
	Status          string    `json:"status" db:"status"`
	Visibility      string    `json:"visibility" db:"visibility"`
	IsActive        bool      `json:"is_active" db:"is_active"`
	AcceptanceRate  float64   `json:"acceptance_rate" db:"acceptance_rate"`
	TotalSubmission int       `json:"total_submissions" db:"total_submissions"`
	TotalAccepted   int       `json:"total_accepted" db:"total_accepted"`
	CreatedBy       *int      `json:"created_by" db:"created_by"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
