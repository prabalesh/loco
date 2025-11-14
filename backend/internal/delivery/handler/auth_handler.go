package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/usecase"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
	logger      *zap.Logger
}

func NewAuthHandler(authUsecase *usecase.AuthUsecase, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
		logger:      logger,
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
