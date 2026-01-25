package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type CacheService interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	DeleteByPrefix(ctx context.Context, prefix string) error
}

type cacheService struct {
	redis  *redis.Client
	logger *zap.Logger
}

func NewCacheService(redisClient *redis.Client, logger *zap.Logger) CacheService {
	return &cacheService{
		redis:  redisClient,
		logger: logger,
	}
}

func (s *cacheService) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := s.redis.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		s.logger.Error("Failed to get from cache", zap.String("key", key), zap.Error(err))
		return err
	}

	return json.Unmarshal(data, dest)
}

func (s *cacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		s.logger.Error("Failed to marshal cache value", zap.String("key", key), zap.Error(err))
		return err
	}

	if err := s.redis.Set(ctx, key, data, expiration).Err(); err != nil {
		s.logger.Error("Failed to set cache", zap.String("key", key), zap.Error(err))
		return err
	}

	return nil
}

func (s *cacheService) Delete(ctx context.Context, key string) error {
	if err := s.redis.Del(ctx, key).Err(); err != nil {
		s.logger.Error("Failed to delete from cache", zap.String("key", key), zap.Error(err))
		return err
	}
	return nil
}

func (s *cacheService) DeleteByPrefix(ctx context.Context, prefix string) error {
	iter := s.redis.Scan(ctx, 0, prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		if err := s.redis.Del(ctx, iter.Val()).Err(); err != nil {
			s.logger.Error("Failed to delete by prefix", zap.String("prefix", prefix), zap.Error(err))
		}
	}
	if err := iter.Err(); err != nil {
		s.logger.Error("Scan error during DeleteByPrefix", zap.String("prefix", prefix), zap.Error(err))
		return err
	}
	return nil
}
