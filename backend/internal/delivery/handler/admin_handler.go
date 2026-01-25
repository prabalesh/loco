// internal/delivery/handler/admin_handler.go
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/prabalesh/loco/backend/internal/delivery/middleware"
	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/domain/dto"
	"github.com/prabalesh/loco/backend/internal/usecase"
	"go.uber.org/zap"
)

type AdminHandler struct {
	adminUsecase *usecase.AdminUsecase
	logger       *zap.Logger
}

func NewAdminHandler(adminUsecase *usecase.AdminUsecase, logger *zap.Logger) *AdminHandler {
	return &AdminHandler{
		adminUsecase: adminUsecase,
		logger:       logger,
	}
}

// ListUsers - Get all users
func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	adminID, _ := middleware.GetUserID(r.Context())

	users, err := h.adminUsecase.GetAllUsers()
	if err != nil {
		h.logger.Error("Failed to fetch users", zap.Error(err), zap.Int("admin_id", adminID))
		RespondError(w, http.StatusInternalServerError, "failed to fetch users")
		return
	}

	h.logger.Info("Admin fetched user list", zap.Int("admin_id", adminID), zap.Int("count", len(users)))
	RespondJSON(w, http.StatusOK, users)
}

// GetUser - Get single user by ID
func (h *AdminHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.adminUsecase.GetUserByID(userID)
	if err != nil {
		RespondError(w, http.StatusNotFound, "user not found")
		return
	}

	RespondJSON(w, http.StatusOK, dto.ToUserResponse(user))
}

// DeleteUser - Delete user
func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	adminID, _ := middleware.GetUserID(r.Context())
	userIDStr := r.PathValue("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := h.adminUsecase.DeleteUser(adminID, userID); err != nil {
		h.logger.Error("Failed to delete user", zap.Error(err), zap.Int("admin_id", adminID))
		switch err.Error() {
		case "cannot delete admin users":
			RespondError(w, http.StatusForbidden, err.Error())

		case "user not found":
			RespondError(w, http.StatusNotFound, err.Error())

		default:
			RespondError(w, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	h.logger.Info("Admin deleted user", zap.Int("admin_id", adminID), zap.Int("deleted_user_id", userID))
	RespondJSON(w, http.StatusOK, map[string]string{"message": "user deleted successfully"})
}

// UpdateUserRole - Change user role
func (h *AdminHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	adminID, _ := middleware.GetUserID(r.Context())
	userIDStr := r.PathValue("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req dto.UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.adminUsecase.UpdateUserRole(adminID, userID, req.Role); err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.Info("Admin updated user role",
		zap.Int("admin_id", adminID),
		zap.Int("user_id", userID),
		zap.String("new_role", req.Role),
	)

	RespondJSON(w, http.StatusOK, map[string]string{"message": "role updated successfully"})
}

// UpdateUserStatus - Activate/Deactivate user
func (h *AdminHandler) UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	adminID, _ := middleware.GetUserID(r.Context())
	userIDStr := r.PathValue("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req dto.UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.adminUsecase.UpdateUserStatus(adminID, userID, req.IsActive); err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.Info("Admin updated user status",
		zap.Int("admin_id", adminID),
		zap.Int("user_id", userID),
		zap.Bool("is_active", req.IsActive),
	)

	RespondJSON(w, http.StatusOK, map[string]string{"message": "status updated successfully"})
}

// GetAnalytics - Get dashboard analytics
func (h *AdminHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	adminID, _ := middleware.GetUserID(r.Context())

	analytics, err := h.adminUsecase.GetAnalytics()
	if err != nil {
		h.logger.Error("Failed to get analytics", zap.Error(err), zap.Int("admin_id", adminID))
		RespondError(w, http.StatusInternalServerError, "failed to get analytics")
		return
	}

	RespondJSON(w, http.StatusOK, analytics)
}

// ListPistonExecutions - Get Piston execution logs
func (h *AdminHandler) ListPistonExecutions(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	executions, total, err := h.adminUsecase.ListPistonExecutions(page, limit)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondPaginatedJSON(w, http.StatusOK, PaginatedResponse[[]domain.PistonExecution]{
		Total: int(total),
		Page:  page,
		Limit: limit,
		Data:  executions,
	})
}

// ListSubmissions - Get all global submissions
func (h *AdminHandler) ListSubmissions(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	submissions, total, err := h.adminUsecase.ListSubmissions(page, limit)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondPaginatedJSON(w, http.StatusOK, PaginatedResponse[[]domain.Submission]{
		Total: int(total),
		Page:  page,
		Limit: limit,
		Data:  submissions,
	})
}
