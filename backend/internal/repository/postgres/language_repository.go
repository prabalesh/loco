package postgres

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/pkg/database"
)

type languageRepository struct {
	db *database.Database
}

func NewLanguageRepository(db *database.Database) domain.LanguageRepository {
	return &languageRepository{db: db}
}

func (r *languageRepository) Create(lang *domain.Language) error {
	result := r.db.DB.Create(lang)
	if result.Error != nil {
		if isUniqueViolation(result.Error) {
			if containsField(result.Error, "slug") || containsField(result.Error, "language_id") {
				return fmt.Errorf("language Slug already exists")
			}
		}
		return fmt.Errorf("failed to create language: %w", result.Error)
	}
	return nil
}

func (r *languageRepository) GetByID(id int) (*domain.Language, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	lang := &domain.Language{}
	err := r.db.DB.WithContext(ctx).First(lang, id).Error

	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("language not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get language: %w", err)
	}

	return lang, nil
}

func (r *languageRepository) GetBySlug(slug string) (*domain.Language, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	lang := &domain.Language{}
	err := r.db.DB.WithContext(ctx).Where("slug = ?", slug).First(lang).Error

	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("language not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get language by slug: %w", err)
	}

	return lang, nil
}

func (r *languageRepository) Update(lang *domain.Language) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Model(&domain.Language{}).Where("id = ?", lang.ID).Updates(lang)

	if result.Error != nil {
		return fmt.Errorf("failed to update language: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("language not found")
	}

	return nil
}

func (r *languageRepository) Delete(id int) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Delete(&domain.Language{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete language: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("language not found")
	}

	return nil
}

func (r *languageRepository) ListActive() ([]domain.Language, error) {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	var languages []domain.Language
	err := r.db.DB.WithContext(ctx).Where("is_active = ?", true).Order("name ASC").Find(&languages).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list active languages: %w", err)
	}
	return languages, nil
}

func (r *languageRepository) GetAll() ([]domain.Language, error) {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	var languages []domain.Language
	err := r.db.DB.WithContext(ctx).Find(&languages).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list languages: %w", err)
	}
	return languages, nil
}
