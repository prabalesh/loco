package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/prabalesh/loco/backend/internal/delivery/cookies"
	"github.com/prabalesh/loco/backend/internal/delivery/middleware"
	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/usecase"
	"github.com/prabalesh/loco/backend/pkg/config"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authUsecase   *usecase.AuthUsecase
	logger        *zap.Logger
	cfg           *config.Config
	cookieManager *cookies.CookieManager
}

func NewAuthHandler(authUsecase *usecase.AuthUsecase, logger *zap.Logger, cfg *config.Config, cookieManager *cookies.CookieManager) *AuthHandler {
	return &AuthHandler{
		authUsecase:   authUsecase,
		logger:        logger,
		cfg:           cfg,
		cookieManager: cookieManager,
	}
}

// register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Invalid JSON in register request", zap.Error(err))
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Call use case
	user, err := h.authUsecase.Register(&req)
	if err != nil {
		// Handle validation errors
		var validationErr *usecase.ValidationError
		if errors.As(err, &validationErr) {
			h.logger.Warn("Registration validation failed",
				zap.Any("errors", validationErr.Errors),
			)
			RespondValidationError(w, validationErr.Errors)
			return
		}

		// Handle business logic errors
		errMsg := err.Error()

		switch errMsg {
		case "email already registered":
			h.logger.Warn("Registration failed: duplicate email", zap.String("error", errMsg))
			RespondError(w, http.StatusConflict, errMsg)
		case "username already taken":
			h.logger.Warn("Registration failed: duplicate username", zap.String("error", errMsg))
			RespondError(w, http.StatusConflict, errMsg)
		default:
			h.logger.Error("Registration failed with unexpected error", zap.Error(err))
			RespondError(w, http.StatusInternalServerError, "failed to create account")
		}
		return
	}

	// Return success response
	h.logger.Info("User registered successfully",
		zap.Int("user_id", user.ID),
		zap.String("email", user.Email),
	)

	response := domain.RegisterResponse{
		Message: "registration successful. we have send an email for verification in your email id. please verify",
		User:    user.ToResponse(),
	}

	RespondJSON(w, http.StatusCreated, response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// parse request
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Invalid JSON in login request", zap.Error(err))
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// handle usecase
	user, tokenPair, err := h.authUsecase.Login(&req)
	if err != nil {
		// handle validation error
		var validationErr *usecase.ValidationError
		if errors.As(err, &validationErr) {
			h.logger.Warn("Registration validation failed",
				zap.Any("errors", validationErr.Errors),
			)
			RespondValidationError(w, validationErr.Errors)
			return
		}

		// handling business logic errors
		errMsg := err.Error()

		switch {
		case err == usecase.ErrEmailNotVerified:
			h.logger.Warn("Login failed: email not verified", zap.String("email", req.Email))
			RespondError(w, http.StatusForbidden, "please verify your email before logging in")
		case errMsg == "invalid email or password":
			h.logger.Warn("Login failed: invalid email or password", zap.String("error", errMsg))
			RespondError(w, http.StatusUnauthorized, errMsg)
		case errMsg == "account is deactivated":
			h.logger.Warn("Login failed: account is deactivated", zap.String("error", errMsg))
			RespondError(w, http.StatusForbidden, errMsg)
		default:
			h.logger.Error("Login failed with unexpected error", zap.Error(err))
			RespondError(w, http.StatusInternalServerError, "failed to create account")
		}
		return
	}

	// return success response
	h.logger.Info("User log in successfully",
		zap.Int("user_id", user.ID),
		zap.String("email", user.Email),
	)

	response := domain.RegisterResponse{
		Message: "login successful",
		User:    user.ToResponse(),
	}

	h.cookieManager.SetSecure(w, "accessToken", tokenPair.AccessToken, int(tokenPair.AccessExpiresAt.Seconds()))
	h.cookieManager.SetSecure(w, "refreshToken", tokenPair.RefreshToken, int(tokenPair.RefreshExpiresAt.Seconds()))

	RespondJSON(w, http.StatusOK, response)
}

// RefreshToken generates new access token from refresh token cookie
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Get refresh token from cookie
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		h.logger.Warn("Refresh token cookie not found")
		RespondError(w, http.StatusUnauthorized, "refresh token required")
		return
	}

	// Generate new access token
	accessToken, expiresAt, err := h.authUsecase.RefreshAccessToken(cookie.Value)
	if err != nil {
		h.logger.Warn("Token refresh failed", zap.Error(err))
		RespondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Set new access token cookie
	maxAge := int(expiresAt.Seconds())
	h.cookieManager.SetSecure(w, "accessToken", accessToken, maxAge)

	h.logger.Info("Access token refreshed successfully")

	RespondJSON(w, http.StatusOK, map[string]string{
		"message": "token refreshed successfully",
	})
}

// Logout revokes refresh token and clears cookies
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get refresh token from cookie
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		// No cookie = already logged out
		RespondJSON(w, http.StatusOK, map[string]string{
			"message": "logged out successfully",
		})
		return
	}

	// Revoke refresh token in database
	if err := h.authUsecase.Logout(cookie.Value); err != nil {
		h.logger.Error("Failed to revoke token", zap.Error(err))
	}

	// Clear cookies
	h.cookieManager.Clear(w, "accessToken")
	h.cookieManager.Clear(w, "refreshToken")

	h.logger.Info("User logged out successfully")

	RespondJSON(w, http.StatusOK, map[string]string{
		"message": "logged out successfully",
	})
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.authUsecase.GetCurrentUser(userID)
	if err != nil {
		h.logger.Error("Failed to get user", zap.Error(err), zap.Int("user_id", userID))
		RespondError(w, http.StatusNotFound, "user not found")
		return
	}

	h.logger.Info("User retrieved successfully", zap.Int("user_id", userID))
	RespondJSON(w, http.StatusOK, user.ToResponse())
}

func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req domain.VerifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Invalid JSON in verify email request", zap.Error(err))
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.authUsecase.VerifyEmail(r.Context(), &req); err != nil {
		switch err {
		case usecase.ErrInvalidToken:
			h.logger.Warn("Invalid token", zap.String("token", req.Token))
			RespondError(w, http.StatusBadRequest, "invalid or expired verification token")
		default:
			h.logger.Error("Failed to verify email", zap.Error(err))
			RespondError(w, http.StatusInternalServerError, "failed to verify email")
		}
		return
	}

	h.logger.Info("Email verified successfully")
	RespondJSON(w, http.StatusOK, map[string]string{
		"message": "email verified successfully",
	})
}

// resend verification email
func (h *AuthHandler) ResendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	var req domain.ResendVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Invalid JSON in resend verification request", zap.Error(err))
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.authUsecase.ResendVerificationEmail(r.Context(), &req); err != nil {
		switch {
		case errors.Is(err, usecase.ErrResendCooldown):
			h.logger.Warn("Resend cooldown active", zap.String("email", req.Email))
			RespondError(w, http.StatusTooManyRequests, err.Error())
		case err == usecase.ErrMaxTokenAttemptsExceeded:
			h.logger.Warn("Max attempts exceeded", zap.String("email", req.Email))
			RespondError(w, http.StatusTooManyRequests, "maximum attempts exceeded")
		default:
			h.logger.Error("Failed to resend verification email", zap.Error(err))
			RespondError(w, http.StatusInternalServerError, "failed to send verification email")
		}
		return
	}

	h.logger.Info("Verification email resent", zap.String("email", req.Email))
	RespondJSON(w, http.StatusOK, map[string]string{
		"message": "verification email sent successfully",
	})
}

// POST /auth/forgot-password
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email string `json:"email"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		RespondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	err := h.authUsecase.ForgotPassword(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, usecase.ErrResendCooldown) {
			RespondError(w, http.StatusTooManyRequests, "Please wait before requesting another reset email")
			return
		}
		h.logger.Error("Failed to process forgot password", zap.Error(err))
		// For security, respond with generic message
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "If the email exists, a reset link has been sent"})
}

// POST /auth/reset-password
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Token == "" || req.NewPassword == "" {
		RespondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if err := h.authUsecase.ResetPassword(r.Context(), req.Token, req.NewPassword); err != nil {
		switch err {
		case usecase.ErrInvalidToken:
			RespondError(w, http.StatusBadRequest, "invalid password reset token")
		default:
			RespondError(w, http.StatusInternalServerError, "failed to reset password")
		}
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "password reset successful"})
}
