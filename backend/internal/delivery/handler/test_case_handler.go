package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/prabalesh/loco/backend/internal/delivery/middleware"
	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/domain/uerror"
	"github.com/prabalesh/loco/backend/internal/usecase"
	"github.com/prabalesh/loco/backend/pkg/config"
	"go.uber.org/zap"
)

type TestCaseHandler struct {
	testCaseUsecase *usecase.TestCaseUsecase
	logger          *zap.Logger
	cfg             *config.Config
}

func NewTestCaseHandler(testCaseUsecase *usecase.TestCaseUsecase, logger *zap.Logger, cfg *config.Config) *TestCaseHandler {
	return &TestCaseHandler{
		testCaseUsecase: testCaseUsecase,
		logger:          logger,
		cfg:             cfg,
	}
}

// ========== ADMIN ENDPOINTS ==========

// CreateTestCase creates a new test case (admin only)
func (h *TestCaseHandler) CreateTestCase(w http.ResponseWriter, r *http.Request) {
	// Get admin ID from context
	adminID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse request
	var req domain.CreateTestCaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Invalid JSON in create test case request", zap.Error(err))
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Create test case
	testCase, err := h.testCaseUsecase.CreateTestCase(&req, adminID)
	if err != nil {
		// Handle validation errors
		var validationErr *uerror.ValidationError
		if errors.As(err, &validationErr) {
			h.logger.Warn("Create test case validation failed",
				zap.Any("errors", validationErr.Errors),
			)
			RespondValidationError(w, validationErr.Errors)
			return
		}

		// Handle business logic errors
		errMsg := err.Error()
		switch errMsg {
		case "problem not found":
			RespondError(w, http.StatusNotFound, errMsg)
		default:
			h.logger.Error("Test case creation failed", zap.Error(err))
			RespondError(w, http.StatusInternalServerError, "failed to create test case")
		}
		return
	}

	h.logger.Info("Test case created successfully",
		zap.Int("test_case_id", testCase.ID),
		zap.Int("problem_id", testCase.ProblemID),
		zap.Int("admin_id", adminID),
	)

	RespondJSON(w, http.StatusCreated, testCase)
}

// UpdateTestCase updates an existing test case (admin only)
func (h *TestCaseHandler) UpdateTestCase(w http.ResponseWriter, r *http.Request) {
	// Get admin ID from context
	adminID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get test case ID from path
	testCaseID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid test case ID")
		return
	}

	// Parse request
	var req domain.UpdateTestCaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Invalid JSON in update test case request", zap.Error(err))
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Update test case
	testCase, err := h.testCaseUsecase.UpdateTestCase(testCaseID, &req, adminID)
	if err != nil {
		// Handle validation errors
		var validationErr *uerror.ValidationError
		if errors.As(err, &validationErr) {
			h.logger.Warn("Update test case validation failed",
				zap.Any("errors", validationErr.Errors),
			)
			RespondValidationError(w, validationErr.Errors)
			return
		}

		errMsg := err.Error()
		switch errMsg {
		case "test case not found":
			RespondError(w, http.StatusNotFound, errMsg)
		default:
			h.logger.Error("Test case update failed", zap.Error(err))
			RespondError(w, http.StatusInternalServerError, "failed to update test case")
		}
		return
	}

	h.logger.Info("Test case updated successfully",
		zap.Int("test_case_id", testCase.ID),
		zap.Int("admin_id", adminID),
	)

	RespondJSON(w, http.StatusOK, testCase)
}

// DeleteTestCase deletes a specific test case (admin only)
func (h *TestCaseHandler) DeleteTestCase(w http.ResponseWriter, r *http.Request) {
	// Get admin ID from context
	adminID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get test case ID from path
	testCaseID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid test case ID")
		return
	}

	// Delete test case
	if err := h.testCaseUsecase.DeleteTestCase(testCaseID, adminID); err != nil {
		errMsg := err.Error()
		switch errMsg {
		case "test case not found":
			RespondError(w, http.StatusNotFound, errMsg)
		default:
			h.logger.Error("Test case deletion failed", zap.Error(err))
			RespondError(w, http.StatusInternalServerError, "failed to delete test case")
		}
		return
	}

	h.logger.Info("Test case deleted successfully",
		zap.Int("test_case_id", testCaseID),
		zap.Int("admin_id", adminID),
	)

	RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Test case deleted successfully",
	})
}

// DeleteAllTestCases deletes all test cases for a problem (admin only)
func (h *TestCaseHandler) DeleteAllTestCases(w http.ResponseWriter, r *http.Request) {
	// Get admin ID from context
	adminID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get problem ID from path
	problemID, err := strconv.Atoi(r.PathValue("problem_id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid problem ID")
		return
	}

	// Delete all test cases
	if err := h.testCaseUsecase.DeleteAllTestCases(problemID, adminID); err != nil {
		errMsg := err.Error()
		switch errMsg {
		case "problem not found":
			RespondError(w, http.StatusNotFound, errMsg)
		default:
			h.logger.Error("Delete all test cases failed", zap.Error(err))
			RespondError(w, http.StatusInternalServerError, "failed to delete test cases")
		}
		return
	}

	h.logger.Info("All test cases deleted successfully",
		zap.Int("problem_id", problemID),
		zap.Int("admin_id", adminID),
	)

	RespondJSON(w, http.StatusOK, map[string]string{
		"message": "All test cases deleted successfully",
	})
}

// ReorderTestCases reorders test cases (admin only)
func (h *TestCaseHandler) ReorderTestCases(w http.ResponseWriter, r *http.Request) {
	// Get admin ID from context
	adminID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse request
	var req domain.ReorderTestCasesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Invalid JSON in reorder test cases request", zap.Error(err))
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Reorder test cases
	if err := h.testCaseUsecase.ReorderTestCases(&req, adminID); err != nil {
		// Handle validation errors
		var validationErr *uerror.ValidationError
		if errors.As(err, &validationErr) {
			h.logger.Warn("Reorder test cases validation failed",
				zap.Any("errors", validationErr.Errors),
			)
			RespondValidationError(w, validationErr.Errors)
			return
		}

		errMsg := err.Error()
		switch errMsg {
		case "test case not found", "test case does not belong to specified problem":
			RespondError(w, http.StatusBadRequest, errMsg)
		default:
			h.logger.Error("Test case reorder failed", zap.Error(err))
			RespondError(w, http.StatusInternalServerError, "failed to reorder test cases")
		}
		return
	}

	h.logger.Info("Test cases reordered successfully",
		zap.Int("problem_id", req.ProblemID),
		zap.Int("count", len(req.TestCases)),
		zap.Int("admin_id", adminID),
	)

	RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Test cases reordered successfully",
	})
}

// ListTestCases lists test cases with filters (admin only)
func (h *TestCaseHandler) ListTestCases(w http.ResponseWriter, r *http.Request) {
	// Get problem ID from path
	problemID, err := strconv.Atoi(r.PathValue("problem_id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid problem ID")
		return
	}

	req := &domain.ListTestCasesRequest{
		ProblemID: problemID,
		IsSample:  parseBoolQuery(r.URL.Query().Get("is_sample")),
		Page:      getIntQuery(r, "page", 1),
		Limit:     getIntQuery(r, "limit", 50),
	}

	testCases, total, err := h.testCaseUsecase.ListTestCases(req)
	if err != nil {
		h.logger.Error("Failed to list test cases", zap.Error(err))
		RespondError(w, http.StatusInternalServerError, "failed to retrieve test cases")
		return
	}

	response := PaginatedResponse[[]*domain.TestCase]{
		Data:  testCases,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}

	RespondPaginatedJSON(w, http.StatusOK, response)
}

// CountTestCasesByProblem returns test case count for a problem (admin only)
func (h *TestCaseHandler) CountTestCasesByProblem(w http.ResponseWriter, r *http.Request) {
	// Get problem ID from path
	problemID, err := strconv.Atoi(r.PathValue("problem_id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid problem ID")
		return
	}

	count, err := h.testCaseUsecase.CountTestCasesByProblem(problemID)
	if err != nil {
		h.logger.Error("Failed to count test cases", zap.Error(err))
		RespondError(w, http.StatusInternalServerError, "failed to count test cases")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]int{
		"problem_id": problemID,
		"count":      count,
	})
}

// ========== USER ENDPOINTS ==========

// GetSampleTestCases gets sample test cases for public problems (user accessible)
func (h *TestCaseHandler) GetSampleTestCases(w http.ResponseWriter, r *http.Request) {
	// Get problem ID from path
	problemID, err := strconv.Atoi(r.PathValue("problem_id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid problem ID")
		return
	}

	testCases, err := h.testCaseUsecase.GetSampleTestCases(problemID)
	if err != nil {
		errMsg := err.Error()
		switch errMsg {
		case "problem not found":
			RespondError(w, http.StatusNotFound, errMsg)
		case "problem not accessible":
			RespondError(w, http.StatusForbidden, "problem not accessible")
		default:
			h.logger.Error("Failed to get sample test cases", zap.Error(err))
			RespondError(w, http.StatusInternalServerError, "failed to retrieve sample test cases")
		}
		return
	}

	RespondJSON(w, http.StatusOK, testCases)
}

// ========== HELPER FUNCTIONS ==========

func parseBoolQuery(value string) *bool {
	if value == "" {
		return nil
	}
	if b, err := strconv.ParseBool(value); err == nil {
		return &b
	}
	return nil
}
