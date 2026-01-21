package v2

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/prabalesh/loco/backend/internal/delivery/handler"
	"github.com/prabalesh/loco/backend/internal/delivery/middleware"
	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/services/validation"
)

type ValidationHandler struct {
	validationService *validation.ValidationService
	languageRepo      domain.LanguageRepository
}

func NewValidationHandler(validationService *validation.ValidationService, languageRepo domain.LanguageRepository) *ValidationHandler {
	return &ValidationHandler{
		validationService: validationService,
		languageRepo:      languageRepo,
	}
}

// POST /api/v2/admin/problems/:id/validate
func (h *ValidationHandler) ValidateReferenceSolution(w http.ResponseWriter, r *http.Request) {
	// Check admin
	role, ok := middleware.GetUserRole(r.Context())
	if !ok || role != "admin" {
		handler.RespondError(w, http.StatusForbidden, "forbidden: admin access required")
		return
	}

	// Get problem ID
	problemIDStr := r.PathValue("id")
	problemID, err := strconv.Atoi(problemIDStr)
	if err != nil {
		handler.RespondError(w, http.StatusBadRequest, "invalid problem ID")
		return
	}

	// Parse request
	var req struct {
		LanguageSlug string `json:"language_slug"`
		Code         string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate inputs
	if req.Code == "" {
		handler.RespondError(w, http.StatusBadRequest, "code is required")
		return
	}
	if req.LanguageSlug == "" {
		handler.RespondError(w, http.StatusBadRequest, "language_slug is required")
		return
	}

	// Get language
	language, err := h.languageRepo.GetBySlug(req.LanguageSlug)
	if err != nil {
		handler.RespondError(w, http.StatusNotFound, "language not found")
		return
	}

	// Validate and save reference solution
	validateReq := validation.ValidateRequest{
		ProblemID:    problemID,
		LanguageSlug: req.LanguageSlug,
		Code:         req.Code,
	}

	referenceSolution, validationResult, err := h.validationService.SaveReferenceSolution(validateReq, language.ID)
	if err != nil {
		handler.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"reference_solution_id": referenceSolution.ID,
		"is_validated":          referenceSolution.IsValidated,
		"validation_result":     validationResult,
	}

	handler.RespondJSON(w, http.StatusOK, response)
}

// GET /api/v2/admin/problems/:id/validation-status
func (h *ValidationHandler) GetValidationStatus(w http.ResponseWriter, r *http.Request) {
	// Check admin
	role, ok := middleware.GetUserRole(r.Context())
	if !ok || role != "admin" {
		handler.RespondError(w, http.StatusForbidden, "forbidden: admin access required")
		return
	}

	// Get problem ID
	problemIDStr := r.PathValue("id")
	problemID, err := strconv.Atoi(problemIDStr)
	if err != nil {
		handler.RespondError(w, http.StatusBadRequest, "invalid problem ID")
		return
	}

	status, err := h.validationService.GetValidationStatus(problemID)
	if err != nil {
		handler.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handler.RespondJSON(w, http.StatusOK, status)
}
