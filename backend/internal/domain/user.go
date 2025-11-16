package domain

import "time"

// User entity - just the database model
type User struct {
	ID                              int        `json:"id" db:"id"`
	Email                           string     `json:"email" db:"email"`
	Username                        string     `json:"username" db:"username"`
	PasswordHash                    string     `json:"-" db:"password_hash"`
	Role                            string     `json:"role" db:"role"`
	IsActive                        bool       `json:"is_active" db:"is_active"`
	EmailVerified                   bool       `json:"email_verified" db:"email_verified"`
	EmailVerificationToken          *string    `json:"-" db:"email_verification_token"`
	EmailVerificationTokenExpiresAt *time.Time `json:"-" db:"email_verification_token_expires_at"`
	EmailVerificationAttempts       int        `json:"-" db:"email_verification_attempts"`
	EmailVerificationLastSentAt     *time.Time `json:"-" db:"email_verification_last_sent_at"`
	CreatedAt                       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                       time.Time  `json:"updated_at" db:"updated_at"`
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

func (u *User) ToUserProfileResponse() UserProfileResponse {
	return UserProfileResponse{
		ID:         u.ID,
		Username:   u.Username,
		IsVerified: u.EmailVerified,
		CreatedAt:  u.CreatedAt,
	}
}
