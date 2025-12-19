package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/usecase/interfaces"
	"github.com/prabalesh/loco/backend/pkg/database"
)

type problemRepository struct {
	db *database.Database
}

func NewProblemRepository(db *database.Database) *problemRepository {
	return &problemRepository{db: db}
}

func (r *problemRepository) Create(problem *domain.Problem) error {
	ctx, cancel := database.WithLongTimeout()
	defer cancel()

	query := `
		INSERT INTO problems (title, description, slug, difficulty, time_limit, memory_limit, validator_type, input_format, output_format, constraints, status, visibility, is_active, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		problem.Title,
		problem.Description,
		problem.Slug,
		problem.Difficulty,
		problem.TimeLimit,
		problem.MemoryLimit,
		problem.ValidatorType,
		problem.InputFormat,
		problem.OutputFormat,
		problem.Constraints,
		problem.Status,
		problem.Visibility,
		problem.IsActive,
		problem.CreatedBy,
	).Scan(&problem.ID, &problem.CreatedAt, &problem.UpdatedAt)

	if err != nil {
		if isUniqueViolation(err) {
			if containsField(err, "title") {
				return fmt.Errorf("title already exists")
			}
			if containsField(err, "slug") {
				return fmt.Errorf("slug already exists")
			}
		}
		return fmt.Errorf("failed to create problem: %w", err)
	}

	return nil
}

func (r *problemRepository) Update(problem *domain.Problem) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	query := `
		UPDATE problems SET
            title = $1,
            slug = $2,
            description = $3,
            difficulty = $4,
            time_limit = $5,
            memory_limit = $6,
            validator_type = $7,
            input_format = $8,
            output_format = $9,
            constraints = $10,
            status = $11,
            visibility = $12,
            is_active = $13,
            updated_at = NOW()
        WHERE id = $14
        RETURNING updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		problem.Title,
		problem.Slug,
		problem.Description,
		problem.Difficulty,
		problem.TimeLimit,
		problem.MemoryLimit,
		problem.ValidatorType,
		problem.InputFormat,
		problem.OutputFormat,
		problem.Constraints,
		problem.Status,
		problem.Visibility,
		problem.IsActive,
		problem.ID,
	).Scan(&problem.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("problem not found")
	}

	if err != nil {
		return fmt.Errorf("failed to update problem: %w", err)
	}

	return nil
}

func (r *problemRepository) Delete(id int) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	query := `DELETE FROM problems WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return fmt.Errorf("failed to delete problem: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("problem not found")
	}

	return nil
}

func (r *problemRepository) GetByID(id int) (*domain.Problem, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	problem := &domain.Problem{}
	var createdBy sql.NullInt64

	query := `
        SELECT id, title, slug, description, difficulty,
               time_limit, memory_limit, validator_type,
               input_format, output_format, constraints,
               status, visibility, is_active,
               acceptance_rate, total_submissions, total_accepted,
               created_by, created_at, updated_at
        FROM problems WHERE id = $1
    `

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&problem.ID,
		&problem.Title,
		&problem.Slug,
		&problem.Description,
		&problem.Difficulty,
		&problem.TimeLimit,
		&problem.MemoryLimit,
		&problem.ValidatorType,
		&problem.InputFormat,
		&problem.OutputFormat,
		&problem.Constraints,
		&problem.Status,
		&problem.Visibility,
		&problem.IsActive,
		&problem.AcceptanceRate,
		&problem.TotalSubmission,
		&problem.TotalAccepted,
		&createdBy,
		&problem.CreatedAt,
		&problem.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("problem not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get problem: %w", err)
	}

	// Handle nullable created_by
	if createdBy.Valid {
		val := int(createdBy.Int64)
		problem.CreatedBy = &val
	}

	return problem, nil
}

func (r *problemRepository) GetBySlug(slug string) (*domain.Problem, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	problem := &domain.Problem{}
	var createdBy sql.NullInt64

	query := `
        SELECT id, title, slug, description, difficulty,
               time_limit, memory_limit, validator_type,
               input_format, output_format, constraints,
               status, visibility, is_active,
               acceptance_rate, total_submissions, total_accepted,
               created_by, created_at, updated_at
        FROM problems WHERE slug = $1
    `

	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&problem.ID,
		&problem.Title,
		&problem.Slug,
		&problem.Description,
		&problem.Difficulty,
		&problem.TimeLimit,
		&problem.MemoryLimit,
		&problem.ValidatorType,
		&problem.InputFormat,
		&problem.OutputFormat,
		&problem.Constraints,
		&problem.Status,
		&problem.Visibility,
		&problem.IsActive,
		&problem.AcceptanceRate,
		&problem.TotalSubmission,
		&problem.TotalAccepted,
		&createdBy,
		&problem.CreatedAt,
		&problem.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("problem not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get problem: %w", err)
	}

	// Handle nullable created_by
	if createdBy.Valid {
		val := int(createdBy.Int64)
		problem.CreatedBy = &val
	}

	return problem, nil
}

func (r *problemRepository) List(filters interfaces.ProblemFilters) ([]*domain.Problem, int, error) {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	// Build WHERE clause dynamically
	whereClauses := []string{}
	args := []interface{}{}
	argPosition := 1

	if filters.Difficulty != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("difficulty = $%d", argPosition))
		args = append(args, filters.Difficulty)
		argPosition++
	}

	if filters.Status != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("status = $%d", argPosition))
		args = append(args, filters.Status)
		argPosition++
	}

	if filters.Visibility != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("visibility = $%d", argPosition))
		args = append(args, filters.Visibility)
		argPosition++
	}

	if filters.Search != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d)", argPosition, argPosition))
		searchPattern := "%" + filters.Search + "%"
		args = append(args, searchPattern)
		argPosition++
	}

	if filters.CreatedBy != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("created_by = $%d", argPosition))
		args = append(args, *filters.CreatedBy)
		argPosition++
	}

	// Tag filtering (if tags provided)
	if len(filters.Tags) > 0 {
		whereClauses = append(whereClauses, fmt.Sprintf(`
            id IN (
                SELECT problem_id FROM problem_tags pt
                JOIN tags t ON pt.tag_id = t.id
                WHERE t.slug = ANY($%d)
            )
        `, argPosition))
		args = append(args, filters.Tags)
		argPosition++
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM problems %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count problems: %w", err)
	}

	// Pagination
	limit := filters.Limit
	if limit == 0 {
		limit = 20 // Default
	}
	offset := (filters.Page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	// Fetch problems
	query := fmt.Sprintf(`
        SELECT id, title, slug, description, difficulty,
               time_limit, memory_limit, current_step, validator_type,
               status, visibility, is_active,
               acceptance_rate, total_submissions, total_accepted,
               created_at, updated_at
        FROM problems
        %s
        ORDER BY created_at DESC
        LIMIT $%d OFFSET $%d
    `, whereClause, argPosition, argPosition+1)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list problems: %w", err)
	}
	defer rows.Close()

	var problems []*domain.Problem
	for rows.Next() {
		var problem domain.Problem
		err := rows.Scan(
			&problem.ID,
			&problem.Title,
			&problem.Slug,
			&problem.Description,
			&problem.Difficulty,
			&problem.TimeLimit,
			&problem.MemoryLimit,
			&problem.CurrentStep,
			&problem.ValidatorType,
			&problem.Status,
			&problem.Visibility,
			&problem.IsActive,
			&problem.AcceptanceRate,
			&problem.TotalSubmission,
			&problem.TotalAccepted,
			&problem.CreatedAt,
			&problem.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan problem: %w", err)
		}
		problems = append(problems, &problem)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating problems: %w", err)
	}

	return problems, total, nil
}

func (r *problemRepository) SlugExists(slug string) (bool, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM problems WHERE slug = $1)`

	err := r.db.QueryRowContext(ctx, query, slug).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check slug: %w", err)
	}

	return exists, nil
}

func (r *problemRepository) TitleExists(title string) (bool, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM problems WHERE title = $1)`

	err := r.db.QueryRowContext(ctx, query, title).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check title: %w", err)
	}

	return exists, nil
}

func (r *problemRepository) UpdateCurrentStep(id int, newCurrentStep int) error {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	query := `
		UPDATE problems
		SET current_step = $1,
		updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, newCurrentStep, id)
	if err != nil {
		return fmt.Errorf("failed to update current step status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("problem not found")
	}

	return nil
}

func (r *problemRepository) UpdateStats(id int, acceptanceRate float64, totalSubmissions, totalAccepted int) error {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	query := `
        UPDATE problems 
        SET acceptance_rate = $1,
            total_submissions = $2,
            total_accepted = $3,
            updated_at = NOW()
        WHERE id = $4
    `

	result, err := r.db.ExecContext(ctx, query, acceptanceRate, totalSubmissions, totalAccepted, id)
	if err != nil {
		return fmt.Errorf("failed to update stats: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("problem not found")
	}

	return nil
}

func (r *problemRepository) UpdateStatus(id int, status string) error {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	query := `
        UPDATE problems 
        SET status = $1, updated_at = NOW() 
        WHERE id = $2
    `

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("problem not found")
	}

	return nil
}

func (r *problemRepository) UpdateVisibility(id int, visibility string) error {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	query := `
        UPDATE problems 
        SET visibility = $1, updated_at = NOW() 
        WHERE id = $2
    `

	result, err := r.db.ExecContext(ctx, query, visibility, id)
	if err != nil {
		return fmt.Errorf("failed to update visibility: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("problem not found")
	}

	return nil
}

func (r *problemRepository) CountProblems() (int, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	var count int
	query := `SELECT COUNT(*) FROM problems`

	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count problems: %w", err)
	}

	return count, nil
}

func (r *problemRepository) CountByStatus(status string) (int, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	var count int
	query := `SELECT COUNT(*) FROM problems WHERE status = $1`

	err := r.db.QueryRowContext(ctx, query, status).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count problems by status: %w", err)
	}

	return count, nil
}

func (r *problemRepository) CountByDifficulty(difficulty string) (int, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	var count int
	query := `SELECT COUNT(*) FROM problems WHERE difficulty = $1`

	err := r.db.QueryRowContext(ctx, query, difficulty).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count problems by difficulty: %w", err)
	}

	return count, nil
}
