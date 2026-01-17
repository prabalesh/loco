package postgres

import (
	"fmt"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/pkg/database"
	"gorm.io/gorm"
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

	result := r.db.DB.WithContext(ctx).Create(testCase)
	if result.Error != nil {
		return fmt.Errorf("failed to create test case: %w", result.Error)
	}

	return nil
}

func (r *testCaseRepository) Update(testCase *domain.TestCase) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	// Update specific fields or all fields?
	// Using Updates with struct updates non-zero fields.
	// If input/output can be empty string, we might need map or Select.
	// Assuming they are provided.

	result := r.db.DB.WithContext(ctx).Model(&domain.TestCase{}).Where("id = ?", testCase.ID).Updates(testCase)

	if result.Error != nil {
		return fmt.Errorf("failed to update test case: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("test case not found")
	}

	return nil
}

func (r *testCaseRepository) Delete(id int) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Delete(&domain.TestCase{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete test case: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("test case not found")
	}

	return nil
}

func (r *testCaseRepository) GetByID(id int) (*domain.TestCase, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	testCase := &domain.TestCase{}
	err := r.db.DB.WithContext(ctx).First(testCase, id).Error

	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("test case not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get test case: %w", err)
	}

	return testCase, nil
}

func (r *testCaseRepository) GetByProblemID(problemID int) ([]domain.TestCase, error) {
	var testCases []domain.TestCase
	err := r.db.DB.Where("problem_id = ?", problemID).Order("order_index ASC").Find(&testCases).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get test cases: %w", err)
	}
	return testCases, nil
}

func (r *testCaseRepository) List(problemID int, filters domain.TestCaseFilters) ([]*domain.TestCase, int, error) {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	var testCases []*domain.TestCase
	var total int64

	query := r.db.DB.WithContext(ctx).Model(&domain.TestCase{}).Where("problem_id = ?", problemID)

	if filters.IsSample != nil {
		query = query.Where("is_sample = ?", *filters.IsSample)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
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
	if err := query.Order("order_index ASC").Limit(limit).Offset(offset).Find(&testCases).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list test cases: %w", err)
	}

	return testCases, int(total), nil
}

func (r *testCaseRepository) Exists(id int) (bool, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	var count int64
	err := r.db.DB.WithContext(ctx).Model(&domain.TestCase{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check test case existence: %w", err)
	}

	return count > 0, nil
}

func (r *testCaseRepository) UpdateOrderIndex(id int, orderIndex int) error {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Model(&domain.TestCase{}).Where("id = ?", id).Update("order_index", orderIndex)
	if result.Error != nil {
		return fmt.Errorf("failed to update order index: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("test case not found")
	}

	return nil
}

func (r *testCaseRepository) DeleteByProblemID(problemID int) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Where("problem_id = ?", problemID).Delete(&domain.TestCase{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete test cases: %w", result.Error)
	}

	return nil
}

func (r *testCaseRepository) CountByProblemID(problemID int) (int, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	var count int64
	err := r.db.DB.WithContext(ctx).Model(&domain.TestCase{}).Where("problem_id = ?", problemID).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count test cases: %w", err)
	}

	return int(count), nil
}

func (r *testCaseRepository) GetSamples(problemID int) ([]domain.TestCase, error) {
	var testCases []domain.TestCase
	err := r.db.DB.Where("problem_id = ? AND is_sample = ?", problemID, true).Order("order_index ASC").Find(&testCases).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get sample test cases: %w", err)
	}
	return testCases, nil
}
