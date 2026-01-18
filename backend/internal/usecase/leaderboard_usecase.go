package usecase

import (
	"fmt"

	"github.com/prabalesh/loco/backend/internal/domain"
	"go.uber.org/zap"
)

type LeaderboardUsecase struct {
	userRepo domain.UserRepository
	logger   *zap.Logger
}

func NewLeaderboardUsecase(userRepo domain.UserRepository, logger *zap.Logger) *LeaderboardUsecase {
	return &LeaderboardUsecase{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (u *LeaderboardUsecase) GetLeaderboard(limit int) ([]domain.LeaderboardEntry, error) {
	if limit <= 0 {
		limit = 100 // Default limit
	}

	entries, err := u.userRepo.GetLeaderboard(limit)
	if err != nil {
		u.logger.Error("Failed to fetch leaderboard", zap.Error(err))
		return nil, fmt.Errorf("could not fetch leaderboard")
	}

	return entries, nil
}
