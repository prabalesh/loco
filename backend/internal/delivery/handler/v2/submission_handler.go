package v2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/prabalesh/loco/backend/internal/delivery/handler"
	"github.com/prabalesh/loco/backend/internal/delivery/middleware"
	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/services/execution"
)

type SubmissionHandler struct {
	executionService *execution.ExecutionService
	problemRepo      domain.ProblemRepository
	testCaseRepo     domain.TestCaseRepository
	languageRepo     domain.LanguageRepository
}

func NewSubmissionHandler(
	executionService *execution.ExecutionService,
	problemRepo domain.ProblemRepository,
	testCaseRepo domain.TestCaseRepository,
	languageRepo domain.LanguageRepository,
) *SubmissionHandler {
	return &SubmissionHandler{
		executionService: executionService,
		problemRepo:      problemRepo,
		testCaseRepo:     testCaseRepo,
		languageRepo:     languageRepo,
	}
}

type SubmitCodeRequest struct {
	Code         string `json:"code"`
	LanguageSlug string `json:"language_slug"`
}

// POST /api/v2/problems/{problem_id}/submit
func (h *SubmissionHandler) SubmitCode(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		handler.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get problem ID
	problemIDStr := r.PathValue("problem_id")
	problemID, err := strconv.Atoi(problemIDStr)
	if err != nil {
		handler.RespondError(w, http.StatusBadRequest, "Invalid problem ID")
		return
	}

	// Parse request
	var req SubmitCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate
	if req.Code == "" {
		handler.RespondError(w, http.StatusBadRequest, "Code is required")
		return
	}
	if req.LanguageSlug == "" {
		handler.RespondError(w, http.StatusBadRequest, "Language is required")
		return
	}

	// Get language
	language, err := h.languageRepo.GetBySlug(req.LanguageSlug)
	if err != nil {
		handler.RespondError(w, http.StatusNotFound, "Language not found")
		return
	}

	// Get test cases
	testCases, err := h.testCaseRepo.GetByProblemID(problemID)
	if err != nil {
		handler.RespondError(w, http.StatusInternalServerError, "Failed to fetch test cases")
		return
	}
	if len(testCases) == 0 {
		handler.RespondError(w, http.StatusNotFound, "No test cases found for this problem")
		return
	}

	// Execute code
	execReq := execution.ExecutionRequest{
		ProblemID:  problemID,
		LanguageID: language.ID,
		UserCode:   req.Code,
		TestCases:  testCases,
	}

	result, err := h.executionService.ExecuteSubmission(execReq, req.LanguageSlug)
	if err != nil {
		handler.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Execution failed: %v", err))
		return
	}

	// Note: In a real scenario, you'd probably save the submission to the DB here.
	// For now, we're focusing on the execution and returning the results directly to the user
	// to enable the "Run/Test" feature.

	handler.RespondJSON(w, http.StatusOK, result)
	_ = userID // Keep userID for future use (e.g. saving to DB)
}
