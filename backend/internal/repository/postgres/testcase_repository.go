package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/usecase/interfaces"
	"github.com/prabalesh/loco/backend/pkg/database"
)

type testCaseRepository struct {
	db *database.Database
}

func NewTestCaseRepository(db *database.Database) *testCaseRepository {
	return &testCaseRepository{db: db}
}

func (r *testCaseRepository) Create(testCase *domain.TestCase) error {
	ctx, cancel := database.WithLongTimeout()
	defer cancel()

	query := `
        INSERT INTO test_cases (problem_id, input, expected_output, is_sample, 
                               validation_config, order_index, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, NOW())
        RETURNING id
    `

	err := r.db.QueryRowContext(ctx, query,
		testCase.ProblemID,
		testCase.Input,
		testCase.ExpectedOutput,
		testCase.IsSample,
		testCase.ValidationConfig,
		testCase.OrderIndex,
	).Scan(&testCase.ID)

	if err != nil {
		return fmt.Errorf("failed to create test case: %w", err)
	}

	return nil
}

func (r *testCaseRepository) Update(testCase *domain.TestCase) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	query := `
        UPDATE test_cases 
        SET input = $1,
            expected_output = $2,
            is_sample = $3,
            validation_config = $4,
            order_index = $5,
            updated_at = NOW()
        WHERE id = $6
        RETURNING updated_at
    `

	var updatedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query,
		testCase.Input,
		testCase.ExpectedOutput,
		testCase.IsSample,
		testCase.ValidationConfig,
		testCase.OrderIndex,
		testCase.ID,
	).Scan(&updatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("test case not found")
	}

	if err != nil {
		return fmt.Errorf("failed to update test case: %w", err)
	}

	return nil
}

func (r *testCaseRepository) Delete(id int) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	query := `DELETE FROM test_cases WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete test case: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("test case not found")
	}

	return nil
}

func (r *testCaseRepository) GetByID(id int) (*domain.TestCase, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	testCase := &domain.TestCase{}
	query := `
        SELECT id, problem_id, input, expected_output, is_sample, 
               validation_config, order_index, created_at
        FROM test_cases WHERE id = $1
    `

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&testCase.ID,
		&testCase.ProblemID,
		&testCase.Input,
		&testCase.ExpectedOutput,
		&testCase.IsSample,
		&testCase.ValidationConfig,
		&testCase.OrderIndex,
		&testCase.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("test case not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get test case: %w", err)
	}

	return testCase, nil
}

func (r *testCaseRepository) GetByProblemID(problemID int, includeSamplesOnly bool) ([]*domain.TestCase, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	query := `
        SELECT id, problem_id, input, expected_output, is_sample, 
               validation_config, order_index, created_at
        FROM test_cases 
        WHERE problem_id = $1
    `
	args := []interface{}{problemID}

	if includeSamplesOnly {
		query += ` AND is_sample = true`
	}

	query += ` ORDER BY order_index ASC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get test cases: %w", err)
	}
	defer rows.Close()

	var testCases []*domain.TestCase
	for rows.Next() {
		testCase := &domain.TestCase{}
		err := rows.Scan(
			&testCase.ID,
			&testCase.ProblemID,
			&testCase.Input,
			&testCase.ExpectedOutput,
			&testCase.IsSample,
			&testCase.ValidationConfig,
			&testCase.OrderIndex,
			&testCase.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan test case: %w", err)
		}
		testCases = append(testCases, testCase)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating test cases: %w", err)
	}

	return testCases, nil
}

func (r *testCaseRepository) List(problemID int, filters interfaces.TestCaseFilters) ([]*domain.TestCase, int, error) {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	whereClauses := []string{fmt.Sprintf("problem_id = $%d", 1)}
	args := []interface{}{problemID}
	argPosition := 2

	if filters.IsSample != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("is_sample = $%d", argPosition))
		args = append(args, *filters.IsSample)
		argPosition++
	}

	whereClause := "WHERE " + strings.Join(whereClauses, " AND ")

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM test_cases %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count test cases: %w", err)
	}

	// Pagination
	limit := filters.Limit
	if limit == 0 {
		limit = 50
	}
	offset := (filters.Page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	// Fetch test cases
	query := fmt.Sprintf(`
        SELECT id, problem_id, input, expected_output, is_sample, 
               validation_config, order_index, created_at
        FROM test_cases
        %s
        ORDER BY order_index ASC
        LIMIT $%d OFFSET $%d
    `, whereClause, argPosition, argPosition+1)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list test cases: %w", err)
	}
	defer rows.Close()

	var testCases []*domain.TestCase
	for rows.Next() {
		testCase := &domain.TestCase{}
		err := rows.Scan(
			&testCase.ID,
			&testCase.ProblemID,
			&testCase.Input,
			&testCase.ExpectedOutput,
			&testCase.IsSample,
			&testCase.ValidationConfig,
			&testCase.OrderIndex,
			&testCase.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan test case: %w", err)
		}
		testCases = append(testCases, testCase)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating test cases: %w", err)
	}

	return testCases, total, nil
}

func (r *testCaseRepository) Exists(id int) (bool, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM test_cases WHERE id = $1)`

	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check test case existence: %w", err)
	}

	return exists, nil
}

func (r *testCaseRepository) UpdateOrderIndex(id int, orderIndex int) error {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	query := `
        UPDATE test_cases 
        SET order_index = $1, updated_at = NOW()
        WHERE id = $2
    `

	result, err := r.db.ExecContext(ctx, query, orderIndex, id)
	if err != nil {
		return fmt.Errorf("failed to update order index: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("test case not found")
	}

	return nil
}

func (r *testCaseRepository) DeleteByProblemID(problemID int) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	query := `DELETE FROM test_cases WHERE problem_id = $1`

	_, err := r.db.ExecContext(ctx, query, problemID)
	if err != nil {
		return fmt.Errorf("failed to delete test cases: %w", err)
	}

	return nil
}

func (r *testCaseRepository) CountByProblemID(problemID int) (int, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	var count int
	query := `SELECT COUNT(*) FROM test_cases WHERE problem_id = $1`

	err := r.db.QueryRowContext(ctx, query, problemID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count test cases: %w", err)
	}

	return count, nil
}

func (r *testCaseRepository) GetSamples(problemID int) ([]*domain.TestCase, error) {
	return r.GetByProblemID(problemID, true)
}
