package usecase

import (
	"errors"
	"fmt"
	"strings"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/domain/validator"
	"github.com/prabalesh/loco/backend/internal/usecase/interfaces"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo interfaces.UserRepository
	logger   *zap.Logger
}

func NewAuthUsecase(userRepo interfaces.UserRepository, logger *zap.Logger) *AuthUsecase {
	return &AuthUsecase{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (u *AuthUsecase) Register(req *domain.RegisterRequest) (*domain.User, error) {
	// validation
	if validationErrors := validator.ValidateRegisterRequest(req); len(validationErrors) > 0 {
		u.logger.Warn("Registration validation failed",
			zap.Any("errors", validationErrors),
		)
		return nil, &ValidationError{Errors: validationErrors}
	}

	existingUser, err := u.userRepo.GetByEmail(req.Email)
	if err != nil && !isNotFoundError(err) {
		// Database error (not "not found")
		u.logger.Error("Failed to check email existence",
			zap.Error(err),
			zap.String("email", req.Email),
		)
		return nil, errors.New("failed to check email availability")
	}
	if existingUser != nil {
		u.logger.Warn("Registration failed: email already exists",
			zap.String("email", req.Email),
		)
		return nil, errors.New("email already registered")
	}

	existingUsername, err := u.userRepo.GetByUsername(req.Username)
	if err != nil && !isNotFoundError(err) {
		u.logger.Error("Failed to check username existence",
			zap.Error(err),
			zap.String("username", req.Username),
		)
		return nil, errors.New("failed to check username availability")
	}
	if existingUsername != nil {
		u.logger.Warn("Registration failed: username already taken",
			zap.String("username", req.Username),
		)
		return nil, errors.New("username already taken")
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		u.logger.Error("Failed to hash password",
			zap.Error(err),
		)
		return nil, errors.New("failed to create account")
	}

	user := &domain.User{
		Email:         req.Email,
		Username:      req.Username,
		PasswordHash:  hashedPassword,
		Role:          "user",
		IsActive:      true,
		EmailVerified: false,
	}

	if err := u.userRepo.Create(user); err != nil {
		u.logger.Error("Failed to create user in database",
			zap.Error(err),
			zap.String("email", req.Email),
			zap.String("username", req.Username),
		)
		return nil, errors.New("failed to create account")
	}

	u.logger.Info("User registered successfully",
		zap.Int("user_id", user.ID),
		zap.String("email", user.Email),
		zap.String("username", user.Username),
	)

	return user, nil
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func isNotFoundError(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "not found")
}

type ValidationError struct {
	Errors map[string]string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %v", e.Errors)
}
