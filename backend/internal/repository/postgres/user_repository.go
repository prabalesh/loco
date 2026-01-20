package postgres

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/pkg/database"
)

type userRepository struct {
	db *database.Database
}

func NewUserRepository(db *database.Database) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Create(user)
	if result.Error != nil {
		if isUniqueViolation(result.Error) {
			if containsField(result.Error, "email") {
				return fmt.Errorf("email already exists")
			}
			if containsField(result.Error, "username") {
				return fmt.Errorf("username already taken")
			}
		}
		return fmt.Errorf("failed to create user: %w", result.Error)
	}

	return nil
}

func (r *userRepository) Update(user *domain.User) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	if err := r.db.DB.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// GetByEmail retrieves user by email with all verification fields
func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	user := &domain.User{}
	err := r.db.DB.WithContext(ctx).Where("email = ?", email).First(user).Error

	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetByUsername(username string) (*domain.User, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	user := &domain.User{}
	err := r.db.DB.WithContext(ctx).Where("username = ?", username).First(user).Error

	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}

	return user, nil
}

// GetByID retrieves user by ID
func (r *userRepository) GetByID(userID int) (*domain.User, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	user := &domain.User{}
	err := r.db.DB.WithContext(ctx).First(user, userID).Error

	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// ========== EMAIL VERIFICATION METHODS ==========

// UpdateVerificationToken sets new OTP token and resets attempts
func (r *userRepository) UpdateVerificationToken(userID int, token string, expiresAt time.Time) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Model(&domain.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"email_verification_token":            token,
		"email_verification_token_expires_at": expiresAt,
		"email_verification_attempts":         0,
	})

	if result.Error != nil {
		return fmt.Errorf("failed to update verification token: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// UpdateVerificationAttempts increments failed verification attempts
func (r *userRepository) UpdateVerificationAttempts(userID int, attempts int) error {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Model(&domain.User{}).Where("id = ?", userID).Update("email_verification_attempts", attempts)

	if result.Error != nil {
		return fmt.Errorf("failed to update verification attempts: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// UpdateLastSentAt records when verification email was last sent
func (r *userRepository) UpdateLastSentAt(userID int, sentAt time.Time) error {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Model(&domain.User{}).Where("id = ?", userID).Update("email_verification_last_sent_at", sentAt)

	if result.Error != nil {
		return fmt.Errorf("failed to update last sent time: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// VerifyEmail marks email as verified and clears verification data
func (r *userRepository) VerifyEmail(userID int) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Model(&domain.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"email_verified":                      true,
		"email_verification_token":            nil,
		"email_verification_token_expires_at": nil,
		"email_verification_attempts":         0,
	})

	if result.Error != nil {
		return fmt.Errorf("failed to verify email: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) SetPasswordResetToken(userID int, token string, expiresAt time.Time) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Model(&domain.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"password_reset_token":            token,
		"password_reset_token_expires_at": expiresAt,
		"password_reset_attempts":         0,
	})

	if result.Error != nil {
		return fmt.Errorf("failed to set reset token: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// ClearPasswordResetToken clears reset token fields after successful reset
func (r *userRepository) ClearPasswordResetToken(userID int) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Model(&domain.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"password_reset_token":            nil,
		"password_reset_token_expires_at": nil,
		"password_reset_attempts":         0,
	})

	if result.Error != nil {
		return fmt.Errorf("failed to clear reset token: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// GetUserByResetToken retrieves user for given reset token
func (r *userRepository) GetUserByResetToken(token string) (*domain.User, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	user := &domain.User{}
	err := r.db.DB.WithContext(ctx).Where("password_reset_token = ?", token).First(user).Error

	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("invalid or expired token")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	return user, nil
}

// Update/reset user password hash
func (r *userRepository) UpdatePassword(userID int, hashedPassword string) error {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Model(&domain.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"password_hash":                   hashedPassword,
		"password_reset_token":            nil,
		"password_reset_token_expires_at": nil,
	})
	return result.Error
}

func (r *userRepository) GetByVerificationToken(token string) (*domain.User, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	user := &domain.User{}
	err := r.db.DB.WithContext(ctx).Where("email_verification_token = ?", token).First(user).Error

	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// Find user by password reset token
func (r *userRepository) GetByPasswordResetToken(token string) (*domain.User, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	user := &domain.User{}
	err := r.db.DB.WithContext(ctx).Where("password_reset_token = ?", token).First(user).Error

	if err != nil {
		return nil, fmt.Errorf("user not found for reset token")
	}
	return user, nil
}

// Update password reset token and expiry, sent time (for initiating reset)
func (r *userRepository) UpdatePasswordResetToken(userID int, token string, expiresAt time.Time, sentAt time.Time) error {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Model(&domain.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"password_reset_token":            token,
		"password_reset_token_expires_at": expiresAt,
		"password_reset_sent_at":          sentAt,
	})
	return result.Error
}

// ========== ADMIN METHODS ==========

// GetAll retrieves all users (admin only)
func (r *userRepository) GetAll() ([]*domain.User, error) {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	var users []*domain.User
	err := r.db.DB.WithContext(ctx).Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	return users, nil
}

// Delete removes a user permanently
func (r *userRepository) Delete(id int) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Delete(&domain.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// UpdateRole changes user role
func (r *userRepository) UpdateRole(id int, role string) error {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Update("role", role)

	if result.Error != nil {
		return fmt.Errorf("failed to update role: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// UpdateActiveStatus activates or deactivates user
func (r *userRepository) UpdateActiveStatus(id int, isActive bool) error {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Update("is_active", isActive)

	if result.Error != nil {
		return fmt.Errorf("failed to update status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// CountUsers returns total number of users
func (r *userRepository) CountUsers() (int, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	var count int64
	err := r.db.DB.WithContext(ctx).Model(&domain.User{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return int(count), nil
}

// CountActiveUsers returns number of active users
func (r *userRepository) CountActiveUsers() (int, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	var count int64
	err := r.db.DB.WithContext(ctx).Model(&domain.User{}).Where("is_active = ?", true).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count active users: %w", err)
	}

	return int(count), nil
}

// CountVerifiedUsers returns number of verified users
func (r *userRepository) CountVerifiedUsers() (int, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	var count int64
	err := r.db.DB.WithContext(ctx).Model(&domain.User{}).Where("email_verified = ?", true).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count verified users: %w", err)
	}

	return int(count), nil
}

func (r *userRepository) GetLeaderboard(limit int) ([]domain.LeaderboardEntry, error) {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	var entries []domain.LeaderboardEntry
	query := `
		SELECT 
			u.id as user_id, 
			u.username, 
			u.xp,
			u.level,
			COALESCE(solved_counts.count, 0) as problems_solved,
			COALESCE(sub_stats.total, 0) as total_submissions,
			COALESCE(sub_stats.rate, 0) as acceptance_rate
		FROM users u
		LEFT JOIN (
			SELECT user_id, COUNT(*) as count 
			FROM user_problem_stats 
			WHERE status = 'solved' 
			GROUP BY user_id
		) solved_counts ON u.id = solved_counts.user_id
		LEFT JOIN (
			SELECT 
				user_id, 
				COUNT(*) as total,
				(COUNT(*) FILTER (WHERE status = 'Accepted')::float / NULLIF(COUNT(*), 0)) * 100 as rate
			FROM submissions
			WHERE is_validation_submission = false
			GROUP BY user_id
		) sub_stats ON u.id = sub_stats.user_id
		ORDER BY problems_solved DESC, acceptance_rate DESC, u.username ASC
		LIMIT ?
	`

	err := r.db.DB.WithContext(ctx).Raw(query, limit).Scan(&entries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch leaderboard: %w", err)
	}

	// Assign ranks
	for i := range entries {
		entries[i].Rank = i + 1
	}

	return entries, nil
}
func (r *userRepository) GetUserRank(userID int) (int, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	// Rank is based on number of problems solved, then acceptance rate
	query := `
		WITH user_stats AS (
			SELECT 
				u.id,
				COALESCE(solved_counts.count, 0) as solved,
				COALESCE(sub_stats.rate, 0) as rate
			FROM users u
			LEFT JOIN (
				SELECT user_id, COUNT(*) as count 
				FROM user_problem_stats 
				WHERE status = 'solved' 
				GROUP BY user_id
			) solved_counts ON u.id = solved_counts.user_id
			LEFT JOIN (
				SELECT 
					user_id, 
					(COUNT(*) FILTER (WHERE status = 'Accepted')::float / NULLIF(COUNT(*), 0)) * 100 as rate
				FROM submissions
				WHERE is_validation_submission = false
				GROUP BY user_id
			) sub_stats ON u.id = sub_stats.user_id
		)
		SELECT rank FROM (
			SELECT id, RANK() OVER (ORDER BY solved DESC, rate DESC) as rank
			FROM user_stats
		) ranked_users
		WHERE id = ?
	`

	var rank int
	err := r.db.DB.WithContext(ctx).Raw(query, userID).Scan(&rank).Error
	if err != nil {
		return 0, fmt.Errorf("failed to get user rank: %w", err)
	}

	return rank, nil
}
