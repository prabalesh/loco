package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/prabalesh/loco/backend/pkg/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisClient struct {
	Client *redis.Client
	Logger *zap.Logger
}

func NewRedisClient(cfg config.RedisConfig, logger *zap.Logger) (*RedisClient, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	opts := &redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		// Connection pool settings
		PoolSize:     10,
		MinIdleConns: 5,
		// Timeouts
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		// Connection pool timeouts
		PoolTimeout: 4 * time.Second,
		// Connection lifecycle
		ConnMaxIdleTime: 5 * time.Minute,
		ConnMaxLifetime: 30 * time.Minute,
	}

	client := redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Error("Failed to connect to Redis", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Redis connected successfully",
		zap.String("host", cfg.Host),
		zap.String("port", cfg.Port),
		zap.Int("db", cfg.DB),
	)

	return &RedisClient{
		Client: client,
		Logger: logger,
	}, nil
}

func (r *RedisClient) Close() error {
	if r.Client != nil {
		return r.Client.Close()
	}
	return nil
}
