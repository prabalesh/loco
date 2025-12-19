package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/prabalesh/loco/backend/internal/delivery/middleware"
	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/domain/uerror"
	"github.com/prabalesh/loco/backend/internal/usecase"
	"github.com/prabalesh/loco/backend/pkg/config"
	"go.uber.org/zap"
)

type ProblemHandler struct {
	problemUsecase *usecase.ProblemUsecase
	logger         *zap.Logger
	cfg            *config.Config
}

func NewProblemHandler(problemUsecase *usecase.ProblemUsecase, logger *zap.Logger, cfg *config.Config) *ProblemHandler {
	return &ProblemHandler{
		problemUsecase: problemUsecase,
		logger:         logger,
		cfg:            cfg,
	}
}

// ========== USER ENDPOINTS ==========

// GetProblem retrieves a single problem (public endpoint)
func (h *ProblemHandler) GetProblem(w http.ResponseWriter, r *http.Request) {
	identifier := r.PathValue("id") // Can be ID or slug

	problem, err := h.problemUsecase.GetProblem(identifier)
	if err != nil {
		h.logger.Warn("Problem not found",
			zap.String("identifier", identifier),
		)
		RespondError(w, http.StatusNotFound, "problem not found")
		return
	}

	h.logger.Info("Problem retrieved successfully",
		zap.Int("problem_id", problem.ID),
	)

	RespondJSON(w, http.StatusOK, problem)
}

// ListProblems retrieves problems with filters (public endpoint)
func (h *ProblemHandler) ListProblems(w http.ResponseWriter, r *http.Request) {
	req := &domain.ListProblemsRequest{
		Page:       getIntQuery(r, "page", 1),
		Limit:      getIntQuery(r, "limit", 20),
		Difficulty: r.URL.Query().Get("difficulty"),
		Search:     r.URL.Query().Get("search"),
		Tags:       r.URL.Query()["tags"], // Multiple tags support
	}

	problems, total, err := h.problemUsecase.ListProblems(req)
	if err != nil {
		h.logger.Error("Failed to list problems", zap.Error(err))
		RespondError(w, http.StatusInternalServerError, "failed to retrieve problems")
		return
	}

	response := PaginatedResponse[[]*domain.Problem]{
		Data:  problems,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}

	fmt.Println(response)

	RespondPaginatedJSON(w, http.StatusOK, response)
}

// ========== ADMIN ENDPOINTS ==========

// CreateProblem creates a new problem (admin only)
func (h *ProblemHandler) CreateProblem(w http.ResponseWriter, r *http.Request) {
	// Get admin ID from context
	adminID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse request
	var req domain.CreateProblemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Invalid JSON in create problem request", zap.Error(err))
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Create problem
	problem, err := h.problemUsecase.CreateProblem(&req, adminID)
	if err != nil {
		// Handle validation errors
		var validationErr *uerror.ValidationError
		if errors.As(err, &validationErr) {
			h.logger.Warn("Create problem validation failed",
				zap.Any("errors", validationErr.Errors),
			)
			RespondValidationError(w, validationErr.Errors)
			return
		}

		// Handle business logic errors
		errMsg := err.Error()

		switch errMsg {
		case "problem with similar title already exists":
			h.logger.Warn("Problem creation failed: duplicate slug", zap.String("error", errMsg))
			RespondError(w, http.StatusConflict, errMsg)
		default:
			h.logger.Error("Problem creation failed with unexpected error", zap.Error(err))
			RespondError(w, http.StatusInternalServerError, "failed to create problem")
		}
		return
	}

	h.logger.Info("Problem created successfully",
		zap.Int("problem_id", problem.ID),
		zap.Int("admin_id", adminID),
	)

	RespondJSON(w, http.StatusCreated, problem)
}

// UpdateProblem updates an existing problem (admin only)
func (h *ProblemHandler) UpdateProblem(w http.ResponseWriter, r *http.Request) {
	// Get admin ID from context
	adminID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get problem ID from path
	problemID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid problem ID")
		return
	}

	// Parse request
	var req domain.UpdateProblemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Invalid JSON in update problem request", zap.Error(err))
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Update problem
	problem, err := h.problemUsecase.UpdateProblem(problemID, &req, adminID)
	if err != nil {
		// Handle validation errors
		var validationErr *uerror.ValidationError
		if errors.As(err, &validationErr) {
			h.logger.Warn("Update problem validation failed",
				zap.Any("errors", validationErr.Errors),
			)
			RespondValidationError(w, validationErr.Errors)
			return
		}

		errMsg := err.Error()

		switch errMsg {
		case "problem not found":
			RespondError(w, http.StatusNotFound, errMsg)
		default:
			h.logger.Error("Problem update failed", zap.Error(err))
			RespondError(w, http.StatusInternalServerError, "failed to update problem")
		}
		return
	}

	h.logger.Info("Problem updated successfully",
		zap.Int("problem_id", problem.ID),
		zap.Int("admin_id", adminID),
	)

	RespondJSON(w, http.StatusOK, problem)
}

func (h *ProblemHandler) ValidateTestCases(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	problemID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid problem ID")
		return
	}

	err = h.problemUsecase.ValidateTestCases(problemID, adminID)
	if err != nil {
		h.logger.Warn("Test case validation failed", zap.Error(err))
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Test cases validated, problem step updated",
	})
}

// DeleteProblem deletes a problem (admin only)
func (h *ProblemHandler) DeleteProblem(w http.ResponseWriter, r *http.Request) {
	// Get admin ID from context
	adminID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get problem ID from path
	problemID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid problem ID")
		return
	}

	// Delete problem
	if err := h.problemUsecase.DeleteProblem(problemID, adminID); err != nil {
		errMsg := err.Error()

		switch errMsg {
		case "problem not found":
			RespondError(w, http.StatusNotFound, errMsg)
		default:
			h.logger.Error("Problem deletion failed", zap.Error(err))
			RespondError(w, http.StatusInternalServerError, "failed to delete problem")
		}
		return
	}

	h.logger.Info("Problem deleted successfully",
		zap.Int("problem_id", problemID),
		zap.Int("admin_id", adminID),
	)

	RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Problem deleted successfully",
	})
}

// ListAllProblems retrieves all problems including drafts (admin only)
func (h *ProblemHandler) ListAllProblems(w http.ResponseWriter, r *http.Request) {
	req := &domain.AdminListProblemsRequest{
		Page:       getIntQuery(r, "page", 1),
		Limit:      getIntQuery(r, "limit", 20),
		Difficulty: r.URL.Query().Get("difficulty"),
		Status:     r.URL.Query().Get("status"),
		Visibility: r.URL.Query().Get("visibility"),
		Search:     r.URL.Query().Get("search"),
		Tags:       r.URL.Query()["tags"],
	}

	problems, total, err := h.problemUsecase.ListAllProblems(req)
	if err != nil {
		h.logger.Error("Failed to list all problems (admin)", zap.Error(err))
		RespondError(w, http.StatusInternalServerError, "failed to retrieve problems")
		return
	}

	response := PaginatedResponse[[]*domain.Problem]{
		Data:  problems,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}

	RespondPaginatedJSON(w, http.StatusOK, response)
}

// PublishProblem changes problem status to published (admin only)
func (h *ProblemHandler) PublishProblem(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	problemID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid problem ID")
		return
	}

	if err := h.problemUsecase.PublishProblem(problemID, adminID); err != nil {
		errMsg := err.Error()

		switch errMsg {
		case "problem not found":
			RespondError(w, http.StatusNotFound, errMsg)
		case "problem is already published":
			RespondError(w, http.StatusBadRequest, errMsg)
		default:
			h.logger.Error("Problem publish failed", zap.Error(err))
			RespondError(w, http.StatusInternalServerError, "failed to publish problem")
		}
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Problem published successfully",
	})
}

// ArchiveProblem changes problem status to archived (admin only)
func (h *ProblemHandler) ArchiveProblem(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	problemID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid problem ID")
		return
	}

	if err := h.problemUsecase.ArchiveProblem(problemID, adminID); err != nil {
		h.logger.Error("Problem archive failed", zap.Error(err))
		RespondError(w, http.StatusInternalServerError, "failed to archive problem")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Problem archived successfully",
	})
}

// GetProblemStats returns problem statistics (admin only)
func (h *ProblemHandler) GetProblemStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.problemUsecase.GetProblemStats()
	if err != nil {
		h.logger.Error("Failed to get problem stats", zap.Error(err))
		RespondError(w, http.StatusInternalServerError, "failed to retrieve statistics")
		return
	}

	RespondJSON(w, http.StatusOK, stats)
}

// ========== HELPER FUNCTIONS ==========

func getIntQuery(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}
