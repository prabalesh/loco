package handler

import (
	"encoding/json"
	"net/http"

	"github.com/prabalesh/loco/backend/internal/usecase"
	"go.uber.org/zap"
)

type AchievementHandler struct {
	achievementUsecase *usecase.AchievementUsecase
	userUsecase        *usecase.UserUsecase
	logger             *zap.Logger
}

func NewAchievementHandler(
	achievementUsecase *usecase.AchievementUsecase,
	userUsecase *usecase.UserUsecase,
	logger *zap.Logger,
) *AchievementHandler {
	return &AchievementHandler{
		achievementUsecase: achievementUsecase,
		userUsecase:        userUsecase,
		logger:             logger,
	}
}

// List returns all available achievements
func (h *AchievementHandler) List(w http.ResponseWriter, r *http.Request) {
	achievements, err := h.achievementUsecase.ListAll()
	if err != nil {
		h.logger.Error("Failed to list achievements", zap.Error(err))
		http.Error(w, "Failed to list achievements", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": achievements,
	})
}

// GetMyAchievements returns achievements unlocked by the current user
func (h *AchievementHandler) GetMyAchievements(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	userAchievements, err := h.achievementUsecase.GetUserProgress(userID)
	if err != nil {
		h.logger.Error("Failed to get user achievements", zap.Error(err))
		http.Error(w, "Failed to get user achievements", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": userAchievements,
	})
}

// GetUserAchievements returns achievements for a specific user (public)
func (h *AchievementHandler) GetUserAchievements(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	user, err := h.userUsecase.GetByUsername(username)
	if err != nil {
		h.logger.Error("User not found", zap.Error(err))
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	userAchievements, err := h.achievementUsecase.GetUserProgress(user.ID)
	if err != nil {
		h.logger.Error("Failed to get user achievements", zap.Error(err))
		http.Error(w, "Failed to get user achievements", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": userAchievements,
	})
}
