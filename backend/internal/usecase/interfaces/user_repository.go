package interfaces

import "github.com/prabalesh/loco/backend/internal/domain"

type UserRepository interface {
	Create(user *domain.User) error
	GetByEmail(email string) (*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
	GetByID(userID int) (*domain.User, error)
}
