package domain

import "time"

// ==================== REQUEST DTOs ====================
// ==================== RESPONSE DTOs ====================

type UserResponse struct {
	ID            int       `json:"id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	Role          string    `json:"role"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	XP            int       `json:"xp"`
	Level         int       `json:"level"`
}

type UserProfileResponse struct {
	ID                 int              `json:"id"`
	Username           string           `json:"username"`
	IsVerified         bool             `json:"is_verified"`
	CreatedAt          time.Time        `json:"created_at"`
	XP                 int              `json:"xp"`
	Level              int              `json:"level"`
	Stats              UserStats        `json:"stats"`
	RecentProblems     []Problem        `json:"recent_problems"`
	SubmissionHeatmap  []HeatmapEntry   `json:"submission_heatmap"`
	SolvedDistribution []DifficultyStat `json:"solved_distribution"`
}

type UserStats struct {
	TotalSubmissions    int              `json:"total_submissions"`
	AcceptedSubmissions int              `json:"accepted_submissions"`
	ProblemsSolved      int              `json:"problems_solved"`
	AcceptanceRate      float64          `json:"acceptance_rate"`
	Rank                int              `json:"rank"`
	Streak              int              `json:"streak"`
	SolvedDistribution  []DifficultyStat `json:"solved_distribution"`
}

type DifficultyStat struct {
	Difficulty string `json:"difficulty"`
	Count      int    `json:"count"`
}

type HeatmapEntry struct {
	Date  string `json:"date"` // YYYY-MM-DD
	Count int    `json:"count"`
}
