package postgres

import (
	"database/sql"
	"fmt"

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

// GetByEmail retrieves user by email
func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	user := &domain.User{}
	query := `
        SELECT id, email, username, password_hash, role, is_active, 
               email_verified, created_at, updated_at
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
	query := `SELECT id, username FROM users WHERE username = $1`

	err := r.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username)

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
