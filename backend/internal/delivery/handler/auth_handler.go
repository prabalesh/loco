package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/prabalesh/loco/backend/internal/delivery/middleware"
	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/usecase"
	"github.com/prabalesh/loco/backend/pkg/config"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
	logger      *zap.Logger
	cfg         *config.Config
}

func NewAuthHandler(authUsecase *usecase.AuthUsecase, logger *zap.Logger, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
		logger:      logger,
		cfg:         cfg,
	}
}

// Register handles user registration
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
		Message: "registration successful",
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

		switch errMsg {
		case "invalid email or password":
			h.logger.Warn("Login failed: invalid email or password", zap.String("error", errMsg))
			RespondError(w, http.StatusUnauthorized, errMsg)
		case "account is deactivated":
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

	h.setTokenCookie(w, "accessToken", tokenPair.AccessToken, int(tokenPair.AccessExpiresAt.Seconds()))
	h.setTokenCookie(w, "refreshToken", tokenPair.RefreshToken, int(tokenPair.RefreshExpiresAt.Seconds()))

	RespondJSON(w, http.StatusCreated, response)
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
	h.setTokenCookie(w, "accessToken", accessToken, maxAge)

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
	h.clearTokenCookie(w, "accessToken")
	h.clearTokenCookie(w, "refreshToken")

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

// setTokenCookie sets an HTTP-only cookie with environment-aware settings
func (h *AuthHandler) setTokenCookie(w http.ResponseWriter, name, value string, maxAge int) {
	// Parse SameSite from config
	sameSite := http.SameSiteLaxMode
	switch h.cfg.Cookie.SameSite {
	case "strict":
		sameSite = http.SameSiteStrictMode
	case "none":
		sameSite = http.SameSiteNoneMode
	case "lax":
		sameSite = http.SameSiteLaxMode
	}

	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   h.cfg.Cookie.Secure,
		SameSite: sameSite,
		Domain:   h.cfg.Cookie.Domain,
	}
	http.SetCookie(w, cookie)
}

// clearTokenCookie removes a cookie
func (h *AuthHandler) clearTokenCookie(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.cfg.Cookie.Secure,
		Domain:   h.cfg.Cookie.Domain,
	}
	http.SetCookie(w, cookie)
}
