package interfaces

import (
	"github.com/prabalesh/loco/backend/internal/domain"
)

// TestCaseRepository defines the interface for test case data access
type TestCaseRepository interface {
	// Create a new test case for a problem
	Create(testCase *domain.TestCase) error

	// Update an existing test case
	Update(testCase *domain.TestCase) error

	// Delete a test case by ID
	Delete(id int) error

	// Get test case by ID
	GetByID(id int) (*domain.TestCase, error)

	// Get all test cases for a specific problem
	GetByProblemID(problemID int, includeSamplesOnly bool) ([]*domain.TestCase, error)

	// List test cases with pagination and filters
	List(problemID int, filters TestCaseFilters) ([]*domain.TestCase, int, error)

	// Check if test case exists by ID
	Exists(id int) (bool, error)

	// Update test case order index (for reordering)
	UpdateOrderIndex(id int, orderIndex int) error

	// Delete all test cases for a problem (cascade delete)
	DeleteByProblemID(problemID int) error

	// Count test cases for a problem
	CountByProblemID(problemID int) (int, error)

	// Get sample test cases only for a problem
	GetSamples(problemID int) ([]*domain.TestCase, error)
}

// TestCaseFilters for listing test cases
type TestCaseFilters struct {
	IsSample *bool
	Limit    int
	Page     int
}
