// internal/usecase/admin_usecase.go
package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type AdminUsecase struct {
	userRepo       domain.UserRepository
	submissionRepo domain.SubmissionRepository
	redis          *redis.Client
	logger         *zap.Logger
}

func NewAdminUsecase(userRepo domain.UserRepository, submissionRepo domain.SubmissionRepository, redis *redis.Client, logger *zap.Logger) *AdminUsecase {
	return &AdminUsecase{
		userRepo:       userRepo,
		submissionRepo: submissionRepo,
		redis:          redis,
		logger:         logger,
	}
}

// GetAllUsers - Fetch all users
func (u *AdminUsecase) GetAllUsers() ([]*domain.User, error) {
	users, err := u.userRepo.GetAll()
	if err != nil {
		u.logger.Error("Failed to fetch all users", zap.Error(err))
		return nil, errors.New("failed to fetch users")
	}
	return users, nil
}

// GetUserByID - Get user by ID
func (u *AdminUsecase) GetUserByID(userID int) (*domain.User, error) {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		u.logger.Error("Failed to fetch user", zap.Error(err), zap.Int("user_id", userID))
		return nil, errors.New("user not found")
	}
	return user, nil
}

// DeleteUser - Delete user (with admin audit logging)
func (u *AdminUsecase) DeleteUser(adminID, userID int) error {
	// Prevent admin from deleting themselves
	if adminID == userID {
		return errors.New("cannot delete your own account")
	}

	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Prevent deleting other admins (optional business rule)
	if user.Role == "admin" {
		u.logger.Warn("Attempt to delete admin user",
			zap.Int("admin_id", adminID),
			zap.Int("target_user_id", userID),
		)
		return errors.New("cannot delete admin users")
	}

	if err := u.userRepo.Delete(userID); err != nil {
		u.logger.Error("Failed to delete user", zap.Error(err), zap.Int("user_id", userID))
		return errors.New("failed to delete user")
	}

	u.logger.Info("User deleted by admin",
		zap.Int("admin_id", adminID),
		zap.Int("deleted_user_id", userID),
	)

	return nil
}

// UpdateUserRole - Change user role
func (u *AdminUsecase) UpdateUserRole(adminID, userID int, newRole string) error {
	// Validate role
	validRoles := map[string]bool{
		"user":      true,
		"admin":     true,
		"moderator": true,
	}

	if !validRoles[newRole] {
		return errors.New("invalid role")
	}

	// Prevent admin from changing their own role
	if adminID == userID {
		return errors.New("cannot change your own role")
	}

	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// No-op if role is same
	if user.Role == newRole {
		return nil
	}

	if err := u.userRepo.UpdateRole(userID, newRole); err != nil {
		u.logger.Error("Failed to update user role", zap.Error(err), zap.Int("user_id", userID))
		return errors.New("failed to update role")
	}

	u.logger.Info("User role updated",
		zap.Int("admin_id", adminID),
		zap.Int("user_id", userID),
		zap.String("old_role", user.Role),
		zap.String("new_role", newRole),
	)

	return nil
}

// UpdateUserStatus - Activate/Deactivate user
func (u *AdminUsecase) UpdateUserStatus(adminID, userID int, isActive bool) error {
	// Prevent admin from deactivating themselves
	if adminID == userID {
		return errors.New("cannot change your own status")
	}

	if err := u.userRepo.UpdateActiveStatus(userID, isActive); err != nil {
		u.logger.Error("Failed to update user status", zap.Error(err), zap.Int("user_id", userID))
		return errors.New("failed to update status")
	}

	action := "deactivated"
	if isActive {
		action = "activated"
	}

	u.logger.Info("User status updated",
		zap.Int("admin_id", adminID),
		zap.Int("user_id", userID),
		zap.String("action", action),
	)

	return nil
}

// GetAnalytics - Dashboard statistics
func (u *AdminUsecase) GetAnalytics() (*domain.AdminAnalytics, error) {
	totalUsers, err := u.userRepo.CountUsers()
	if err != nil {
		return nil, errors.New("failed to get analytics")
	}

	activeUsers, err := u.userRepo.CountActiveUsers()
	if err != nil {
		return nil, errors.New("failed to get analytics")
	}

	verifiedUsers, err := u.userRepo.CountVerifiedUsers()
	if err != nil {
		return nil, errors.New("failed to get analytics")
	}

	totalSubmissions, err := u.submissionRepo.CountTotal()
	if err != nil {
		u.logger.Error("Failed to count submissions", zap.Error(err))
	}

	pendingSubmissions, _ := u.submissionRepo.CountPending()

	var queueSize int64
	if u.redis != nil {
		queueSize, _ = u.redis.LLen(context.Background(), "submission:queue").Result()
	}

	// Calculate oldest pending submission age
	oldestAge := int64(0)
	if pendingSubmissions > 0 {
		oldestSubmissions, err := u.submissionRepo.GetOldestPending(1)
		if err == nil && len(oldestSubmissions) > 0 {
			oldestAge = int64(time.Since(oldestSubmissions[0].CreatedAt).Seconds())
		}
	}

	// Determine queue health
	queueHealthStatus := "healthy"
	if pendingSubmissions > 0 && queueSize >= pendingSubmissions {
		// No workers processing
		queueHealthStatus = "critical"
	} else if oldestAge > 300 {
		// Submissions waiting > 5 minutes
		queueHealthStatus = "critical"
	} else if oldestAge > 120 || queueSize > 10 {
		// Submissions waiting > 2 minutes or queue backed up
		queueHealthStatus = "warning"
	}

	// Fetch daily stats (last 7 days)
	dailyStats, err := u.submissionRepo.GetDailyStats(7)
	if err != nil {
		u.logger.Error("Failed to get daily stats", zap.Error(err))
		dailyStats = []domain.DailySubmissionStat{}
	}

	activeWorkers, err := u.countActiveWorkers(context.Background())
	if err != nil {
		u.logger.Error("Failed to count active workers", zap.Error(err))
		activeWorkers = 0
	}

	// Fetch trending problems (last 7 days, top 5)
	trendingProblems, err := u.submissionRepo.GetTrendingProblems(5, 7)
	if err != nil {
		u.logger.Error("Failed to get trending problems", zap.Error(err))
		trendingProblems = []domain.TrendingProblem{}
	}

	// Fetch language stats
	languageStats, err := u.submissionRepo.GetLanguageStats()
	if err != nil {
		u.logger.Error("Failed to get language stats", zap.Error(err))
		languageStats = []domain.LanguageStat{}
	}

	analytics := &domain.AdminAnalytics{
		TotalUsers:         totalUsers,
		ActiveUsers:        activeUsers,
		InactiveUsers:      totalUsers - activeUsers,
		VerifiedUsers:      verifiedUsers,
		TotalSubmissions:   int(totalSubmissions),
		PendingSubmissions: int(pendingSubmissions),
		ActiveWorkers:      activeWorkers,
		QueueSize:          queueSize,
		OldestPendingAge:   oldestAge,
		QueueHealthStatus:  queueHealthStatus,
		SubmissionHistory:  dailyStats,
		TrendingProblems:   trendingProblems,
		LanguageStats:      languageStats,
	}

	return analytics, nil
}

func (u *AdminUsecase) countActiveWorkers(ctx context.Context) (int, error) {
	keys, err := u.redis.Keys(ctx, "worker:*:heartbeat").Result()
	if err != nil {
		return 0, err
	}
	return len(keys), nil
}
