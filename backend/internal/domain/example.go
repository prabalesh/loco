package domain

import "time"

// ProblemExample represents a visible example for a problem
type ProblemExample struct {
	ID          int       `json:"id" db:"id"`
	ProblemID   int       `json:"problem_id" db:"problem_id"`
	Input       string    `json:"input" db:"input"`
	Output      string    `json:"output" db:"output"`
	Explanation string    `json:"explanation" db:"explanation"`
	OrderIndex  int       `json:"order_index" db:"order_index"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
