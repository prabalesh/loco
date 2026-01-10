package postgres

import (
	"fmt"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/pkg/database"
)

type problemLanguageRepository struct {
	db *database.Database
}

func NewProblemLanguageRepository(db *database.Database) *problemRepository {
	return &problemRepository{db: db}
}

func (r *problemLanguageRepository) Create(problemLanguage *domain.ProblemLanguage) error {
	ctx, cancel := database.WithLongTimeout()
	defer cancel()

	query := `
		INSERT INTO problem_language (problem_id, language_id, function_code, main_code, solution_code)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		problemLanguage.ProblemID,
		problemLanguage.LanguageID,
		problemLanguage.FunctionCode,
		problemLanguage.MainCode,
		problemLanguage.SolutionCode,
	).Scan(&problemLanguage.CreatedAt, &problemLanguage.UpdatedAt)

	if err != nil {
		if isUniqueViolation(err) {
			if containsField(err, "problem_id") {
				return fmt.Errorf("problem_id already exists")
			}
			if containsField(err, "language_id") {
				return fmt.Errorf("language_id already exists")
			}
		}
		return fmt.Errorf("failed to create problem language: %w", err)
	}

	return nil
}
