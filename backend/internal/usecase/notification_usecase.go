package usecase

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/pkg/redis"
	"go.uber.org/zap"
)

type NotificationUsecase struct {
	redis     *redis.RedisClient
	logger    *zap.Logger
	clients   map[int]chan domain.NotificationEvent
	clientsMu sync.RWMutex
}

func NewNotificationUsecase(redis *redis.RedisClient, logger *zap.Logger) *NotificationUsecase {
	u := &NotificationUsecase{
		redis:   redis,
		logger:  logger,
		clients: make(map[int]chan domain.NotificationEvent),
	}

	// Start listening to Redis Pub/Sub in a background goroutine
	go u.listenToEvents()

	return u
}

// AddClient registers a new SSE client for a user
func (u *NotificationUsecase) AddClient(userID int) chan domain.NotificationEvent {
	u.clientsMu.Lock()
	defer u.clientsMu.Unlock()

	// If client already exists, we might want to close the old one or support multiple
	// For simplicity, we'll support one connection per user for now
	ch := make(chan domain.NotificationEvent, 10)
	u.clients[userID] = ch

	u.logger.Info("New notification client registered", zap.Int("user_id", userID), zap.Int("active_clients", len(u.clients)))
	return ch
}

// RemoveClient unregisters an SSE client
func (u *NotificationUsecase) RemoveClient(userID int) {
	u.clientsMu.Lock()
	defer u.clientsMu.Unlock()

	if ch, ok := u.clients[userID]; ok {
		close(ch)
		delete(u.clients, userID)
		u.logger.Debug("Notification client removed", zap.Int("user_id", userID))
	}
}

// listenToEvents subscribes to Redis channels and broadcasts to relevant clients
func (u *NotificationUsecase) listenToEvents() {
	ctx := context.Background()
	pubsub := u.redis.Client.Subscribe(ctx, domain.AchievementEventChannel)
	defer pubsub.Close()

	ch := pubsub.Channel()

	u.logger.Info("Notification listener started", zap.String("channel", domain.AchievementEventChannel))

	for msg := range ch {
		u.logger.Info("Received message from Redis", zap.String("channel", msg.Channel), zap.String("payload", msg.Payload))
		var event domain.NotificationEvent
		if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
			u.logger.Error("Failed to unmarshal notification event", zap.Error(err))
			continue
		}

		u.broadcast(event)
	}
}

// broadcast sends the event to the correct user channel
func (u *NotificationUsecase) broadcast(event domain.NotificationEvent) {
	u.clientsMu.RLock()
	defer u.clientsMu.RUnlock()

	// Type-based routing
	switch event.Type {
	case domain.EventAchievementUnlocked:
		// Data is interface{}, need to extract userID
		// In a real app, we'd have a more robust way to handle this mapping
		// For now, we'll use a hack or re-unmarshal specifically if needed
		// Let's re-unmarshal to get the data reliably
		dataJson, _ := json.Marshal(event.Data)
		var data domain.AchievementUnlockedEvent
		json.Unmarshal(dataJson, &data)

		u.logger.Info("Attempting to broadcast achievement", zap.Int("target_user_id", data.UserID), zap.String("slug", data.Slug))

		if ch, ok := u.clients[data.UserID]; ok {
			select {
			case ch <- event:
				u.logger.Info("Notification sent to user channel", zap.Int("user_id", data.UserID), zap.String("type", event.Type))
			default:
				u.logger.Warn("Notification channel full for user", zap.Int("user_id", data.UserID))
			}
		} else {
			u.logger.Warn("No active notification client found for user", zap.Int("user_id", data.UserID), zap.Int("total_clients", len(u.clients)))
		}
	}
}
