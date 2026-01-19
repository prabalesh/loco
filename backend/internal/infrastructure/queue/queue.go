package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/prabalesh/loco/backend/pkg/redis"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	SubmissionQueueName  = "submission:queue"
	AchievementQueueName = "achievement:queue"
	QueueTimeout         = 0 // 0 means block indefinitely
)

type JobQueue interface {
	EnqueueSubmission(ctx context.Context, submissionID int) error
	DequeueSubmission(ctx context.Context) (*SubmissionJob, error)
	EnqueueAchievement(ctx context.Context, submissionID int) error
	DequeueAchievement(ctx context.Context) (*AchievementJob, error)
}

type SubmissionJob struct {
	SubmissionID int       `json:"submission_id"`
	EnqueuedAt   time.Time `json:"enqueued_at"`
}

type AchievementJob struct {
	SubmissionID int       `json:"submission_id"`
	EnqueuedAt   time.Time `json:"enqueued_at"`
}

type jobQueue struct {
	redis  *redis.RedisClient
	logger *zap.Logger
}

func NewJobQueue(redisClient *redis.RedisClient, logger *zap.Logger) JobQueue {
	return &jobQueue{
		redis:  redisClient,
		logger: logger,
	}
}

// EnqueueSubmission pushes a submission job to the Redis queue
func (q *jobQueue) EnqueueSubmission(ctx context.Context, submissionID int) error {
	job := SubmissionJob{
		SubmissionID: submissionID,
		EnqueuedAt:   time.Now(),
	}

	jobData, err := json.Marshal(job)
	if err != nil {
		q.logger.Error("Failed to marshal submission job", zap.Error(err), zap.Int("submission_id", submissionID))
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	if err := q.redis.Client.LPush(ctx, SubmissionQueueName, jobData).Err(); err != nil {
		q.logger.Error("Failed to enqueue submission", zap.Error(err), zap.Int("submission_id", submissionID))
		return fmt.Errorf("failed to enqueue submission: %w", err)
	}

	q.logger.Info("Submission enqueued successfully",
		zap.Int("submission_id", submissionID),
		zap.String("queue", SubmissionQueueName),
	)

	return nil
}

// DequeueSubmission pulls a submission job from the Redis queue using BRPOP (blocking)
func (q *jobQueue) DequeueSubmission(ctx context.Context) (*SubmissionJob, error) {
	// BRPOP blocks until an item is available or timeout
	result, err := q.redis.Client.BRPop(ctx, time.Duration(QueueTimeout)*time.Second, SubmissionQueueName).Result()
	if err != nil {
		if err == goredis.Nil {
			// Timeout or no items - this is expected behavior
			return nil, nil
		}
		q.logger.Error("Failed to dequeue submission", zap.Error(err))
		return nil, fmt.Errorf("failed to dequeue submission: %w", err)
	}

	if len(result) < 2 {
		return nil, fmt.Errorf("invalid queue result format")
	}

	var job SubmissionJob
	if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
		q.logger.Error("Failed to unmarshal submission job", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}

	q.logger.Debug("Submission dequeued successfully",
		zap.Int("submission_id", job.SubmissionID),
	)

	return &job, nil
}

// EnqueueAchievement pushes an achievement evaluation job to the Redis queue
func (q *jobQueue) EnqueueAchievement(ctx context.Context, submissionID int) error {
	job := AchievementJob{
		SubmissionID: submissionID,
		EnqueuedAt:   time.Now(),
	}

	jobData, err := json.Marshal(job)
	if err != nil {
		q.logger.Error("Failed to marshal achievement job", zap.Error(err), zap.Int("submission_id", submissionID))
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	if err := q.redis.Client.LPush(ctx, AchievementQueueName, jobData).Err(); err != nil {
		q.logger.Error("Failed to enqueue achievement job", zap.Error(err), zap.Int("submission_id", submissionID))
		return fmt.Errorf("failed to enqueue achievement: %w", err)
	}

	q.logger.Info("Achievement job enqueued successfully",
		zap.Int("submission_id", submissionID),
		zap.String("queue", AchievementQueueName),
	)

	return nil
}

// DequeueAchievement pulls an achievement job from the Redis queue using BRPOP (blocking)
func (q *jobQueue) DequeueAchievement(ctx context.Context) (*AchievementJob, error) {
	result, err := q.redis.Client.BRPop(ctx, time.Duration(QueueTimeout)*time.Second, AchievementQueueName).Result()
	if err != nil {
		if err == goredis.Nil {
			return nil, nil
		}
		q.logger.Error("Failed to dequeue achievement job", zap.Error(err))
		return nil, fmt.Errorf("failed to dequeue achievement: %w", err)
	}

	if len(result) < 2 {
		return nil, fmt.Errorf("invalid queue result format")
	}

	var job AchievementJob
	if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
		q.logger.Error("Failed to unmarshal achievement job", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}

	q.logger.Debug("Achievement job dequeued successfully",
		zap.Int("submission_id", job.SubmissionID),
	)

	return &job, nil
}
