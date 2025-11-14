package validator

import (
	"regexp"
	"strings"
)

var (
	EmailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	UsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
)

// IsValidEmail checks if email format is valid
func IsValidEmail(email string) bool {
	return EmailRegex.MatchString(email)
}

// IsValidUsername checks if username format is valid
func IsValidUsername(username string) bool {
	return UsernameRegex.MatchString(username)
}

// HasUpperCase checks if string contains uppercase letter
func HasUpperCase(s string) bool {
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return true
		}
	}
	return false
}

// HasLowerCase checks if string contains lowercase letter
func HasLowerCase(s string) bool {
	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			return true
		}
	}
	return false
}

// HasDigit checks if string contains digit
func HasDigit(s string) bool {
	for _, r := range s {
		if r >= '0' && r <= '9' {
			return true
		}
	}
	return false
}

// NormalizeEmail converts email to lowercase and trims spaces
func NormalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

// NormalizeUsername trims spaces from username
func NormalizeUsername(username string) string {
	return strings.TrimSpace(username)
}
