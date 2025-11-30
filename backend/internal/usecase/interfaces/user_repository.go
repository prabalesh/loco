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
	SetPasswordResetToken(userID int, token string, expiresAt time.Time) error
	ClearPasswordResetToken(userID int) error
	GetUserByResetToken(token string) (*domain.User, error)
	UpdatePassword(userID int, newPasswordHash string) error
	GetByPasswordResetToken(token string) (*domain.User, error)
	GetByVerificationToken(token string) (*domain.User, error)
	UpdatePasswordResetToken(userID int, token string, expiresAt time.Time, sentAt time.Time) error

	GetAll() ([]*domain.User, error)
	Delete(id int) error
	UpdateRole(id int, role string) error
	UpdateActiveStatus(id int, isActive bool) error
	CountUsers() (int, error)
	CountActiveUsers() (int, error)
	CountVerifiedUsers() (int, error)
}
