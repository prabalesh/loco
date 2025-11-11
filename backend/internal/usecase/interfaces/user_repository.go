package interfaces

import "github.com/prabalesh/loco/backend/internal/domain"

type UserRepository interface {
	Create(user *domain.User) error
	GetByEmail(email string) (*domain.User, error)
	GetByID(id int) (*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
}
