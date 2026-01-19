package worker

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/prabalesh/loco/backend/internal/infrastructure/queue"
	"github.com/prabalesh/loco/backend/pkg/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type mockQueue struct{}

func (m *mockQueue) EnqueueSubmission(ctx context.Context, submissionID int) error {
	return nil
}

func (m *mockQueue) DequeueSubmission(ctx context.Context) (*queue.SubmissionJob, error) {
	<-ctx.Done()
	return nil, nil
}

func (m *mockQueue) EnqueueAchievement(ctx context.Context, submissionID int) error {
	return nil
}

func (m *mockQueue) DequeueAchievement(ctx context.Context) (*queue.AchievementJob, error) {
	<-ctx.Done()
	return nil, nil
}

func TestWorkerHeartbeat(t *testing.T) {
	// Start miniredis
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to create miniredis: %v", err)
	}
	defer s.Close()

	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	// Create logger
	logger := zap.NewNop()

	// Shorten heartbeat interval for test
	oldInterval := HeartbeatInterval
	HeartbeatInterval = 100 * time.Millisecond
	defer func() { HeartbeatInterval = oldInterval }()

	// Initialize Worker with minimal dependencies
	w := NewWorker(
		&mockQueue{}, // queue
		nil,          // submissionRepo
		nil,          // problemRepo
		nil,          // testCaseRepo
		nil,          // languageRepo
		nil,          // problemLanguageRepo
		nil,          // pistonService
		nil,          // userStatsRepo
		logger,
		rdb,
		&config.Config{
			Worker: config.WorkerConfig{
				MaxConcurrentSubmissions: 4,
				MaxConcurrentTestCases:   5,
			},
		},
	)

	// Context for test
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start worker in goroutine
	go w.Start(ctx)

	// Wait for heartbeat
	time.Sleep(200 * time.Millisecond)

	// Verify key exists
	keys, err := rdb.Keys(context.Background(), "worker:*:heartbeat").Result()
	if err != nil {
		t.Fatalf("Failed to list keys: %v", err)
	}

	if len(keys) != 1 {
		t.Errorf("Expected 1 heartbeat key, got %d", len(keys))
	}

	// Stop worker
	w.Stop()
	cancel()
	// Wait specifically for the cleanup or just wait a bit
	time.Sleep(100 * time.Millisecond)

	// Verify key is gone
	keys, err = rdb.Keys(context.Background(), "worker:*:heartbeat").Result()
	if err != nil {
		t.Fatalf("Failed to list keys: %v", err)
	}

	if len(keys) != 0 {
		t.Errorf("Expected 0 heartbeat keys after stop, got %d", len(keys))
	}
}
