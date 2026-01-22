package handler

import (
	"encoding/json"
	"net/http"

	"github.com/prabalesh/loco/backend/internal/delivery/middleware"
	"github.com/prabalesh/loco/backend/internal/services/bulk"
)

type BulkHandler struct {
	bulkService *bulk.BulkImportService
}

func NewBulkHandler(bulkService *bulk.BulkImportService) *BulkHandler {
	return &BulkHandler{
		bulkService: bulkService,
	}
}

// POST /api/v2/admin/problems/bulk
func (h *BulkHandler) BulkImportProblems(w http.ResponseWriter, r *http.Request) {
	// Check admin
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	role, ok := middleware.GetUserRole(r.Context())
	if !ok || role != "admin" {
		RespondError(w, http.StatusForbidden, "forbidden: admin access required")
		return
	}

	// Parse request
	var req bulk.BulkImportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate batch size
	if len(req.Problems) == 0 {
		RespondError(w, http.StatusBadRequest, "no problems provided")
		return
	}
	if len(req.Problems) > 100 {
		RespondError(w, http.StatusBadRequest, "maximum 100 problems per batch")
		return
	}

	// Import problems
	result, err := h.bulkService.BulkImport(req, userID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return result
	statusCode := http.StatusOK
	if result.TotalFailed > 0 && result.TotalCreated == 0 {
		statusCode = http.StatusBadRequest // All failed
	} else if result.TotalFailed > 0 {
		statusCode = http.StatusPartialContent // Partial success
	}

	RespondJSON(w, statusCode, result)
}

// POST /api/v2/admin/problems/bulk-async
func (h *BulkHandler) BulkImportProblemsAsync(w http.ResponseWriter, r *http.Request) {
	// Check admin
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	role, ok := middleware.GetUserRole(r.Context())
	if !ok || role != "admin" {
		RespondError(w, http.StatusForbidden, "forbidden: admin access required")
		return
	}

	// Parse request
	var req bulk.BulkImportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate batch size
	if len(req.Problems) > 1000 {
		RespondError(w, http.StatusBadRequest, "maximum 1000 problems for async import")
		return
	}

	// Start async import
	jobID, err := h.bulkService.BulkImportAsync(req, userID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return job ID
	RespondJSON(w, http.StatusAccepted, map[string]interface{}{
		"message": "Import job started",
		"job_id":  jobID,
		"status":  "processing",
	})
}
