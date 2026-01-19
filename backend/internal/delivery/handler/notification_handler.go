package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/prabalesh/loco/backend/internal/delivery/middleware"
	"github.com/prabalesh/loco/backend/internal/usecase"
	"go.uber.org/zap"
)

type NotificationHandler struct {
	notificationUsecase *usecase.NotificationUsecase
	logger              *zap.Logger
}

func NewNotificationHandler(notificationUsecase *usecase.NotificationUsecase, logger *zap.Logger) *NotificationHandler {
	return &NotificationHandler{
		notificationUsecase: notificationUsecase,
		logger:              logger,
	}
}

func (h *NotificationHandler) Stream(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	// Create a channel for this user
	clientChan := h.notificationUsecase.AddClient(userID)
	defer h.notificationUsecase.RemoveClient(userID)

	// Keep connection open
	h.logger.Info("Checking flusher compatibility", zap.String("writer_type", fmt.Sprintf("%T", w)))
	flusher, ok := w.(http.Flusher)
	if !ok {
		h.logger.Error("Streaming unsupported: writer does not implement http.Flusher", zap.String("type", fmt.Sprintf("%T", w)))
		RespondError(w, http.StatusInternalServerError, "Streaming unsupported")
		return
	}

	// Send initial "connected" event
	h.logger.Info("Sending initial connected event to user", zap.Int("user_id", userID))
	fmt.Fprintf(w, "data: %s\n\n", "{\"type\":\"connected\"}")
	flusher.Flush()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	h.logger.Info("SSE connection opened", zap.Int("user_id", userID))

	for {
		select {
		case <-r.Context().Done():
			h.logger.Info("SSE connection closed by client", zap.Int("user_id", userID))
			return

		case event := <-clientChan:
			eventData, err := json.Marshal(event)
			if err != nil {
				h.logger.Error("Failed to marshal notification event", zap.Error(err))
				continue
			}

			fmt.Fprintf(w, "data: %s\n\n", eventData)
			flusher.Flush()

		case <-ticker.C:
			// Send keep-alive ping
			fmt.Fprintf(w, ": keep-alive\n\n")
			flusher.Flush()
		}
	}
}
