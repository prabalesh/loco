package postgres

import (
	"github.com/prabalesh/loco/backend/internal/domain"

	"gorm.io/gorm"
)

type typeImplementationRepository struct {
	db *gorm.DB
}

func NewTypeImplementationRepository(db *gorm.DB) domain.TypeImplementationRepository {
	return &typeImplementationRepository{db: db}
}

func (r *typeImplementationRepository) Create(impl *domain.TypeImplementation) error {
	return r.db.Create(impl).Error
}

func (r *typeImplementationRepository) GetByTypeAndLanguage(customTypeID, languageID int) (*domain.TypeImplementation, error) {
	var impl domain.TypeImplementation
	if err := r.db.Where("custom_type_id = ? AND language_id = ?", customTypeID, languageID).First(&impl).Error; err != nil {
		return nil, err
	}
	return &impl, nil
}

func (r *typeImplementationRepository) GetByTypeAndLanguageSlug(typeName, languageSlug string) (*domain.TypeImplementation, error) {
	var impl domain.TypeImplementation
	err := r.db.
		Joins("JOIN custom_types ON custom_types.id = type_implementations.custom_type_id").
		Joins("JOIN languages ON languages.id = type_implementations.language_id").
		Where("custom_types.name = ? AND languages.slug = ?", typeName, languageSlug).
		First(&impl).Error
	if err != nil {
		return nil, err
	}
	return &impl, nil
}
