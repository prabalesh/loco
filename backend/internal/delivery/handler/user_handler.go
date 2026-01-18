package handler

import (
	"net/http"

	"github.com/prabalesh/loco/backend/internal/delivery/middleware"
	"github.com/prabalesh/loco/backend/internal/usecase"
	"go.uber.org/zap"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
	logger      *zap.Logger
}

func NewUserHandler(userUsecase *usecase.UserUsecase, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
		logger:      logger,
	}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	targetUserID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.userUsecase.GetUserProfile(targetUserID)
	if err != nil {
		h.logger.Warn("Failed to get user profile",
			zap.Error(err),
			zap.Int("user_id", targetUserID),
		)
		RespondError(w, http.StatusNotFound, "user not found")
		return
	}

	h.logger.Info("User profile retrieved",
		zap.Int("user_id", targetUserID),
	)

	RespondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) GetProfileByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	if username == "" {
		RespondError(w, http.StatusNotFound, "usern not found")
		return
	}

	user, err := h.userUsecase.GetUserProfileByUsername(username)
	if err != nil {
		h.logger.Warn("Failed to get user profile",
			zap.Error(err),
			zap.String("username", username),
		)
		RespondError(w, http.StatusNotFound, "user not found")
		return
	}

	h.logger.Info("User profile retrieved",
		zap.String("usernaeme", username),
	)

	RespondJSON(w, http.StatusOK, user)
}
