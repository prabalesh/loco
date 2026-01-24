package dto

import (
	"time"

	"github.com/prabalesh/loco/backend/internal/domain"
)

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
	ID                 int                      `json:"id"`
	Username           string                   `json:"username"`
	Email              string                   `json:"email"`
	IsVerified         bool                     `json:"is_verified"`
	CreatedAt          time.Time                `json:"created_at"`
	XP                 int                      `json:"xp"`
	Level              int                      `json:"level"`
	Stats              UserStats                `json:"stats"`
	RecentProblems     []domain.Problem         `json:"recent_problems"`
	SubmissionHeatmap  []domain.HeatmapEntry    `json:"submission_heatmap"`
	SolvedDistribution []domain.DifficultyStat  `json:"solved_distribution"`
	Achievements       []domain.UserAchievement `json:"achievements"`
}

type UserStats struct {
	TotalSubmissions    int                     `json:"total_submissions"`
	AcceptedSubmissions int                     `json:"accepted_submissions"`
	ProblemsSolved      int                     `json:"problems_solved"`
	AcceptanceRate      float64                 `json:"acceptance_rate"`
	Rank                int                     `json:"rank"`
	Streak              int                     `json:"streak"`
	SolvedDistribution  []domain.DifficultyStat `json:"solved_distribution"`
}
