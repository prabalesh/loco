package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/prabalesh/loco/backend/internal/delivery/cookies"
	"github.com/prabalesh/loco/backend/internal/delivery/middleware"
	"github.com/prabalesh/loco/backend/internal/domain/dto"
	"github.com/prabalesh/loco/backend/internal/domain/uerror"
	"github.com/prabalesh/loco/backend/internal/usecase"
	"github.com/prabalesh/loco/backend/pkg/config"
	"go.uber.org/zap"
)

type AdminAuthHandler struct {
	authUsecase   *usecase.AuthUsecase
	logger        *zap.Logger
	cfg           *config.Config
	cookieManager *cookies.CookieManager
}

func NewAdminAuthHandler(authUsecase *usecase.AuthUsecase, logger *zap.Logger, cfg *config.Config, cookieManager *cookies.CookieManager) *AdminAuthHandler {
	return &AdminAuthHandler{
		authUsecase:   authUsecase,
		logger:        logger,
		cfg:           cfg,
		cookieManager: cookieManager,
	}
}

// AdminLogin - Admin login endpoint with role verification
func (h *AdminAuthHandler) AdminLogin(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Invalid JSON in admin login request", zap.Error(err))
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Call login usecase
	user, tokenPair, err := h.authUsecase.Login(&req)
	if err != nil {
		var validationErr *uerror.ValidationError
		if errors.As(err, &validationErr) {
			RespondValidationError(w, validationErr.Errors)
			return
		}

		errMsg := err.Error()
		switch errMsg {
		case "invalid email or password":
			h.logger.Warn("Admin login failed: invalid credentials", zap.String("email", req.Email))
			RespondError(w, http.StatusUnauthorized, errMsg)
		case "account is deactivated":
			RespondError(w, http.StatusForbidden, errMsg)
		default:
			h.logger.Error("Admin login failed", zap.Error(err))
			RespondError(w, http.StatusInternalServerError, "login failed")
		}
		return
	}

	// Verify admin role
	if user.Role != "admin" {
		h.logger.Warn("Non-admin attempted admin login",
			zap.String("email", req.Email),
			zap.String("role", user.Role),
		)
		RespondError(w, http.StatusForbidden, "admin access required")
		return
	}

	// Verify email is verified
	if !user.EmailVerified {
		h.logger.Warn("Unverified admin attempted login", zap.String("email", req.Email))
		RespondError(w, http.StatusForbidden, "email not verified")
		return
	}

	// Set admin cookies with shorter expiry
	h.cookieManager.SetSecure(w, "adminAccessToken", tokenPair.AccessToken, 600)
	h.cookieManager.SetSecure(w, "adminRefreshToken", tokenPair.RefreshToken, 28800)

	h.logger.Info("Admin logged in successfully",
		zap.Int("admin_id", user.ID),
		zap.String("email", user.Email),
	)

	response := dto.LoginResponse{
		Message: "admin login successful",
		User:    dto.ToUserResponse(user),
	}

	RespondJSON(w, http.StatusOK, response)
}

// AdminRefreshToken - Refresh admin access token
func (h *AdminAuthHandler) AdminRefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("adminRefreshToken")
	if err != nil {
		h.logger.Warn("Admin refresh token cookie not found")
		RespondError(w, http.StatusUnauthorized, "refresh token required")
		return
	}

	// Generate new access token
	accessToken, _, err := h.authUsecase.RefreshAccessToken(cookie.Value)
	if err != nil {
		h.logger.Warn("Admin token refresh failed", zap.Error(err))
		RespondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Set new access token cookie (10 minutes)
	h.cookieManager.SetSecure(w, "adminAccessToken", accessToken, 28800)

	h.logger.Info("Admin access token refreshed successfully")

	RespondJSON(w, http.StatusOK, map[string]string{
		"message": "token refreshed successfully",
	})
}

// AdminLogout - Logout admin and clear cookies
func (h *AdminAuthHandler) AdminLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("adminRefreshToken")
	if err != nil {
		RespondJSON(w, http.StatusOK, map[string]string{
			"message": "logged out successfully",
		})
		return
	}

	// Revoke refresh token
	if err := h.authUsecase.Logout(cookie.Value); err != nil {
		h.logger.Error("Failed to revoke admin token", zap.Error(err))
	}

	// Clear cookies
	h.cookieManager.Clear(w, "adminAccessToken")
	h.cookieManager.Clear(w, "adminRefreshToken")

	h.logger.Info("Admin logged out successfully")

	RespondJSON(w, http.StatusOK, map[string]string{
		"message": "logged out successfully",
	})
}

// GetAdminProfile - Get current admin user info
func (h *AdminAuthHandler) GetAdminProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.authUsecase.GetCurrentUser(userID)
	if err != nil {
		h.logger.Error("Failed to get admin profile", zap.Error(err), zap.Int("user_id", userID))
		RespondError(w, http.StatusNotFound, "user not found")
		return
	}

	RespondJSON(w, http.StatusOK, dto.ToUserResponse(user))
}
