package worker

import (
	"context"
	"time"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/infrastructure/queue"
	"github.com/prabalesh/loco/backend/internal/usecase"
	"go.uber.org/zap"
)

type AchievementWorker struct {
	queue              queue.JobQueue
	achievementUsecase *usecase.AchievementUsecase
	submissionRepo     domain.SubmissionRepository
	logger             *zap.Logger
	stopChan           chan struct{}
}

func NewAchievementWorker(
	queue queue.JobQueue,
	achievementUsecase *usecase.AchievementUsecase,
	submissionRepo domain.SubmissionRepository,
	logger *zap.Logger,
) *AchievementWorker {
	return &AchievementWorker{
		queue:              queue,
		achievementUsecase: achievementUsecase,
		submissionRepo:     submissionRepo,
		logger:             logger,
		stopChan:           make(chan struct{}),
	}
}

func (w *AchievementWorker) Start(ctx context.Context) {
	w.logger.Info("Achievement Worker started, waiting for jobs...")

	for {
		select {
		case <-w.stopChan:
			w.logger.Info("Achievement Worker stopped")
			return
		case <-ctx.Done():
			w.logger.Info("Achievement Worker context cancelled")
			return
		default:
			job, err := w.queue.DequeueAchievement(ctx)
			if err != nil {
				w.logger.Error("Failed to dequeue achievement job", zap.Error(err))
				time.Sleep(1 * time.Second)
				continue
			}

			if job == nil {
				continue
			}

			w.processAchievementJob(ctx, job.SubmissionID)
		}
	}
}

func (w *AchievementWorker) Stop() {
	close(w.stopChan)
}

func (w *AchievementWorker) processAchievementJob(ctx context.Context, submissionID int) {
	w.logger.Debug("Processing achievement job", zap.Int("submission_id", submissionID))

	// 1. Fetch submission
	submission, err := w.submissionRepo.GetByID(submissionID)
	if err != nil {
		w.logger.Error("Failed to fetch submission for achievement evaluation",
			zap.Error(err),
			zap.Int("submission_id", submissionID),
		)
		return
	}

	// Skip achievements for admin/validation submissions
	if submission.IsAdminSubmission || submission.IsValidationSubmission {
		w.logger.Debug("Skipping achievements for admin/validation submission", zap.Int("submission_id", submissionID))
		return
	}

	// 2. Fetch comprehensive stats for evaluation
	// We reuse the logic from the submission usecase for now, but encapsulated here
	totalSolved, err := w.submissionRepo.CountProblemsSolvedByUser(submission.UserID)
	if err != nil {
		w.logger.Error("Failed to count solved problems", zap.Error(err))
		return
	}

	totalAccepted, err := w.submissionRepo.CountAcceptedByUser(submission.UserID)
	if err != nil {
		w.logger.Error("Failed to count accepted submissions", zap.Error(err))
		return
	}

	streak, err := w.submissionRepo.GetCurrentStreak(submission.UserID)
	if err != nil {
		w.logger.Error("Failed to get current streak", zap.Error(err))
	}

	distribution, err := w.submissionRepo.GetSolvedDistribution(submission.UserID)
	if err != nil {
		w.logger.Error("Failed to get solved distribution", zap.Error(err))
	}

	stats := &domain.UserStats{
		AcceptedSubmissions: int(totalAccepted),
		ProblemsSolved:      int(totalSolved),
		Streak:              streak,
		SolvedDistribution:  distribution,
	}

	// 3. Evaluate
	if err := w.achievementUsecase.EvaluateSubmissionAchievements(submission, stats); err != nil {
		w.logger.Error("Failed to evaluate achievements",
			zap.Error(err),
			zap.Int("submission_id", submissionID),
			zap.Int("user_id", submission.UserID),
		)
	}

	w.logger.Info("Achievement evaluation completed",
		zap.Int("submission_id", submissionID),
		zap.Int("user_id", submission.UserID),
	)
}
