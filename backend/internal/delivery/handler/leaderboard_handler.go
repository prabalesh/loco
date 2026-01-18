package handler

import (
	"net/http"
	"strconv"

	"github.com/prabalesh/loco/backend/internal/usecase"
	"go.uber.org/zap"
)

type LeaderboardHandler struct {
	leaderboardUsecase *usecase.LeaderboardUsecase
	logger             *zap.Logger
}

func NewLeaderboardHandler(leaderboardUsecase *usecase.LeaderboardUsecase, logger *zap.Logger) *LeaderboardHandler {
	return &LeaderboardHandler{
		leaderboardUsecase: leaderboardUsecase,
		logger:             logger,
	}
}

func (h *LeaderboardHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	entries, err := h.leaderboardUsecase.GetLeaderboard(limit)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, entries)
}
