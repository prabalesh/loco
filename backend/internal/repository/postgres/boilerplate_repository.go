package postgres

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/pkg/database"
)

type boilerplateRepository struct {
	db *database.Database
}

func NewBoilerplateRepository(db *database.Database) domain.BoilerplateRepository {
	return &boilerplateRepository{db: db}
}

func (r *boilerplateRepository) Create(boilerplate *domain.ProblemBoilerplate) error {
	result := r.db.DB.Create(boilerplate)
	if result.Error != nil {
		if isUniqueViolation(result.Error) {
			return fmt.Errorf("boilerplate already exists for this problem and language")
		}
		return fmt.Errorf("failed to create boilerplate: %w", result.Error)
	}
	return nil
}

func (r *boilerplateRepository) GetByProblemAndLanguage(problemID, languageID int) (*domain.ProblemBoilerplate, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	bp := &domain.ProblemBoilerplate{}
	err := r.db.DB.WithContext(ctx).
		Where("problem_id = ? AND language_id = ?", problemID, languageID).
		First(bp).Error

	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("boilerplate not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get boilerplate: %w", err)
	}

	return bp, nil
}

func (r *boilerplateRepository) GetByProblemID(problemID int) ([]domain.ProblemBoilerplate, error) {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	var boilerplates []domain.ProblemBoilerplate
	err := r.db.DB.WithContext(ctx).
		Where("problem_id = ?", problemID).
		Preload("Language").
		Find(&boilerplates).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get boilerplates for problem: %w", err)
	}

	return boilerplates, nil
}

func (r *boilerplateRepository) Update(boilerplate *domain.ProblemBoilerplate) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Save(boilerplate)
	if result.Error != nil {
		return fmt.Errorf("failed to update boilerplate: %w", result.Error)
	}

	return nil
}

func (r *boilerplateRepository) DeleteByProblemID(problemID int) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).
		Where("problem_id = ?", problemID).
		Delete(&domain.ProblemBoilerplate{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete boilerplates: %w", result.Error)
	}

	return nil
}

func (r *boilerplateRepository) Exists(problemID, languageID int) (bool, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	var count int64
	err := r.db.DB.WithContext(ctx).
		Model(&domain.ProblemBoilerplate{}).
		Where("problem_id = ? AND language_id = ?", problemID, languageID).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check boilerplate existence: %w", err)
	}

	return count > 0, nil
}
