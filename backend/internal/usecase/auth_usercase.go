package usecase

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/domain/validator"
	"github.com/prabalesh/loco/backend/internal/infrastructure/auth"
	"github.com/prabalesh/loco/backend/internal/usecase/interfaces"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo   interfaces.UserRepository
	jwtService *auth.JWTService
	logger     *zap.Logger
}

func NewAuthUsecase(userRepo interfaces.UserRepository, jwtService *auth.JWTService, logger *zap.Logger) *AuthUsecase {
	return &AuthUsecase{
		userRepo:   userRepo,
		jwtService: jwtService,
		logger:     logger,
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

type TokenPair struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  time.Duration
	RefreshExpiresAt time.Duration
}

func (u *AuthUsecase) Login(req *domain.LoginRequest) (*domain.User, *TokenPair, error) {
	// validation
	if validationErrors := validator.ValidateLoginRequest(req); len(validationErrors) > 0 {
		u.logger.Warn("Registration validation failed",
			zap.Any("errors", validationErrors),
		)
		return nil, nil, &ValidationError{Errors: validationErrors}
	}

	// get user by email
	existingUser, err := u.userRepo.GetByEmail(req.Email)
	if err != nil && !isNotFoundError(err) {
		return nil, nil, errors.New("internal server error")
	}

	// verify password
	if err != nil || !verifyPassword(existingUser.PasswordHash, req.Password) {
		u.logger.Warn("Login failed: invalid password", zap.String("email", req.Email))
		return nil, nil, errors.New("invalid email or password")
	}

	// check if account is active
	if !existingUser.IsActive {
		u.logger.Warn("Login failed: account deactivated", zap.String("email", req.Email))
		return nil, nil, errors.New("account is deactivated")
	}

	// generate tokens
	accessToken, accessTokenExpires, err := u.jwtService.GenerateAccessToken(existingUser.ID, existingUser.Email, existingUser.Role)
	if err != nil {
		u.logger.Warn("Login failed: error occured while creating access token", zap.Error(err))
		return nil, nil, errors.New("internal server error")
	}

	refreshToken, refreshTokenExpires, err := u.jwtService.GenerateRefreshToken(existingUser.ID, existingUser.Email)
	if err != nil {
		u.logger.Warn("Login failed: error occured while creating refresh token", zap.Error(err))
		return nil, nil, errors.New("internal server error")
	}

	tokenPair := TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresAt:  accessTokenExpires,
		RefreshExpiresAt: refreshTokenExpires,
	}

	return existingUser, &tokenPair, nil
}

// RefreshAccessToken generates new access token from refresh token
func (u *AuthUsecase) RefreshAccessToken(refreshToken string) (string, time.Duration, error) {
	// Validate refresh token
	claims, err := u.jwtService.ValidateToken(refreshToken, true)
	if err != nil {
		return "", 0, errors.New("invalid refresh token")
	}

	// Check if token exists in database and not revoked (if you're storing them)
	// tokenHash := auth.HashToken(refreshToken)
	// dbToken, err := u.refreshTokenRepo.GetByTokenHash(tokenHash)
	// if err != nil || dbToken == nil {
	//     return "", 0, errors.New("refresh token not found or revoked")
	// }

	// Get user (to get latest role in case it changed)
	user, err := u.userRepo.GetByID(claims.UserID)
	if err != nil {
		return "", 0, errors.New("user not found")
	}

	// Check if user is still active
	if !user.IsActive {
		return "", 0, errors.New("account is deactivated")
	}

	// Generate new access token
	accessToken, expiresAt, err := u.jwtService.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		u.logger.Error("Failed to generate access token", zap.Error(err))
		return "", 0, errors.New("failed to refresh token")
	}

	return accessToken, expiresAt, nil
}

// Logout revokes refresh token
func (u *AuthUsecase) Logout(refreshToken string) error {
	// If you're storing refresh tokens in DB:
	// tokenHash := auth.HashToken(refreshToken)
	// return u.refreshTokenRepo.Revoke(tokenHash)

	// For now, just log
	u.logger.Info("User logged out")
	return nil
}

// GetCurrentUser returns user info by ID
func (u *AuthUsecase) GetCurrentUser(userID int) (*domain.User, error) {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func verifyPassword(hashedPassword, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false
	}
	return true
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
