package domain

import "time"

// ==================== REQUEST DTOs ====================

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type VerifyEmailRequest struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateProfileRequest struct {
	Username string `json:"username,omitempty"`
}

// ==================== RESPONSE DTOs ====================

type RegisterResponse struct {
	Message string       `json:"message"`
	User    UserResponse `json:"user"`
}

type ResendVerificationRequest struct {
	Email string `json:"email"`
}

type LoginResponse struct {
	Message string       `json:"message`
	User    UserResponse `json:"user"`
}

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
}

const (
	AccessTokenCookie  = "access_token"
	RefreshTokenCookie = "refresh_token"
	AccessTokenMaxAge  = 15 * 60          // 15 minutes
	RefreshTokenMaxAge = 7 * 24 * 60 * 60 // 7 days
)
