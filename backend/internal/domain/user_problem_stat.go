package domain

import "time"

// UserProblemStats tracks a user's progress on a specific problem
type UserProblemStats struct {
	UserID           int        `json:"user_id" db:"user_id"`
	ProblemID        int        `json:"problem_id" db:"problem_id"`
	Status           string     `json:"status" db:"status"` // solved, attempted, unsolved
	Attempts         int        `json:"attempts" db:"attempts"`
	FirstSolvedAt    *time.Time `json:"first_solved_at,omitempty" db:"first_solved_at"`       // nullable
	BestSubmissionID *int       `json:"best_submission_id,omitempty" db:"best_submission_id"` // nullable
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}
