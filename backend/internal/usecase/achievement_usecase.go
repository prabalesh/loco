package usecase

import (
	"fmt"

	"github.com/prabalesh/loco/backend/internal/domain"
	"go.uber.org/zap"
)

type AchievementUsecase struct {
	achievementRepo domain.AchievementRepository
	userRepo        domain.UserRepository
	submissionRepo  domain.SubmissionRepository
	problemRepo     domain.ProblemRepository
	logger          *zap.Logger
}

func NewAchievementUsecase(
	achievementRepo domain.AchievementRepository,
	userRepo domain.UserRepository,
	submissionRepo domain.SubmissionRepository,
	problemRepo domain.ProblemRepository,
	logger *zap.Logger,
) *AchievementUsecase {
	return &AchievementUsecase{
		achievementRepo: achievementRepo,
		userRepo:        userRepo,
		submissionRepo:  submissionRepo,
		problemRepo:     problemRepo,
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

		// Problem Solving counts (Solver I-X)
		solverSlugs := []string{"solver-i", "solver-ii", "solver-iii", "solver-iv", "solver-v", "solver-vi", "solver-vii", "solver-viii", "solver-ix", "solver-x"}
		solverCounts := []int{1, 10, 25, 50, 100, 200, 300, 400, 500, 1000}
		for i, count := range solverCounts {
			if stats.ProblemsSolved >= count {
				_ = u.CheckAndUnlock(userID, solverSlugs[i])
			}
		}

		// Difficulty based (Easy Peasy, Medium Well, Hard Core)
		for _, dist := range stats.SolvedDistribution {
			var slugs []string
			var counts []int
			switch dist.Difficulty {
			case "Easy":
				slugs = []string{"easy-peasy-i", "easy-peasy-ii", "easy-peasy-iii", "easy-peasy-iv", "easy-peasy-v"}
				counts = []int{10, 50, 100, 200, 500}
			case "Medium":
				slugs = []string{"medium-well-i", "medium-well-ii", "medium-well-iii", "medium-well-iv", "medium-well-v"}
				counts = []int{10, 50, 100, 200, 500}
			case "Hard":
				slugs = []string{"hard-core-i", "hard-core-ii", "hard-core-iii", "hard-core-iv", "hard-core-v"}
				counts = []int{10, 50, 100, 200, 500}
			}

			for i, count := range counts {
				if dist.Count >= count {
					_ = u.CheckAndUnlock(userID, slugs[i])
				}
			}
		}

		// One Shot (Solved on first attempt)
		attempts, err := u.submissionRepo.CountByUserProblem(userID, submission.ProblemID)
		if err == nil && attempts == 1 {
			_ = u.CheckAndUnlock(userID, "one-shot")
		}

		// Persistence (Solved after 10+ attempts)
		// We subtract 1 because the current successful submission is included
		if err == nil && (attempts-1) >= 10 {
			_ = u.CheckAndUnlock(userID, "persistence")
		}

	} else if submission.Status == domain.SubmissionStatusWrongAnswer {
		_ = u.CheckAndUnlock(userID, "bug-hunter")
	} else if submission.Status == domain.SubmissionStatusTimeLimitExceeded {
		_ = u.CheckAndUnlock(userID, "speed-demon")
	} else if submission.Status == domain.SubmissionStatusMemoryLimitExceeded {
		_ = u.CheckAndUnlock(userID, "memory-leak")
	}

	// 3. Streak based
	streakSlugs := []string{"getting-serious", "weekly-warrior", "fortnight-fighter", "monthly-master", "century-club", "year-of-code"}
	streakValues := []int{3, 7, 14, 30, 100, 365}
	for i, val := range streakValues {
		if stats.Streak >= val {
			_ = u.CheckAndUnlock(userID, streakSlugs[i])
		}
	}

	return nil
}

func (u *AchievementUsecase) ListAll() ([]domain.Achievement, error) {
	return u.achievementRepo.GetAll()
}

func (u *AchievementUsecase) GetUserProgress(userID int) ([]domain.UserAchievement, error) {
	return u.achievementRepo.GetUnlockedByUser(userID)
}
