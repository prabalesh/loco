package dto

import (
	"github.com/prabalesh/loco/backend/internal/domain"
)

// ToUserResponse converts User entity to UserResponse DTO
func ToUserResponse(u *domain.User) UserResponse {
	return UserResponse{
		ID:            u.ID,
		Email:         u.Email,
		Username:      u.Username,
		Role:          u.Role,
		EmailVerified: u.EmailVerified,
		CreatedAt:     u.CreatedAt,
		XP:            u.XP,
		Level:         u.Level,
	}
}

func ToUserProfileResponse(u *domain.User, stats UserStats, recentProblems []domain.Problem, heatmap []domain.HeatmapEntry, distribution []domain.DifficultyStat, achievements []domain.UserAchievement) UserProfileResponse {
	return UserProfileResponse{
		ID:                 u.ID,
		Username:           u.Username,
		Email:              u.Email,
		IsVerified:         u.EmailVerified,
		CreatedAt:          u.CreatedAt,
		XP:                 u.XP,
		Level:              u.Level,
		Stats:              stats,
		RecentProblems:     recentProblems,
		SubmissionHeatmap:  heatmap,
		SolvedDistribution: distribution,
		Achievements:       achievements,
	}
}
