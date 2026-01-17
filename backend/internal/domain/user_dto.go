package domain

import "time"

// ==================== REQUEST DTOs ====================
// ==================== RESPONSE DTOs ====================

type UserResponse struct {
	ID            int       `json:"id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	Role          string    `json:"role"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
}

type UserProfileResponse struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	Stats      UserStats `json:"stats"`
}
