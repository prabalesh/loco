package interfaces

import (
	"time"

	"github.com/prabalesh/loco/backend/internal/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetByEmail(email string) (*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
	GetByID(userID int) (*domain.User, error)
	UpdateVerificationToken(userID int, token string, expiresAt time.Time) error
	UpdateVerificationAttempts(userID int, attempts int) error
	UpdateLastSentAt(userID int, sentAt time.Time) error
	VerifyEmail(userID int) error
}
