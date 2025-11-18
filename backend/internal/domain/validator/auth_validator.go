package validator

import (
	"github.com/prabalesh/loco/backend/internal/domain"
)

func ValidateRegisterRequest(req *domain.RegisterRequest) map[string]string {
	errors := make(map[string]string)

	// Normalize inputs
	req.Email = NormalizeEmail(req.Email)
	req.Username = NormalizeUsername(req.Username)

	// Validate email
	if req.Email == "" {
		errors["email"] = "email is required"
	} else if !IsValidEmail(req.Email) {
		errors["email"] = "invalid email format"
	} else if len(req.Email) > 255 {
		errors["email"] = "email too long (max 255 characters)"
	}

	// Validate username
	if req.Username == "" {
		errors["username"] = "username is required"
	} else if len(req.Username) < 3 {
		errors["username"] = "username must be at least 3 characters"
	} else if len(req.Username) > 50 {
		errors["username"] = "username too long (max 50 characters)"
	} else if !IsValidUsername(req.Username) {
		errors["username"] = "username can only contain letters, numbers, and underscores"
	}

	// Validate password
	if req.Password == "" {
		errors["password"] = "password is required"
	} else if len(req.Password) < 8 {
		errors["password"] = "password must be at least 8 characters"
	} else if len(req.Password) > 72 {
		errors["password"] = "password too long (max 72 characters)"
	} else if !HasUpperCase(req.Password) {
		errors["password"] = "password must contain at least one uppercase letter"
	} else if !HasLowerCase(req.Password) {
		errors["password"] = "password must contain at least one lowercase letter"
	} else if !HasDigit(req.Password) {
		errors["password"] = "password must contain at least one number"
	}

	return errors
}

// ValidateLoginRequest validates login request
func ValidateLoginRequest(req *domain.LoginRequest) map[string]string {
	errors := make(map[string]string)

	req.Email = NormalizeEmail(req.Email)

	if req.Email == "" {
		errors["email"] = "email is required"
	} else if !IsValidEmail(req.Email) {
		errors["email"] = "invalid email format"
	}

	if req.Password == "" {
		errors["password"] = "password is required"
	}

	return errors
}

func ValidateResetPasswordRequest(password string) map[string]string {
	errors := make(map[string]string)

	if password == "" {
		errors["new_password"] = "password is required"
	} else if len(password) < 8 {
		errors["new_password"] = "password must be at least 8 characters"
	} else if len(password) > 72 {
		errors["new_password"] = "password too long (max 72 characters)"
	} else if !HasUpperCase(password) {
		errors["new_password"] = "password must contain at least one uppercase letter"
	} else if !HasLowerCase(password) {
		errors["new_password"] = "password must contain at least one lowercase letter"
	} else if !HasDigit(password) {
		errors["new_password"] = "password must contain at least one number"
	}

	return errors
}
