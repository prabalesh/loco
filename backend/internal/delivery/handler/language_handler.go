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

type LanguageHandler struct {
	languageUsecase *usecase.LanguageUsecase
	logger          *zap.Logger
	cfg             *config.Config
}

func NewLanguageHandler(languageUsecase *usecase.LanguageUsecase, logger *zap.Logger, cfg *config.Config) *LanguageHandler {
	return &LanguageHandler{
		languageUsecase: languageUsecase,
		logger:          logger,
		cfg:             cfg,
	}
}

// GetLanguage retrieves a single language (public endpoint)
func (h *LanguageHandler) GetLanguage(w http.ResponseWriter, r *http.Request) {
	identifier := r.PathValue("id") // Can be ID or language_id

	language, err := h.languageUsecase.GetLanguage(identifier)
	if err != nil {
		h.logger.Warn("Language not found",
			zap.String("identifier", identifier),
		)
		RespondError(w, http.StatusNotFound, "language not found")
		return
	}

	h.logger.Info("Language retrieved successfully",
		zap.Int("language_id", language.ID),
	)

	RespondJSON(w, http.StatusOK, language)
}

// ListActiveLanguages retrieves all active languages (public endpoint for code editor)
func (h *LanguageHandler) ListActiveLanguages(w http.ResponseWriter, r *http.Request) {
	languages, err := h.languageUsecase.ListActiveLanguages()
	if err != nil {
		h.logger.Error("Failed to list active languages", zap.Error(err))
		RespondError(w, http.StatusInternalServerError, "failed to retrieve languages")
		return
	}

	RespondJSON(w, http.StatusOK, languages)
}

func (h *LanguageHandler) ListLanguages(w http.ResponseWriter, r *http.Request) {
	languages, err := h.languageUsecase.ListLanguages()
	if err != nil {
		h.logger.Error("Failed to list languages", zap.Error(err))
		RespondError(w, http.StatusInternalServerError, "failed to retrieve languages")
		return
	}

	RespondJSON(w, http.StatusOK, languages)
}

// ========== ADMIN ENDPOINTS ==========

// CreateLanguage creates a new language (admin only)
func (h *LanguageHandler) CreateLanguage(w http.ResponseWriter, r *http.Request) {
	// Get admin ID from context
	adminID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse request
	var req domain.CreateLanguageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Invalid JSON in create language request", zap.Error(err))
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Create language
	language, err := h.languageUsecase.CreateLanguage(&req, adminID)
	if err != nil {
		// Handle validation errors
		var validationErr *uerror.ValidationError
		if errors.As(err, &validationErr) {
			h.logger.Warn("Create language validation failed",
				zap.Any("errors", validationErr.Errors),
			)
			RespondValidationError(w, validationErr.Errors)
			return
		}

		// Handle business logic errors
		errMsg := err.Error()

		switch errMsg {
		case "language already exists":
			h.logger.Warn("Language creation failed: duplicate language_id", zap.String("error", errMsg))
			RespondError(w, http.StatusConflict, errMsg)
		default:
			h.logger.Error("Language creation failed with unexpected error", zap.Error(err))
			RespondError(w, http.StatusInternalServerError, "failed to create language")
		}
		return
	}

	h.logger.Info("Language created successfully",
		zap.Int("language_id", language.ID),
		zap.String("language_id", language.LanguageID),
		zap.Int("admin_id", adminID),
	)

	RespondJSON(w, http.StatusCreated, language)
}

// UpdateLanguage updates an existing language (admin only)
func (h *LanguageHandler) UpdateLanguage(w http.ResponseWriter, r *http.Request) {
	// Get admin ID from context
	adminID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get language ID from path
	languageID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid language ID")
		return
	}

	// Parse request
	var req domain.UpdateLanguageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Invalid JSON in update language request", zap.Error(err))
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Update language
	language, err := h.languageUsecase.UpdateLanguage(languageID, &req, adminID)
	if err != nil {
		// Handle validation errors
		var validationErr *uerror.ValidationError
		if errors.As(err, &validationErr) {
			h.logger.Warn("Update language validation failed",
				zap.Any("errors", validationErr.Errors),
			)
			RespondValidationError(w, validationErr.Errors)
			return
		}

		errMsg := err.Error()

		switch errMsg {
		case "language not found":
			RespondError(w, http.StatusNotFound, errMsg)
		case "language_id already exists":
			RespondError(w, http.StatusConflict, errMsg)
		default:
			h.logger.Error("Language update failed", zap.Error(err))
			RespondError(w, http.StatusInternalServerError, "failed to update language")
		}
		return
	}

	h.logger.Info("Language updated successfully",
		zap.Int("language_id", language.ID),
		zap.String("name", language.Name),
		zap.Int("admin_id", adminID),
	)

	RespondJSON(w, http.StatusOK, language)
}

// DeleteLanguage deletes a language (admin only)
func (h *LanguageHandler) DeleteLanguage(w http.ResponseWriter, r *http.Request) {
	// Get admin ID from context
	adminID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get language ID from path
	languageID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid language ID")
		return
	}

	// Delete language
	if err := h.languageUsecase.DeleteLanguage(languageID, adminID); err != nil {
		errMsg := err.Error()

		switch errMsg {
		case "language not found":
			RespondError(w, http.StatusNotFound, errMsg)
		default:
			h.logger.Error("Language deletion failed", zap.Error(err))
			RespondError(w, http.StatusInternalServerError, "failed to delete language")
		}
		return
	}

	h.logger.Info("Language deleted successfully",
		zap.Int("language_id", languageID),
		zap.Int("admin_id", adminID),
	)

	RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Language deleted successfully",
	})
}

// ActivateLanguage activates a language (admin only)
func (h *LanguageHandler) ActivateLanguage(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	languageID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid language ID")
		return
	}

	if err := h.languageUsecase.ActivateLanguage(languageID, adminID); err != nil {
		errMsg := err.Error()

		switch errMsg {
		case "language not found":
			RespondError(w, http.StatusNotFound, errMsg)
		case "language is already active":
			RespondError(w, http.StatusBadRequest, errMsg)
		default:
			h.logger.Error("Language activation failed", zap.Error(err))
			RespondError(w, http.StatusInternalServerError, "failed to activate language")
		}
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Language activated successfully",
	})
}

// DeactivateLanguage deactivates a language (admin only)
func (h *LanguageHandler) DeactivateLanguage(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.GetUserID(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	languageID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid language ID")
		return
	}

	if err := h.languageUsecase.DeactivateLanguage(languageID, adminID); err != nil {
		errMsg := err.Error()

		switch errMsg {
		case "language not found":
			RespondError(w, http.StatusNotFound, errMsg)
		case "language is already inactive":
			RespondError(w, http.StatusBadRequest, errMsg)
		default:
			h.logger.Error("Language deactivation failed", zap.Error(err))
			RespondError(w, http.StatusInternalServerError, "failed to deactivate language")
		}
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Language deactivated successfully",
	})
}
