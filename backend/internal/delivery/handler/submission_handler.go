package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/prabalesh/loco/backend/internal/delivery/middleware"
	"github.com/prabalesh/loco/backend/internal/domain/dto"
	"github.com/prabalesh/loco/backend/internal/usecase"
	"go.uber.org/zap"
)

type SubmissionHandler struct {
	submissionUsecase *usecase.SubmissionUsecase
	logger            *zap.Logger
}

func NewSubmissionHandler(submissionUsecase *usecase.SubmissionUsecase, logger *zap.Logger) *SubmissionHandler {
	return &SubmissionHandler{
		submissionUsecase: submissionUsecase,
		logger:            logger,
	}
}

func (h *SubmissionHandler) Submit(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	problemIDStr := r.PathValue("problem_id")
	problemID, err := strconv.Atoi(problemIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid problem id")
		return
	}

	var req dto.CreateSubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	submission, err := h.submissionUsecase.Submit(userID, problemID, &req)
	if err != nil {
		h.logger.Error("Submission failed", zap.Error(err))
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, dto.ToSubmissionResponse(submission))
}

func (h *SubmissionHandler) GetSubmission(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid submission id")
		return
	}

	submission, err := h.submissionUsecase.GetSubmission(id)
	if err != nil {
		h.logger.Error("Failed to get submission", zap.Int("id", id), zap.Error(err))
		RespondError(w, http.StatusNotFound, "submission not found")
		return
	}

	// Sanitize Results (Hide hidden test cases)
	isAdmin := false
	if role, ok := middleware.GetUserRole(r.Context()); ok && role == "admin" {
		isAdmin = true
	}

	if !isAdmin {
		submission.Sanitize()
	}

	RespondJSON(w, http.StatusOK, dto.ToSubmissionResponse(submission))
}

func (h *SubmissionHandler) ListUserProblemSubmissions(w http.ResponseWriter, r *http.Request) {
	problemIDStr := r.PathValue("problem_id")
	problemID, err := strconv.Atoi(problemIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid problem id")
		return
	}
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	offset := (page - 1) * limit
	submissions, count, err := h.submissionUsecase.GetUserProblemSubmissions(userID, problemID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to list submissions", zap.Error(err))
		RespondError(w, http.StatusInternalServerError, "failed to list submissions")
		return
	}

	// Sanitize Results (Hide hidden test cases)
	isAdmin := false
	if role, ok := middleware.GetUserRole(r.Context()); ok && role == "admin" {
		isAdmin = true
	}

	if !isAdmin {
		for i := range submissions {
			submissions[i].Sanitize()
		}
	}

	submissionResponses := make([]dto.SubmissionResponse, len(submissions))
	for i := range submissions {
		submissionResponses[i] = dto.ToSubmissionResponse(&submissions[i])
	}

	RespondPaginatedJSON(w, http.StatusOK, PaginatedResponse[[]dto.SubmissionResponse]{
		Data:  submissionResponses,
		Total: int(count),
		Page:  page,
		Limit: limit,
	})
}

func (h *SubmissionHandler) ListUserSubmissions(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	offset := (page - 1) * limit
	submissions, count, err := h.submissionUsecase.GetUserSubmissions(userID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to list user submissions", zap.Error(err))
		RespondError(w, http.StatusInternalServerError, "failed to list submissions")
		return
	}

	// Sanitize Results (Hide hidden test cases)
	isAdmin := false
	if role, ok := middleware.GetUserRole(r.Context()); ok && role == "admin" {
		isAdmin = true
	}

	if !isAdmin {
		for i := range submissions {
			submissions[i].Sanitize()
		}
	}

	submissionResponses := make([]dto.SubmissionResponse, len(submissions))
	for i := range submissions {
		submissionResponses[i] = dto.ToSubmissionResponse(&submissions[i])
	}

	RespondPaginatedJSON(w, http.StatusOK, PaginatedResponse[[]dto.SubmissionResponse]{
		Data:  submissionResponses,
		Total: int(count),
		Page:  page,
		Limit: limit,
	})
}

func (h *SubmissionHandler) ListAdminUserSubmissions(w http.ResponseWriter, r *http.Request) {
	problemIDStr := r.PathValue("problem_id")
	problemID, err := strconv.Atoi(problemIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid problem id")
		return
	}
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	offset := (page - 1) * limit
	submissions, count, err := h.submissionUsecase.GetUserProblemSubmissions(userID, problemID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to list submissions", zap.Error(err))
		RespondError(w, http.StatusInternalServerError, "failed to list submissions")
		return
	}

	// Sanitize Results (Hide hidden test cases)
	isAdmin := false
	if role, ok := middleware.GetUserRole(r.Context()); ok && role == "admin" {
		isAdmin = true
	}

	if !isAdmin {
		for i := range submissions {
			submissions[i].Sanitize()
		}
	}

	submissionResponses := make([]dto.SubmissionResponse, len(submissions))
	for i := range submissions {
		submissionResponses[i] = dto.ToSubmissionResponse(&submissions[i])
	}

	RespondPaginatedJSON(w, http.StatusOK, PaginatedResponse[[]dto.SubmissionResponse]{
		Data:  submissionResponses,
		Total: int(count),
		Page:  page,
		Limit: limit,
	})
}

// AdminSubmit handles admin test submissions for problems
func (h *SubmissionHandler) AdminSubmit(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	problemIDStr := r.PathValue("id")
	problemID, err := strconv.Atoi(problemIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid problem id")
		return
	}

	var req struct {
		LanguageID int    `json:"language_id"`
		Code       string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	submissionReq := &dto.CreateSubmissionRequest{
		ProblemID:  problemID,
		LanguageID: req.LanguageID,
		Code:       req.Code,
	}

	submission, err := h.submissionUsecase.AdminSubmit(adminID, submissionReq)
	if err != nil {
		h.logger.Error("Admin submission failed", zap.Error(err), zap.Int("admin_id", adminID))
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.logger.Info("Admin submission created", zap.Int("submission_id", submission.ID), zap.Int("admin_id", adminID))
	RespondJSON(w, http.StatusCreated, dto.ToSubmissionResponse(submission))
}

// ListProblemSubmissions lists all submissions for a specific problem (admin only)
func (h *SubmissionHandler) ListProblemSubmissions(w http.ResponseWriter, r *http.Request) {
	problemIDStr := r.PathValue("id")
	problemID, err := strconv.Atoi(problemIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid problem id")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 50 {
		limit = 20
	}

	offset := (page - 1) * limit
	submissions, err := h.submissionUsecase.GetProblemSubmissions(problemID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to list problem submissions", zap.Error(err))
		RespondError(w, http.StatusInternalServerError, "failed to list submissions")
		return
	}

	submissionResponses := make([]dto.SubmissionResponse, len(submissions))
	for i := range submissions {
		submissionResponses[i] = dto.ToSubmissionResponse(&submissions[i])
	}

	RespondJSON(w, http.StatusOK, map[string]interface{}{
		"data":  submissionResponses,
		"page":  page,
		"limit": limit,
	})
}

// RunCode executes code against public test cases without creating a submission
func (h *SubmissionHandler) RunCode(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	problemIDStr := r.PathValue("problem_id")
	problemID, err := strconv.Atoi(problemIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid problem id")
		return
	}

	var req dto.RunCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	submission, err := h.submissionUsecase.RunCode(userID, problemID, &req)
	if err != nil {
		h.logger.Error("Run code failed", zap.Error(err), zap.Int("user_id", userID))
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusAccepted, dto.ToSubmissionResponse(submission))
}
