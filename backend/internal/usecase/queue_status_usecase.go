package usecase

import (
	"context"
	"time"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/infrastructure/queue"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type QueueStatusUsecase struct {
	submissionRepo domain.SubmissionRepository
	redisClient    *redis.Client
	logger         *zap.Logger
}

func NewQueueStatusUsecase(
	submissionRepo domain.SubmissionRepository,
	redisClient *redis.Client,
	logger *zap.Logger,
) *QueueStatusUsecase {
	return &QueueStatusUsecase{
		submissionRepo: submissionRepo,
		redisClient:    redisClient,
		logger:         logger,
	}
}

// GetQueueStatus returns the overall status of the submission queue
func (u *QueueStatusUsecase) GetQueueStatus() (*domain.QueueStatus, error) {
	ctx := context.Background()

	// Get queue size from Redis
	queueSize, err := u.redisClient.LLen(ctx, queue.SubmissionQueueName).Result()
	if err != nil {
		u.logger.Error("Failed to get queue size", zap.Error(err))
		queueSize = 0
	}

	// Get pending submissions count from database
	pendingCount, err := u.submissionRepo.CountPending()
	if err != nil {
		u.logger.Error("Failed to count pending submissions", zap.Error(err))
		pendingCount = 0
	}

	// Get oldest pending submission to calculate age
	oldestAge := int64(0)
	if pendingCount > 0 {
		oldestSubmission, err := u.getOldestPendingSubmission()
		if err == nil && oldestSubmission != nil {
			oldestAge = int64(time.Since(oldestSubmission.CreatedAt).Seconds())
		}
	}

	// Determine worker count (simplified - checking if any jobs are being processed)
	// In a real system, you'd have worker heartbeats
	activeWorkers := u.estimateActiveWorkers(pendingCount, queueSize)

	// Determine health status
	healthStatus, warningMsg := u.determineHealthStatus(activeWorkers, queueSize, oldestAge)

	// Estimate wait time (rough estimate: 30 seconds per submission if workers are active)
	estimatedWaitTime := int64(0)
	if activeWorkers > 0 && queueSize > 0 {
		estimatedWaitTime = (queueSize * 30) / int64(activeWorkers)
	}

	return &domain.QueueStatus{
		QueueSize:         queueSize,
		ActiveWorkers:     activeWorkers,
		OldestPendingAge:  oldestAge,
		HealthStatus:      healthStatus,
		WarningMessage:    warningMsg,
		PendingCount:      pendingCount,
		EstimatedWaitTime: estimatedWaitTime,
	}, nil
}

// GetSubmissionQueueInfo returns queue-specific info for a submission
func (u *QueueStatusUsecase) GetSubmissionQueueInfo(submissionID int) (*domain.SubmissionQueueInfo, error) {
	// Get submission
	submission, err := u.submissionRepo.GetByID(submissionID)
	if err != nil {
		return nil, err
	}

	// Only calculate queue info for pending submissions
	if submission.Status != domain.SubmissionStatusPending {
		return &domain.SubmissionQueueInfo{
			WorkersActive: true, // Assume workers are active if submission was processed
		}, nil
	}

	// Get queue status
	queueStatus, err := u.GetQueueStatus()
	if err != nil {
		u.logger.Error("Failed to get queue status", zap.Error(err))
		return &domain.SubmissionQueueInfo{
			WorkersActive: false,
		}, nil
	}

	// Calculate approximate queue position based on creation time
	position := u.calculateQueuePosition(submission)

	// Estimate wait time for this specific submission
	estimatedWait := int64(0)
	if queueStatus.ActiveWorkers > 0 && position > 0 {
		estimatedWait = (int64(position) * 30) / int64(queueStatus.ActiveWorkers)
	}

	return &domain.SubmissionQueueInfo{
		QueuePosition:     position,
		EstimatedWaitTime: estimatedWait,
		WorkersActive:     queueStatus.ActiveWorkers > 0,
	}, nil
}

// Helper methods

func (u *QueueStatusUsecase) getOldestPendingSubmission() (*domain.Submission, error) {
	// Get the oldest pending submission
	submissions, err := u.submissionRepo.GetOldestPending(1)
	if err != nil || len(submissions) == 0 {
		return nil, err
	}
	return &submissions[0], nil
}

func (u *QueueStatusUsecase) estimateActiveWorkers(pendingCount, queueSize int64) int {
	// Simple heuristic: if there are pending submissions but queue is empty or small,
	// workers might be processing them
	// If queue is growing and old submissions exist, workers are likely inactive

	// For now, we'll check if queue size matches pending count
	// If they're very different, workers might be processing
	if pendingCount == 0 {
		return 0
	}

	// If queue size is much smaller than pending count, workers are likely active
	if queueSize < pendingCount/2 {
		return 1 // Estimate at least 1 worker is active
	}

	// If queue size equals pending count, no workers are processing
	if queueSize >= pendingCount {
		return 0
	}

	return 1
}

func (u *QueueStatusUsecase) determineHealthStatus(workers int, queueSize, oldestAge int64) (domain.QueueHealthStatus, string) {
	// Critical: No workers active
	if workers == 0 && queueSize > 0 {
		return domain.QueueHealthCritical, "No workers are currently active. Submissions will not be processed."
	}

	// Critical: Submissions waiting too long (> 5 minutes)
	if oldestAge > 300 {
		return domain.QueueHealthCritical, "Submissions have been pending for over 5 minutes. Please check worker status."
	}

	// Warning: Queue is backing up (> 10 submissions)
	if queueSize > 10 {
		return domain.QueueHealthWarning, "Queue is backed up with multiple pending submissions."
	}

	// Warning: Submissions waiting moderately long (> 2 minutes)
	if oldestAge > 120 {
		return domain.QueueHealthWarning, "Some submissions have been pending for over 2 minutes."
	}

	return domain.QueueHealthHealthy, ""
}

func (u *QueueStatusUsecase) calculateQueuePosition(submission *domain.Submission) int {
	// Count how many pending submissions were created before this one
	// This is a simplified approach - in production you'd want to track actual queue positions
	// For now, we'll estimate based on creation time
	olderPending, err := u.submissionRepo.CountPendingBefore(submission.CreatedAt)
	if err != nil {
		u.logger.Error("Failed to count older pending submissions", zap.Error(err))
		return 0
	}

	return int(olderPending) + 1 // +1 because position is 1-indexed
}
