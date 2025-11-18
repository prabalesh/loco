package uerror

import (
	"errors"
	"fmt"
	"strings"
)

// variables
var (
	ErrEmailNotVerified         = errors.New("email not verified")
	ErrInvalidToken             = errors.New("invalid or expired token")
	ErrMaxTokenAttemptsExceeded = errors.New("maximum token attempts exceeded")
	ErrResendCooldown           = errors.New("please wait before requesting a new token")
)

// struct
type ValidationError struct {
	Errors map[string]string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %v", e.Errors)
}

// functions
func IsNotFoundError(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "not found")
}
