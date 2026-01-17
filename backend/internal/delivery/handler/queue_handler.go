package handler

import (
	"net/http"

	"github.com/prabalesh/loco/backend/internal/usecase"
	"go.uber.org/zap"
)

type QueueHandler struct {
	queueStatusUsecase *usecase.QueueStatusUsecase
	logger             *zap.Logger
}

func NewQueueHandler(queueStatusUsecase *usecase.QueueStatusUsecase, logger *zap.Logger) *QueueHandler {
	return &QueueHandler{
		queueStatusUsecase: queueStatusUsecase,
		logger:             logger,
	}
}

// GetQueueStatus returns the overall status of the submission queue
func (h *QueueHandler) GetQueueStatus(w http.ResponseWriter, r *http.Request) {
	queueStatus, err := h.queueStatusUsecase.GetQueueStatus()
	if err != nil {
		h.logger.Error("Failed to get queue status", zap.Error(err))
		RespondError(w, http.StatusInternalServerError, "failed to get queue status")
		return
	}

	RespondJSON(w, http.StatusOK, queueStatus)
}
