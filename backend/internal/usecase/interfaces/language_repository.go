package interfaces

import (
	"context"

	"github.com/prabalesh/loco/backend/internal/domain"
)

// LanguageRepository defines methods for managing Language entities
type LanguageRepository interface {
	Create(ctx context.Context, lang *domain.Language) error
	GetByID(ctx context.Context, id int) (*domain.Language, error)
	GetByLanguageID(ctx context.Context, languageID string) (*domain.Language, error)
	Update(ctx context.Context, lang *domain.Language) error
	Delete(ctx context.Context, id int) error
	ListActive() ([]*domain.Language, error)
	List() ([]*domain.Language, error)
}
