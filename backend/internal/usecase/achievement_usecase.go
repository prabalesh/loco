package usecase

import (
	"fmt"

	"github.com/prabalesh/loco/backend/internal/domain"
	"go.uber.org/zap"
)

type AchievementUsecase struct {
	achievementRepo domain.AchievementRepository
	userRepo        domain.UserRepository
	logger          *zap.Logger
}

func NewAchievementUsecase(
	achievementRepo domain.AchievementRepository,
	userRepo domain.UserRepository,
	logger *zap.Logger,
) *AchievementUsecase {
	return &AchievementUsecase{
		achievementRepo: achievementRepo,
		userRepo:        userRepo,
		logger:          logger,
	}
}

// CheckAndUnlock checks if a specific achievement condition is met and unlocks it
func (u *AchievementUsecase) CheckAndUnlock(userID int, slug string) error {
	// 1. Check if already unlocked
	achievement, err := u.achievementRepo.GetBySlug(slug)
	if err != nil {
		return fmt.Errorf("achievement not found: %s", slug)
	}

	unlocked, err := u.achievementRepo.HasUnlocked(userID, achievement.ID)
	if err != nil {
		return err
	}
	if unlocked {
		return nil // Already unlocked
	}

	// 2. Unlock
	if err := u.achievementRepo.Unlock(userID, achievement.ID); err != nil {
		return err
	}

	// 3. Award XP
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	user.XP += achievement.XPReward
	// Simple level formula: Level = 1 + (XP / 100)
	user.Level = 1 + (user.XP / 100)

	if err := u.userRepo.Update(user); err != nil {
		return err
	}

	u.logger.Info("Achievement unlocked",
		zap.Int("user_id", userID),
		zap.String("slug", slug),
		zap.Int("xp_awarded", achievement.XPReward),
	)

	return nil
}

// EvaluateSubmissionAchievements checks for achievements related to a submission
func (u *AchievementUsecase) EvaluateSubmissionAchievements(submission *domain.Submission, stats *domain.UserStats) error {
	userID := submission.UserID

	// 1. Hello World (First Submission)
	_ = u.CheckAndUnlock(userID, "hello-world")

	// 2. Conditions based on Status
	if submission.Status == domain.SubmissionStatusAccepted {
		if err := u.CheckAndUnlock(userID, "first-blood"); err != nil {
			u.logger.Error("Failed to unlock first-blood", zap.Error(err))
		}

		// One Shot (First attempt is AC)
		if submission.TotalTestCases > 0 && submission.PassedTestCases == submission.TotalTestCases {
			// Logic for One Shot would require checking previous submissions for this problem
		}

		// Count based achievements
		if stats.ProblemsSolved >= 1 {
			_ = u.CheckAndUnlock(userID, "solver-i")
		}
		if stats.ProblemsSolved >= 10 {
			_ = u.CheckAndUnlock(userID, "solver-ii")
		}
		// ... maps to other counts

	} else if submission.Status == domain.SubmissionStatusWrongAnswer {
		_ = u.CheckAndUnlock(userID, "bug-hunter")
	} else if submission.Status == domain.SubmissionStatusTimeLimitExceeded {
		_ = u.CheckAndUnlock(userID, "speed-demon")
	} else if submission.Status == domain.SubmissionStatusMemoryLimitExceeded {
		_ = u.CheckAndUnlock(userID, "memory-leak")
	}

	return nil
}

func (u *AchievementUsecase) ListAll() ([]domain.Achievement, error) {
	return u.achievementRepo.GetAll()
}

func (u *AchievementUsecase) GetUserProgress(userID int) ([]domain.UserAchievement, error) {
	return u.achievementRepo.GetUnlockedByUser(userID)
}
