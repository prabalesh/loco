package postgres

import (
	"database/sql"
	"fmt"
	"time"

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

	query := `
		INSERT INTO users (email, username, password_hash, role, is_active, email_verified)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.Role,
		user.IsActive,
		user.EmailVerified,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		// Check for unique constraint violation
		if isUniqueViolation(err) {
			if containsField(err, "email") {
				return fmt.Errorf("email already exists")
			}
			if containsField(err, "username") {
				return fmt.Errorf("username already taken")
			}
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByEmail retrieves user by email with all verification fields
func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	user := &domain.User{}
	query := `
		SELECT id, email, username, password_hash, role, is_active, 
		       email_verified, email_verification_token, 
		       email_verification_token_expires_at, email_verification_attempts,
		       email_verification_last_sent_at, created_at, updated_at
		FROM users WHERE email = $1
	`

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.EmailVerified,
		&user.EmailVerificationToken,
		&user.EmailVerificationTokenExpiresAt,
		&user.EmailVerificationAttempts,
		&user.EmailVerificationLastSentAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
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
	query := `
		SELECT id, email, username, password_hash, role, is_active,
		       email_verified, created_at, updated_at
		FROM users WHERE username = $1
	`

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
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
	query := `
		SELECT id, email, username, password_hash, role, is_active, 
		       email_verified, created_at, updated_at
		FROM users WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
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

	query := `
		UPDATE users 
		SET email_verification_token = $1,
		    email_verification_token_expires_at = $2,
		    email_verification_attempts = 0,
		    updated_at = NOW()
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, token, expiresAt, userID)
	if err != nil {
		return fmt.Errorf("failed to update verification token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// UpdateVerificationAttempts increments failed verification attempts
func (r *userRepository) UpdateVerificationAttempts(userID int, attempts int) error {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	query := `
		UPDATE users 
		SET email_verification_attempts = $1,
		    updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, attempts, userID)
	if err != nil {
		return fmt.Errorf("failed to update verification attempts: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// UpdateLastSentAt records when verification email was last sent
func (r *userRepository) UpdateLastSentAt(userID int, sentAt time.Time) error {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	query := `
		UPDATE users 
		SET email_verification_last_sent_at = $1,
		    updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, sentAt, userID)
	if err != nil {
		return fmt.Errorf("failed to update last sent time: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// VerifyEmail marks email as verified and clears verification data
func (r *userRepository) VerifyEmail(userID int) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	query := `
		UPDATE users 
		SET email_verified = true,
		    email_verification_token = NULL,
		    email_verification_token_expires_at = NULL,
		    email_verification_attempts = 0,
		    updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) SetPasswordResetToken(userID int, token string, expiresAt time.Time) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	query := `
        UPDATE users
        SET password_reset_token = $1,
            password_reset_token_expires_at = $2,
            password_reset_attempts = 0,
            updated_at = NOW()
        WHERE id = $3
    `

	result, err := r.db.ExecContext(ctx, query, token, expiresAt, userID)
	if err != nil {
		return fmt.Errorf("failed to set reset token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// ClearPasswordResetToken clears reset token fields after successful reset
func (r *userRepository) ClearPasswordResetToken(userID int) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	query := `
        UPDATE users
        SET password_reset_token = NULL,
            password_reset_token_expires_at = NULL,
            password_reset_attempts = 0,
            updated_at = NOW()
        WHERE id = $1
    `

	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to clear reset token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// GetUserByResetToken retrieves user for given reset token
func (r *userRepository) GetUserByResetToken(token string) (*domain.User, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	user := &domain.User{}
	query := `
        SELECT id, email, username, password_hash, role, is_active,
               email_verified, created_at, updated_at,
               password_reset_token, password_reset_token_expires_at, password_reset_attempts
        FROM users WHERE password_reset_token = $1
    `

	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.PasswordResetToken,
		&user.PasswordResetTokenExpiresAt,
	)

	if err == sql.ErrNoRows {
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

	query := `UPDATE users SET password_hash = $1, password_reset_token = NULL, password_reset_token_expires_at = NULL WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, hashedPassword, userID)
	return err
}

func (r *userRepository) GetByVerificationToken(token string) (*domain.User, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	user := &domain.User{}
	query := `
		SELECT id, email, username, password_hash, role, is_active, 
		       email_verified, email_verification_token, 
		       email_verification_token_expires_at, email_verification_attempts,
		       email_verification_last_sent_at, created_at, updated_at
		FROM users WHERE email_verification_token = $1
	`

	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.EmailVerified,
		&user.EmailVerificationToken,
		&user.EmailVerificationTokenExpiresAt,
		&user.EmailVerificationAttempts,
		&user.EmailVerificationLastSentAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
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
	query := `
		SELECT id, email, username, password_hash, role, is_active,
		       password_reset_token, password_reset_token_expires_at,
		       password_reset_sent_at, created_at, updated_at
		FROM users
		WHERE password_reset_token = $1
	`
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.Role, &user.IsActive,
		&user.PasswordResetToken, &user.PasswordResetTokenExpiresAt,
		&user.PasswordResetSentAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("user not found for reset token")
	}
	return user, nil
}

// Update password reset token and expiry, sent time (for initiating reset)
func (r *userRepository) UpdatePasswordResetToken(userID int, token string, expiresAt time.Time, sentAt time.Time) error {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	query := `UPDATE users SET password_reset_token = $1, password_reset_token_expires_at = $2, password_reset_sent_at = $3 WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query, token, expiresAt, sentAt, userID)
	return err
}
