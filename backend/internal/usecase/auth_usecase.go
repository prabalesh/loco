package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/domain/dto"
	"github.com/prabalesh/loco/backend/internal/domain/uerror"
	"github.com/prabalesh/loco/backend/internal/domain/validator"
	"github.com/prabalesh/loco/backend/internal/infrastructure/auth"
	"github.com/prabalesh/loco/backend/internal/infrastructure/email"
	"github.com/prabalesh/loco/backend/pkg/config"
	"github.com/prabalesh/loco/backend/pkg/utils"
	"go.uber.org/zap"
)

type AuthUsecase struct {
	userRepo     domain.UserRepository
	jwtService   *auth.JWTService
	emailService *email.EmailService
	cfg          *config.Config
	logger       *zap.Logger
}

func NewAuthUsecase(userRepo domain.UserRepository, jwtService *auth.JWTService, emailService *email.EmailService, cfg *config.Config, logger *zap.Logger) *AuthUsecase {
	return &AuthUsecase{
		userRepo:     userRepo,
		jwtService:   jwtService,
		emailService: emailService,
		cfg:          cfg,
		logger:       logger,
	}
}

func (u *AuthUsecase) Register(req *dto.RegisterRequest) (*domain.User, error) {
	// validation
	if validationErrors := validator.ValidateRegisterRequest(req); len(validationErrors) > 0 {
		u.logger.Warn("Registration validation failed",
			zap.Any("errors", validationErrors),
		)
		return nil, &uerror.ValidationError{Errors: validationErrors}
	}

	// checking if email already exists
	existingUser, err := u.userRepo.GetByEmail(req.Email)
	if err != nil && !uerror.IsNotFoundError(err) {
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

	// check username exists
	existingUsername, err := u.userRepo.GetByUsername(req.Username)
	if err != nil && !uerror.IsNotFoundError(err) {
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

	// hash password
	hashedPassword, err := utils.HashPassword(req.Password)
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

	// send verification email
	if err := u.sendVerificationEmail(context.Background(), user); err != nil {
		u.logger.Error("Failed to send verification email", zap.Error(err))
		// Don't fail registration if email fails
	}

	u.logger.Info("User registered successfully",
		zap.Int("user_id", user.ID),
		zap.String("email", user.Email),
		zap.String("username", user.Username),
	)

	return user, nil
}

func (u *AuthUsecase) VerifyEmail(ctx context.Context, req *dto.VerifyEmailRequest) error {
	user, err := u.userRepo.GetByVerificationToken(req.Token)
	if err != nil {
		u.logger.Warn("User not found for token")
		return errors.New("user not found")
	}

	if user.EmailVerified {
		u.logger.Info("Email already verified", zap.Int("user_id", user.ID))
		return nil // Already verified
	}

	// Check max attempts
	if user.EmailVerificationAttempts >= u.cfg.Email.MaxTokenAttempts {
		return uerror.ErrMaxTokenAttemptsExceeded
	}

	// Check if token exists and hasn't expired
	if user.EmailVerificationToken == nil || user.EmailVerificationTokenExpiresAt == nil {
		return uerror.ErrInvalidToken
	}

	if time.Now().After(*user.EmailVerificationTokenExpiresAt) {
		return uerror.ErrInvalidToken
	}

	// Verify token
	if *user.EmailVerificationToken != req.Token {
		// Increment attempts
		newAttempts := user.EmailVerificationAttempts + 1
		u.userRepo.UpdateVerificationAttempts(user.ID, newAttempts)

		if newAttempts >= u.cfg.Email.MaxTokenAttempts {
			return uerror.ErrMaxTokenAttemptsExceeded
		}

		return uerror.ErrInvalidToken
	}

	// Mark email as verified
	if err := u.userRepo.VerifyEmail(user.ID); err != nil {
		u.logger.Error("Failed to verify email", zap.Error(err))
		return errors.New("failed to verify email")
	}

	u.logger.Info("Email verified successfully", zap.String("email", user.Email))
	return nil
}

// resends email verification link with cooldown
func (u *AuthUsecase) ResendVerificationEmail(ctx context.Context, req *dto.ResendVerificationRequest) error {
	user, err := u.userRepo.GetByEmail(req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	if user.EmailVerified {
		return nil // Already verified
	}

	// Check cooldown
	if user.EmailVerificationLastSentAt != nil {
		cooldownDuration := time.Duration(u.cfg.Email.ResendCooldownMinutes) * time.Minute
		nextAllowedTime := user.EmailVerificationLastSentAt.Add(cooldownDuration)

		if time.Now().Before(nextAllowedTime) {
			remainingSeconds := int(time.Until(nextAllowedTime).Seconds())
			return fmt.Errorf("%w: %d seconds remaining", uerror.ErrResendCooldown, remainingSeconds)
		}
	}

	// Check max attempts
	if user.EmailVerificationAttempts >= u.cfg.Email.MaxTokenAttempts {
		return uerror.ErrMaxTokenAttemptsExceeded
	}

	return u.sendVerificationEmail(ctx, user)
}

func (u *AuthUsecase) Login(req *dto.LoginRequest) (*domain.User, *dto.TokenPair, error) {
	// validation
	if validationErrors := validator.ValidateLoginRequest(req); len(validationErrors) > 0 {
		u.logger.Warn("Registration validation failed",
			zap.Any("errors", validationErrors),
		)
		return nil, nil, &uerror.ValidationError{Errors: validationErrors}
	}

	// get user by email
	existingUser, err := u.userRepo.GetByEmail(req.Email)
	if err != nil && !uerror.IsNotFoundError(err) {
		return nil, nil, errors.New("internal server error")
	}

	// verify password
	if err != nil || !utils.VerifyPassword(existingUser.PasswordHash, req.Password) {
		u.logger.Warn("Login failed: invalid password", zap.String("email", req.Email))
		return nil, nil, errors.New("invalid email or password")
	}

	if !existingUser.EmailVerified {
		u.logger.Warn("Login attempt with unverified email", zap.String("email", req.Email))
		return nil, nil, uerror.ErrEmailNotVerified
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

	tokenPair := dto.TokenPair{
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
		fmt.Println(err)
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

// Starts the password reset process: finds user by email, generates token, sends email
func (u *AuthUsecase) ForgotPassword(ctx context.Context, email string) error {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		// For security, do not reveal if user exists or not
		return nil
	}

	// Implement rate limiting using sentAt timestamp
	if user.PasswordResetSentAt != nil && time.Since(*user.PasswordResetSentAt) < time.Duration(u.cfg.Email.ResendCooldownMinutes)*time.Minute {
		return uerror.ErrResendCooldown
	}

	// Generate secure token
	token, err := utils.GenerateToken(64)
	if err != nil {
		return err
	}

	expiresAt := time.Now().Add(time.Duration(u.cfg.Email.PasswordResetExpiryMinutes) * time.Minute)

	// Update DB with token and expiry
	if err := u.userRepo.UpdatePasswordResetToken(user.ID, token, expiresAt, time.Now()); err != nil {
		return err
	}

	// Send reset email
	if err := u.emailService.SendPasswordResetEmail(ctx, user.Email, user.Username, token); err != nil {
		return err
	}

	return nil
}

// Resets password using the token; validates token and expiration, hashes new password
func (u *AuthUsecase) ResetPassword(ctx context.Context, token string, newPassword string) error {
	// validation
	if validationErrors := validator.ValidateResetPasswordRequest(newPassword); len(validationErrors) > 0 {
		u.logger.Warn("Reset password validation failed",
			zap.Any("errors", validationErrors),
		)
		return &uerror.ValidationError{Errors: validationErrors}
	}

	user, err := u.userRepo.GetByPasswordResetToken(token)
	if err != nil {
		return uerror.ErrInvalidToken
	}

	if user.PasswordResetTokenExpiresAt == nil || time.Now().After(*user.PasswordResetTokenExpiresAt) {
		return uerror.ErrInvalidToken
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update user password and clear reset token
	if err := u.userRepo.UpdatePassword(user.ID, hashedPassword); err != nil {
		return err
	}

	return nil
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

func (u *AuthUsecase) sendVerificationEmail(ctx context.Context, user *domain.User) error {
	// Generate OTP
	token, err := utils.GenerateToken(64)
	fmt.Printf("%s %d\n", token, len(token))
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	// Set expiration
	expiresAt := time.Now().Add(time.Duration(u.cfg.Email.TokenExpirationMinutes) * time.Minute)

	// Save OTP to database
	if err := u.userRepo.UpdateVerificationToken(user.ID, token, expiresAt); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	// Send email
	if err := u.emailService.SendVerificationEmail(ctx, user.Email, user.Username, token); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Update last sent time
	if err := u.userRepo.UpdateLastSentAt(user.ID, time.Now()); err != nil {
		u.logger.Error("Failed to update last sent time", zap.Error(err))
	}

	return nil
}
