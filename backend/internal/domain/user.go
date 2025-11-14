package domain

import "time"

// User entity - just the database model
type User struct {
	ID            int       `json:"id" db:"id"`
	Email         string    `json:"email" db:"email"`
	Username      string    `json:"username" db:"username"`
	PasswordHash  string    `json:"-" db:"password_hash"`
	Role          string    `json:"role" db:"role"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	EmailVerified bool      `json:"email_verified" db:"email_verified"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// ToUserResponse converts User entity to UserResponse DTO
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:            u.ID,
		Email:         u.Email,
		Username:      u.Username,
		Role:          u.Role,
		EmailVerified: u.EmailVerified,
		CreatedAt:     u.CreatedAt,
	}
}
