package postgres

import (
	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/pkg/database"
)

type pistonExecutionRepository struct {
	db *database.Database
}

func NewPistonExecutionRepository(db *database.Database) domain.PistonExecutionRepository {
	return &pistonExecutionRepository{db: db}
}

func (r *pistonExecutionRepository) Create(execution *domain.PistonExecution) error {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	return r.db.DB.WithContext(ctx).Create(execution).Error
}

func (r *pistonExecutionRepository) List(limit, offset int) ([]domain.PistonExecution, int64, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	var executions []domain.PistonExecution
	var total int64

	err := r.db.DB.WithContext(ctx).Model(&domain.PistonExecution{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.DB.WithContext(ctx).
		Preload("Problem").
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&executions).Error

	return executions, total, err
}

func (r *pistonExecutionRepository) GetByProblemID(problemID int, limit, offset int) ([]domain.PistonExecution, int64, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	var executions []domain.PistonExecution
	var total int64

	query := r.db.DB.WithContext(ctx).Model(&domain.PistonExecution{}).Where("problem_id = ?", problemID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Problem").
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&executions).Error

	return executions, total, err
}
