package usecase

import (
	"errors"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/usecase/interfaces"
	"go.uber.org/zap"
)

type UserUsecase struct {
	userRepo interfaces.UserRepository
	logger   *zap.Logger
}

func NewUserUsecase(userRepo interfaces.UserRepository, logger *zap.Logger) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
		logger:   logger,
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

func (u *UserUsecase) GetUserProfileByUsername(username string) (*domain.User, error) {
	user, err := u.userRepo.GetByUsername(username)
	if err != nil {
		u.logger.Error("Failed to get user by username",
			zap.Error(err),
			zap.String("username", username),
		)
		return nil, errors.New("user not found")
	}

	return user, nil
}
