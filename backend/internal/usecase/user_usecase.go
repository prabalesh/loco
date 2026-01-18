package usecase

import (
	"errors"

	"github.com/prabalesh/loco/backend/internal/domain"
	"go.uber.org/zap"
)

type UserUsecase struct {
	userRepo       domain.UserRepository
	submissionRepo domain.SubmissionRepository
	logger         *zap.Logger
}

func NewUserUsecase(userRepo domain.UserRepository, submissionRepo domain.SubmissionRepository, logger *zap.Logger) *UserUsecase {
	return &UserUsecase{
		userRepo:       userRepo,
		submissionRepo: submissionRepo,
		logger:         logger,
	}
}

// GetUserProfile returns full user profile response for the given user ID
func (u *UserUsecase) GetUserProfile(userID int) (*domain.UserProfileResponse, error) {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		u.logger.Error("Failed to get user by ID",
			zap.Error(err),
			zap.Int("user_id", userID),
		)
		return nil, errors.New("user not found")
	}

	// Fetch rank
	rank, err := u.userRepo.GetUserRank(user.ID)
	if err != nil {
		u.logger.Error("Failed to get user rank", zap.Error(err))
	}

	// Fetch stats
	stats, err := u.getUserStats(user.ID, rank)
	if err != nil {
		u.logger.Error("Failed to get user stats", zap.Error(err))
	}

	// Fetch recent problems
	recentProblems, err := u.submissionRepo.FindSolvedProblemsByUser(user.ID, 5) // Limit to 5
	if err != nil {
		u.logger.Error("Failed to get recent problems", zap.Error(err))
		recentProblems = []domain.Problem{}
	}

	// Fetch distribution
	distribution, err := u.submissionRepo.GetSolvedDistribution(user.ID)
	if err != nil {
		u.logger.Error("Failed to get solved distribution", zap.Error(err))
		distribution = []domain.DifficultyStat{}
	}

	// Fetch heatmap
	heatmap, err := u.submissionRepo.GetSubmissionHeatmap(user.ID)
	if err != nil {
		u.logger.Error("Failed to get submission heatmap", zap.Error(err))
		heatmap = []domain.HeatmapEntry{}
	}

	resp := user.ToUserProfileResponse(stats, recentProblems, heatmap, distribution)
	return &resp, nil
}

func (u *UserUsecase) GetByUsername(username string) (*domain.User, error) {
	user, err := u.userRepo.GetByUsername(username)
	if err != nil {
		u.logger.Error("Failed to get user by username",
			zap.Error(err),
			zap.String("username", username),
		)
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (u *UserUsecase) GetUserProfileByUsername(username string) (*domain.UserProfileResponse, error) {
	user, err := u.userRepo.GetByUsername(username)
	if err != nil {
		u.logger.Error("Failed to get user by username",
			zap.Error(err),
			zap.String("username", username),
		)
		return nil, errors.New("user not found")
	}

	// Fetch rank
	rank, err := u.userRepo.GetUserRank(user.ID)
	if err != nil {
		u.logger.Error("Failed to get user rank", zap.Error(err))
	}

	// Fetch stats
	stats, err := u.getUserStats(user.ID, rank)
	if err != nil {
		u.logger.Error("Failed to get user stats", zap.Error(err))
	}

	// Fetch recent problems
	recentProblems, err := u.submissionRepo.FindSolvedProblemsByUser(user.ID, 5) // Limit to 5
	if err != nil {
		u.logger.Error("Failed to get recent problems", zap.Error(err))
		recentProblems = []domain.Problem{}
	}

	// Fetch distribution
	distribution, err := u.submissionRepo.GetSolvedDistribution(user.ID)
	if err != nil {
		u.logger.Error("Failed to get solved distribution", zap.Error(err))
		distribution = []domain.DifficultyStat{}
	}

	// Fetch heatmap
	heatmap, err := u.submissionRepo.GetSubmissionHeatmap(user.ID)
	if err != nil {
		u.logger.Error("Failed to get submission heatmap", zap.Error(err))
		heatmap = []domain.HeatmapEntry{}
	}

	resp := user.ToUserProfileResponse(stats, recentProblems, heatmap, distribution)
	return &resp, nil
}

func (u *UserUsecase) getUserStats(userID int, rank int) (domain.UserStats, error) {
	totalSubmissions, err := u.submissionRepo.CountByUser(userID)
	if err != nil {
		return domain.UserStats{}, err
	}

	acceptedSubmissions, err := u.submissionRepo.CountAcceptedByUser(userID)
	if err != nil {
		return domain.UserStats{}, err
	}

	problemsSolved, err := u.submissionRepo.CountProblemsSolvedByUser(userID)
	if err != nil {
		return domain.UserStats{}, err
	}

	streak, err := u.submissionRepo.GetCurrentStreak(userID)
	if err != nil {
		u.logger.Error("Failed to get user streak", zap.Error(err))
	}

	distribution, err := u.submissionRepo.GetSolvedDistribution(userID)
	if err != nil {
		u.logger.Error("Failed to get solved distribution", zap.Error(err))
	}

	acceptanceRate := 0.0
	if totalSubmissions > 0 {
		acceptanceRate = float64(acceptedSubmissions) / float64(totalSubmissions) * 100
	}

	return domain.UserStats{
		TotalSubmissions:    int(totalSubmissions),
		AcceptedSubmissions: int(acceptedSubmissions),
		ProblemsSolved:      int(problemsSolved),
		AcceptanceRate:      acceptanceRate,
		Rank:                rank,
		Streak:              streak,
		SolvedDistribution:  distribution,
	}, nil
}
