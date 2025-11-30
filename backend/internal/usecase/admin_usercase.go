// internal/usecase/admin_usecase.go
package usecase

import (
	"errors"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/usecase/interfaces"
	"go.uber.org/zap"
)

type AdminUsecase struct {
	userRepo interfaces.UserRepository
	logger   *zap.Logger
}

func NewAdminUsecase(userRepo interfaces.UserRepository, logger *zap.Logger) *AdminUsecase {
	return &AdminUsecase{
		userRepo: userRepo,
		logger:   logger,
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

	analytics := &domain.AdminAnalytics{
		TotalUsers:    totalUsers,
		ActiveUsers:   activeUsers,
		InactiveUsers: totalUsers - activeUsers,
		VerifiedUsers: verifiedUsers,
	}

	return analytics, nil
}
