package domain

import "time"

// sturcts
type TokenPair struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  time.Duration
	RefreshExpiresAt time.Duration
}

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

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
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
	Message string       `json:"message"`
	User    UserResponse `json:"user"`
}
