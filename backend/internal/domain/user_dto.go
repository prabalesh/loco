package domain

import "time"

// ==================== REQUEST DTOs ====================

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
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

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int          `json:"expires_in"`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	ID            int       `json:"id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	Role          string    `json:"role"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
}
