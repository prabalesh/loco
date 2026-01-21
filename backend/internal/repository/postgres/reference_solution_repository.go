package postgres

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/pkg/database"
)

type referenceSolutionRepository struct {
	db *database.Database
}

// NewReferenceSolutionRepository creates a new postgres reference solution repository
func NewReferenceSolutionRepository(db *database.Database) domain.ReferenceSolutionRepository {
	return &referenceSolutionRepository{db: db}
}

func (r *referenceSolutionRepository) Create(solution *domain.ProblemReferenceSolution) error {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	if err := r.db.DB.WithContext(ctx).Create(solution).Error; err != nil {
		return fmt.Errorf("failed to create reference solution: %w", err)
	}
	return nil
}

func (r *referenceSolutionRepository) GetByProblemAndLanguage(problemID, languageID int) (*domain.ProblemReferenceSolution, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	var solution domain.ProblemReferenceSolution
	err := r.db.DB.WithContext(ctx).
		Where("problem_id = ? AND language_id = ?", problemID, languageID).
		Preload("Language").
		First(&solution).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get reference solution: %w", err)
	}
	return &solution, nil
}

func (r *referenceSolutionRepository) GetAllByProblemID(problemID int) ([]domain.ProblemReferenceSolution, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	var solutions []domain.ProblemReferenceSolution
	err := r.db.DB.WithContext(ctx).
		Where("problem_id = ?", problemID).
		Preload("Language").
		Find(&solutions).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get reference solutions: %w", err)
	}
	return solutions, nil
}

func (r *referenceSolutionRepository) Update(solution *domain.ProblemReferenceSolution) error {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	if err := r.db.DB.WithContext(ctx).Save(solution).Error; err != nil {
		return fmt.Errorf("failed to update reference solution: %w", err)
	}
	return nil
}

func (r *referenceSolutionRepository) Delete(id int) error {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	if err := r.db.DB.WithContext(ctx).Delete(&domain.ProblemReferenceSolution{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete reference solution: %w", err)
	}
	return nil
}

func (r *referenceSolutionRepository) Exists(problemID, languageID int) (bool, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	var count int64
	err := r.db.DB.WithContext(ctx).Model(&domain.ProblemReferenceSolution{}).
		Where("problem_id = ? AND language_id = ?", problemID, languageID).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check reference solution existence: %w", err)
	}
	return count > 0, nil
}
