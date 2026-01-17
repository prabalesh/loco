package usecase

import (
	"errors"

	"github.com/prabalesh/loco/backend/internal/domain"
	"go.uber.org/zap"
)

type UserUsecase struct {
	userRepo       domain.UserRepository
	submissionRepo domain.SubmissionRepository
	logger         *zap.Logger
}

func NewUserUsecase(userRepo domain.UserRepository, submissionRepo domain.SubmissionRepository, logger *zap.Logger) *UserUsecase {
	return &UserUsecase{
		userRepo:       userRepo,
		submissionRepo: submissionRepo,
		logger:         logger,
	}
}

// GetUserProfile returns user profile by ID
func (u *UserUsecase) GetUserProfile(userID int) (*domain.User, error) {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		u.logger.Error("Failed to get user by ID",
			zap.Error(err),
			zap.Int("user_id", userID),
		)
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (u *UserUsecase) GetUserProfileByUsername(username string) (*domain.UserProfileResponse, error) {
	user, err := u.userRepo.GetByUsername(username)
	if err != nil {
		u.logger.Error("Failed to get user by username",
			zap.Error(err),
			zap.String("username", username),
		)
		return nil, errors.New("user not found")
	}

	// Fetch stats
	stats, err := u.getUserStats(user.ID)
	if err != nil {
		u.logger.Error("Failed to get user stats", zap.Error(err))
	}

	resp := user.ToUserProfileResponse(stats)
	return &resp, nil
}

func (u *UserUsecase) getUserStats(userID int) (domain.UserStats, error) {
	totalSubmissions, err := u.submissionRepo.CountByUser(userID)
	if err != nil {
		return domain.UserStats{}, err
	}

	acceptedSubmissions, err := u.submissionRepo.CountAcceptedByUser(userID)
	if err != nil {
		return domain.UserStats{}, err
	}

	problemsSolved, err := u.submissionRepo.CountProblemsSolvedByUser(userID)
	if err != nil {
		return domain.UserStats{}, err
	}

	acceptanceRate := 0.0
	if totalSubmissions > 0 {
		acceptanceRate = float64(acceptedSubmissions) / float64(totalSubmissions) * 100
	}

	return domain.UserStats{
		TotalSubmissions:    int(totalSubmissions),
		AcceptedSubmissions: int(acceptedSubmissions),
		ProblemsSolved:      int(problemsSolved),
		AcceptanceRate:      acceptanceRate,
	}, nil
}
