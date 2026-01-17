package usecase

import (
	"errors"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/domain/uerror"
	"github.com/prabalesh/loco/backend/internal/domain/validator"

	"github.com/prabalesh/loco/backend/pkg/config"
	"go.uber.org/zap"
)

type TestCaseUsecase struct {
	testCaseRepo domain.TestCaseRepository
	problemRepo  domain.ProblemRepository
	cfg          *config.Config
	logger       *zap.Logger
}

func NewTestCaseUsecase(
	testCaseRepo domain.TestCaseRepository,
	problemRepo domain.ProblemRepository,
	cfg *config.Config,
	logger *zap.Logger,
) *TestCaseUsecase {
	return &TestCaseUsecase{
		testCaseRepo: testCaseRepo,
		problemRepo:  problemRepo,
		cfg:          cfg,
		logger:       logger,
	}
}

// ========== ADMIN OPERATIONS ==========

// CreateTestCase creates a new test case for a problem
func (u *TestCaseUsecase) CreateTestCase(req *domain.CreateTestCaseRequest, adminID int) (*domain.TestCase, error) {
	// Validate request (custom validation - no external package)
	if validationErrors := validator.ValidateCreateTestCaseRequest(req); len(validationErrors) > 0 {
		u.logger.Warn("Create test case validation failed", zap.Any("errors", validationErrors))
		return nil, &uerror.ValidationError{Errors: validationErrors}
	}

	// Verify problem exists
	_, err := u.problemRepo.GetByID(req.ProblemID)
	if err != nil {
		u.logger.Warn("Problem not found for test case creation", zap.Int("problem_id", req.ProblemID))
		return nil, errors.New("problem not found")
	}

	// Set defaults
	testCase := &domain.TestCase{
		ProblemID:        req.ProblemID,
		Input:            req.Input,
		ExpectedOutput:   req.ExpectedOutput,
		IsSample:         req.IsSample,
		ValidationConfig: req.ValidationConfig,
		OrderIndex:       req.OrderIndex,
	}

	if err := u.testCaseRepo.Create(testCase); err != nil {
		u.logger.Error("Failed to create test case",
			zap.Error(err),
			zap.Int("problem_id", req.ProblemID),
			zap.Int("admin_id", adminID),
		)
		return nil, errors.New("failed to create test case")
	}

	u.logger.Info("Test case created successfully",
		zap.Int("test_case_id", testCase.ID),
		zap.Int("problem_id", testCase.ProblemID),
		zap.Int("admin_id", adminID),
	)

	return testCase, nil
}

// UpdateTestCase updates an existing test case
func (u *TestCaseUsecase) UpdateTestCase(testCaseID int, req *domain.UpdateTestCaseRequest, adminID int) (*domain.TestCase, error) {
	// Validate request
	if validationErrors := validator.ValidateUpdateTestCaseRequest(req); len(validationErrors) > 0 {
		u.logger.Warn("Update test case validation failed", zap.Any("errors", validationErrors))
		return nil, &uerror.ValidationError{Errors: validationErrors}
	}

	// Get existing test case
	testCase, err := u.testCaseRepo.GetByID(testCaseID)
	if err != nil {
		u.logger.Warn("Test case not found", zap.Int("test_case_id", testCaseID))
		return nil, errors.New("test case not found")
	}

	// Update fields if provided
	if req.Input != "" {
		testCase.Input = req.Input
	}
	if req.ExpectedOutput != "" {
		testCase.ExpectedOutput = req.ExpectedOutput
	}
	if req.IsSample != nil {
		testCase.IsSample = *req.IsSample
	}
	if len(req.ValidationConfig) > 0 {
		testCase.ValidationConfig = req.ValidationConfig
	}
	if req.OrderIndex != nil {
		testCase.OrderIndex = *req.OrderIndex
	}

	if err := u.testCaseRepo.Update(testCase); err != nil {
		u.logger.Error("Failed to update test case",
			zap.Error(err),
			zap.Int("test_case_id", testCaseID),
		)
		return nil, errors.New("failed to update test case")
	}

	u.logger.Info("Test case updated successfully",
		zap.Int("test_case_id", testCase.ID),
		zap.Int("admin_id", adminID),
	)

	return testCase, nil
}

// DeleteTestCase deletes a specific test case
func (u *TestCaseUsecase) DeleteTestCase(testCaseID int, adminID int) error {
	// Verify test case exists
	_, err := u.testCaseRepo.GetByID(testCaseID)
	if err != nil {
		u.logger.Warn("Test case not found for deletion", zap.Int("test_case_id", testCaseID))
		return errors.New("test case not found")
	}

	if err := u.testCaseRepo.Delete(testCaseID); err != nil {
		u.logger.Error("Failed to delete test case",
			zap.Error(err),
			zap.Int("test_case_id", testCaseID),
		)
		return errors.New("failed to delete test case")
	}

	u.logger.Info("Test case deleted successfully",
		zap.Int("test_case_id", testCaseID),
		zap.Int("admin_id", adminID),
	)

	return nil
}

// DeleteAllTestCases deletes all test cases for a problem
func (u *TestCaseUsecase) DeleteAllTestCases(problemID int, adminID int) error {
	// Verify problem exists
	_, err := u.problemRepo.GetByID(problemID)
	if err != nil {
		u.logger.Warn("Problem not found for test case deletion", zap.Int("problem_id", problemID))
		return errors.New("problem not found")
	}

	if err := u.testCaseRepo.DeleteByProblemID(problemID); err != nil {
		u.logger.Error("Failed to delete all test cases",
			zap.Error(err),
			zap.Int("problem_id", problemID),
		)
		return errors.New("failed to delete test cases")
	}

	u.logger.Info("All test cases deleted successfully",
		zap.Int("problem_id", problemID),
		zap.Int("admin_id", adminID),
	)

	return nil
}

// ReorderTestCases reorders test cases for a problem
func (u *TestCaseUsecase) ReorderTestCases(req *domain.ReorderTestCasesRequest, adminID int) error {
	// Custom validation
	if validationErrors := validator.ValidateReorderTestCasesRequest(req); len(validationErrors) > 0 {
		u.logger.Warn("Reorder test cases validation failed", zap.Any("errors", validationErrors))
		return &uerror.ValidationError{Errors: validationErrors}
	}

	// Verify all test cases belong to the same problem
	for _, tc := range req.TestCases {
		testCase, err := u.testCaseRepo.GetByID(tc.ID)
		if err != nil {
			return errors.New("test case not found")
		}
		if testCase.ProblemID != req.ProblemID {
			return errors.New("test case does not belong to specified problem")
		}
	}

	// Update order indices
	for _, tc := range req.TestCases {
		if err := u.testCaseRepo.UpdateOrderIndex(tc.ID, tc.OrderIndex); err != nil {
			u.logger.Error("Failed to update test case order",
				zap.Error(err),
				zap.Int("test_case_id", tc.ID),
			)
			return errors.New("failed to reorder test cases")
		}
	}

	u.logger.Info("Test cases reordered successfully",
		zap.Int("problem_id", req.ProblemID),
		zap.Int("count", len(req.TestCases)),
		zap.Int("admin_id", adminID),
	)

	return nil
}

// ========== USER/GET OPERATIONS ==========

// GetTestCase retrieves a specific test case by ID
func (u *TestCaseUsecase) GetTestCase(testCaseID int) (*domain.TestCase, error) {
	testCase, err := u.testCaseRepo.GetByID(testCaseID)
	if err != nil {
		u.logger.Warn("Test case not found", zap.Int("test_case_id", testCaseID))
		return nil, errors.New("test case not found")
	}
	return testCase, nil
}

// GetTestCasesByProblem gets all test cases for a problem
func (u *TestCaseUsecase) GetTestCasesByProblem(problemID int, includeSamplesOnly bool) ([]*domain.TestCase, error) {
	// Verify problem exists and is published/public for user access
	problem, err := u.problemRepo.GetByID(problemID)
	if err != nil {
		return nil, errors.New("problem not found")
	}
	if problem.Status != "published" || problem.Visibility != "public" {
		return nil, errors.New("problem not accessible")
	}

	var testCases []domain.TestCase

	if includeSamplesOnly {
		testCases, err = u.testCaseRepo.GetSamples(problemID)
	} else {
		testCases, err = u.testCaseRepo.GetByProblemID(problemID)
	}

	if err != nil {
		u.logger.Error("Failed to get test cases",
			zap.Error(err),
			zap.Int("problem_id", problemID),
		)
		return nil, errors.New("failed to retrieve test cases")
	}

	result := make([]*domain.TestCase, len(testCases))
	for i := range testCases {
		result[i] = &testCases[i]
	}

	return result, nil
}

// ListTestCases lists test cases with pagination
func (u *TestCaseUsecase) ListTestCases(req *domain.ListTestCasesRequest) ([]*domain.TestCase, int, error) {
	if req.ProblemID <= 0 {
		return nil, 0, errors.New("problem_id is required")
	}

	filters := domain.TestCaseFilters{
		IsSample: req.IsSample,
		Limit:    req.Limit,
		Page:     req.Page,
	}

	testCases, total, err := u.testCaseRepo.List(req.ProblemID, filters)
	if err != nil {
		u.logger.Error("Failed to list test cases",
			zap.Error(err),
			zap.Int("problem_id", req.ProblemID),
		)
		return nil, 0, errors.New("failed to retrieve test cases")
	}
	return testCases, total, nil
}

// GetSampleTestCases gets only sample test cases
func (u *TestCaseUsecase) GetSampleTestCases(problemID int) ([]*domain.TestCase, error) {
	return u.GetTestCasesByProblem(problemID, true)
}

// CountTestCasesByProblem returns count of test cases for a problem
func (u *TestCaseUsecase) CountTestCasesByProblem(problemID int) (int, error) {
	count, err := u.testCaseRepo.CountByProblemID(problemID)
	if err != nil {
		u.logger.Error("Failed to count test cases",
			zap.Error(err),
			zap.Int("problem_id", problemID),
		)
		return 0, errors.New("failed to count test cases")
	}
	return count, nil
}
