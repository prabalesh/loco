package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/prabalesh/loco/backend/pkg/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RateLimitMiddleware struct {
	redisClient *redis.Client
	logger      *zap.Logger
	cfg         *config.RateLimitConfig
}

func NewRateLimitMiddleware(redisClient *redis.Client, logger *zap.Logger, cfg *config.RateLimitConfig) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		redisClient: redisClient,
		logger:      logger,
		cfg:         cfg,
	}
}

func NewSubmissionRateLimitMiddleware(redisClient *redis.Client, logger *zap.Logger, cfg *config.SubmissionRateLimitConfig) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		redisClient: redisClient,
		logger:      logger,
		cfg:         &config.RateLimitConfig{Limit: cfg.Limit, Window: cfg.Window},
	}
}

func NewRunCodeRateLimitMiddleware(redisClient *redis.Client, logger *zap.Logger, cfg *config.RunCodeRateLimitConfig) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		redisClient: redisClient,
		logger:      logger,
		cfg:         &config.RateLimitConfig{Limit: cfg.Limit, Window: cfg.Window},
	}
}

func (m *RateLimitMiddleware) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserID(r.Context())
		if !ok {
			// If not authenticated, skip rate limiting or limit by IP (skipping for now as this is applied to protected routes)
			next.ServeHTTP(w, r)
			return
		}

		key := fmt.Sprintf("rate_limit:submission:%d", userID)
		limit := m.cfg.Limit
		window := time.Duration(m.cfg.Window) * time.Second

		ctx := context.Background()

		// Redis Incrementation
		count, err := m.redisClient.Incr(ctx, key).Result()
		if err != nil {
			m.logger.Error("Failed to increment rate limit counter", zap.Error(err))
			// Fail open
			next.ServeHTTP(w, r)
			return
		}

		// Set expiration on first request
		if count == 1 {
			m.redisClient.Expire(ctx, key, window)
		}

		// Check limit
		if count > int64(limit) {
			m.logger.Warn("Rate limit exceeded", zap.Int("user_id", userID))
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(window).Unix(), 10))
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
		w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(int64(limit)-count, 10))

		next.ServeHTTP(w, r)
	})
}
