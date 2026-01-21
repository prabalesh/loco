package v2

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/prabalesh/loco/backend/internal/delivery/handler"
	"github.com/prabalesh/loco/backend/internal/delivery/middleware"
	"github.com/prabalesh/loco/backend/internal/services/problem"
)

type ProblemHandler struct {
	problemService *problem.ProblemService
}

func NewProblemHandler(problemService *problem.ProblemService) *ProblemHandler {
	return &ProblemHandler{
		problemService: problemService,
	}
}

// POST /api/v2/admin/problems
func (h *ProblemHandler) CreateProblem(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		handler.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	role, ok := middleware.GetUserRole(r.Context())
	if !ok || role != "admin" {
		handler.RespondError(w, http.StatusForbidden, "forbidden: admin access required")
		return
	}

	// Parse request
	var req problem.CreateProblemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Create problem
	createdProblem, err := h.problemService.CreateProblem(req, userID)
	if err != nil {
		handler.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Return response
	handler.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Problem created successfully",
		"problem": createdProblem,
	})
}

// GET /api/v2/problems/:slug
func (h *ProblemHandler) GetProblem(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	// Check if admin
	role, _ := middleware.GetUserRole(r.Context())
	isAdmin := role == "admin"

	// Get problem detail
	problemDetail, err := h.problemService.GetProblemDetail(slug, isAdmin)
	if err != nil {
		handler.RespondError(w, http.StatusNotFound, "problem not found")
		return
	}

	handler.RespondJSON(w, http.StatusOK, problemDetail)
}

// GET /api/v2/problems
func (h *ProblemHandler) ListProblems(w http.ResponseWriter, r *http.Request) {
	// Parse query params
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	filters := map[string]interface{}{}
	if difficulty := r.URL.Query().Get("difficulty"); difficulty != "" {
		filters["difficulty"] = difficulty
	}
	if status := r.URL.Query().Get("status"); status != "" {
		filters["status"] = status
	}

	// Get problems
	problems, total, err := h.problemService.GetAllProblems(filters, page, limit)
	if err != nil {
		handler.RespondError(w, http.StatusInternalServerError, "failed to fetch problems")
		return
	}

	totalPages := (int(total) + limit - 1) / limit

	handler.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"data":        problems,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	})
}

// GET /api/v2/admin/problems
func (h *ProblemHandler) AdminListProblems(w http.ResponseWriter, r *http.Request) {
	// Re-use ListProblems logic but ensure admin check
	role, ok := middleware.GetUserRole(r.Context())
	if !ok || role != "admin" {
		handler.RespondError(w, http.StatusForbidden, "forbidden: admin access required")
		return
	}
	h.ListProblems(w, r)
}

// GET /api/v2/admin/problems/:id
func (h *ProblemHandler) AdminGetProblem(w http.ResponseWriter, r *http.Request) {
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

	problem, err := h.problemService.GetAdminByID(problemID)
	if err != nil {
		handler.RespondError(w, http.StatusNotFound, "problem not found")
		return
	}

	handler.RespondJSON(w, http.StatusOK, problem)
}

// DELETE /api/v2/admin/problems/:id
func (h *ProblemHandler) DeleteProblem(w http.ResponseWriter, r *http.Request) {
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

	if err := h.problemService.DeleteProblem(problemID); err != nil {
		handler.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handler.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Problem deleted successfully",
	})
}

// POST /api/v2/admin/problems/:id/publish
func (h *ProblemHandler) PublishProblem(w http.ResponseWriter, r *http.Request) {
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

	if err := h.problemService.PublishProblem(problemID); err != nil {
		handler.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	handler.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "problem published successfully",
	})
}

// GET /api/v2/admin/custom-types
func (h *ProblemHandler) GetCustomTypes(w http.ResponseWriter, r *http.Request) {
	customTypes, err := h.problemService.GetCustomTypes()
	if err != nil {
		handler.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handler.RespondJSON(w, http.StatusOK, customTypes)
}

// POST /api/v2/admin/problems/:id/boilerplates
func (h *ProblemHandler) RegenerateBoilerplates(w http.ResponseWriter, r *http.Request) {
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

	if err := h.problemService.RegenerateBoilerplates(problemID); err != nil {
		handler.RespondError(w, http.StatusInternalServerError, "failed to regenerate boilerplates: "+err.Error())
		return
	}

	handler.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "boilerplates regenerated successfully",
	})
}
