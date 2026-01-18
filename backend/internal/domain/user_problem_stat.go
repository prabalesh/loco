package domain

import "time"

// UserProblemStats tracks a user's progress on a specific problem
type UserProblemStats struct {
	UserID           int        `json:"user_id" gorm:"primaryKey;autoIncrement:false"`
	ProblemID        int        `json:"problem_id" gorm:"primaryKey;autoIncrement:false"`
	Status           string     `json:"status" gorm:"size:50;default:'unsolved'"` // solved, attempted, unsolved
	Attempts         int        `json:"attempts" gorm:"default:0"`
	FirstSolvedAt    *time.Time `json:"first_solved_at,omitempty"`
	BestSubmissionID *int       `json:"best_submission_id,omitempty"`
	CreatedAt        time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}
