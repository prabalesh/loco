package postgres

import (
	"fmt"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/pkg/database"
)

type problemLanguageRepository struct {
	db *database.Database
}

func NewProblemLanguageRepository(db *database.Database) *problemLanguageRepository {
	return &problemLanguageRepository{db: db}
}

func (r *problemLanguageRepository) Create(problemLanguage *domain.ProblemLanguage) error {
	ctx, cancel := database.WithLongTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Create(problemLanguage)
	if result.Error != nil {
		return fmt.Errorf("failed to create problem language: %w", result.Error)
	}

	return nil
}

func (r *problemLanguageRepository) GetByProblemID(problemID int) ([]domain.ProblemLanguage, error) {
	ctx, cancel := database.WithLongTimeout()
	defer cancel()

	var problemLanguages []domain.ProblemLanguage
	result := r.db.DB.WithContext(ctx).
		Table("problem_languages").
		Select("problem_languages.*, languages.name as language_name, languages.version as language_version").
		Joins("join languages on languages.id = problem_languages.language_id").
		Where("problem_languages.problem_id = ?", problemID).
		Preload("Language").
		Find(&problemLanguages)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get problem languages: %w", result.Error)
	}

	return problemLanguages, nil
}

func (r *problemLanguageRepository) GetByProblemAndLanguage(problemID int, languageID int) (*domain.ProblemLanguage, error) {
	ctx, cancel := database.WithLongTimeout()
	defer cancel()

	var problemLanguage domain.ProblemLanguage
	result := r.db.DB.WithContext(ctx).
		Table("problem_languages").
		Select("problem_languages.*, languages.name as language_name, languages.version as language_version").
		Joins("join languages on languages.id = problem_languages.language_id").
		Where("problem_languages.problem_id = ? AND problem_languages.language_id = ?", problemID, languageID).
		Preload("Language").
		First(&problemLanguage)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get problem language: %w", result.Error)
	}

	return &problemLanguage, nil
}

func (r *problemLanguageRepository) Update(problemLanguage *domain.ProblemLanguage) error {
	ctx, cancel := database.WithLongTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Save(problemLanguage)
	if result.Error != nil {
		return fmt.Errorf("failed to update problem language: %w", result.Error)
	}

	return nil
}

func (r *problemLanguageRepository) Delete(problemID int, languageID int) error {
	ctx, cancel := database.WithLongTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Where("problem_id = ? AND language_id = ?", problemID, languageID).Delete(&domain.ProblemLanguage{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete problem language: %w", result.Error)
	}

	return nil
}
