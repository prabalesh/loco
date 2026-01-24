package usecase

import (
	"errors"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/domain/dto"
	"github.com/prabalesh/loco/backend/internal/domain/uerror"
	"github.com/prabalesh/loco/backend/internal/domain/validator"
	"github.com/prabalesh/loco/backend/pkg/config"
	"github.com/prabalesh/loco/backend/pkg/utils"
	"go.uber.org/zap"
)

type LanguageUsecase struct {
	languageRepo domain.LanguageRepository
	cfg          *config.Config
	logger       *zap.Logger
}

func NewLanguageUsecase(languageRepo domain.LanguageRepository, cfg *config.Config, logger *zap.Logger) *LanguageUsecase {
	return &LanguageUsecase{
		languageRepo: languageRepo,
		cfg:          cfg,
		logger:       logger,
	}
}

// ========== ADMIN OPERATIONS ==========

// CreateLanguage creates a new programming language support
func (u *LanguageUsecase) CreateLanguage(req *dto.CreateLanguageRequest, adminID int) (*domain.Language, error) {
	// Validation
	if validationErrors := validator.ValidateCreateLanguageRequest(req); len(validationErrors) > 0 {
		u.logger.Warn("Create language validation failed",
			zap.Any("errors", validationErrors),
		)
		return nil, &uerror.ValidationError{Errors: validationErrors}
	}

	// Check if language_id already exists
	_, err := u.languageRepo.GetBySlug(req.LanguageID)
	if err == nil {
		u.logger.Warn("Language creation failed: language_id already exists",
			zap.String("language_id", req.LanguageID),
		)
		return nil, errors.New("language already exists")
	}

	lang := &domain.Language{
		Slug:            req.LanguageID,
		Name:            req.Name,
		Version:         req.Version,
		Extension:       req.Extension,
		DefaultTemplate: req.DefaultTemplate,
		IsActive:        true, // Default active
		ExecutorConfig:  req.ExecutorConfig,
	}

	if err := u.languageRepo.Create(lang); err != nil {
		u.logger.Error("Failed to create language in database",
			zap.Error(err),
			zap.String("language_id", req.LanguageID),
			zap.Int("admin_id", adminID),
		)
		return nil, errors.New("failed to create language")
	}

	u.logger.Info("Language created successfully",
		zap.Int("language_id", lang.ID),
		zap.String("name", lang.Name),
		zap.Int("created_by", adminID),
	)

	return lang, nil
}

// UpdateLanguage updates an existing language
func (u *LanguageUsecase) UpdateLanguage(languageID int, req *dto.UpdateLanguageRequest, adminID int) (*domain.Language, error) {
	// Validation
	if validationErrors := validator.ValidateUpdateLanguageRequest(req); len(validationErrors) > 0 {
		u.logger.Warn("Update language validation failed",
			zap.Any("errors", validationErrors),
		)
		return nil, &uerror.ValidationError{Errors: validationErrors}
	}

	// Get existing language
	lang, err := u.languageRepo.GetByID(languageID)
	if err != nil {
		u.logger.Warn("Language not found for update",
			zap.Int("language_id", languageID),
		)
		return nil, errors.New("language not found")
	}

	// Update fields (only non-empty fields)
	if req.LanguageID != "" && req.LanguageID != lang.Slug {
		// Check if new language_id already exists
		existing, _ := u.languageRepo.GetBySlug(req.LanguageID)
		if existing != nil {
			return nil, errors.New("language_id already exists")
		}
		lang.Slug = req.LanguageID
	}

	if req.Name != "" {
		lang.Name = req.Name
	}

	if req.Version != "" {
		lang.Version = req.Version
	}

	if req.Extension != "" {
		lang.Extension = req.Extension
	}

	if req.DefaultTemplate != "" {
		lang.DefaultTemplate = req.DefaultTemplate
	}

	if req.ExecutorConfig != nil {
		lang.ExecutorConfig = req.ExecutorConfig
	}

	lang.IsActive = req.IsActive

	if err := u.languageRepo.Update(lang); err != nil {
		u.logger.Error("Failed to update language",
			zap.Error(err),
			zap.Int("language_id", languageID),
		)
		return nil, errors.New("failed to update language")
	}

	u.logger.Info("Language updated successfully",
		zap.Int("language_id", lang.ID),
		zap.String("name", lang.Name),
		zap.Int("updated_by", adminID),
	)

	return lang, nil
}

// DeleteLanguage deletes a language
func (u *LanguageUsecase) DeleteLanguage(languageID int, adminID int) error {
	// Check if language exists
	_, err := u.languageRepo.GetByID(languageID)
	if err != nil {
		u.logger.Warn("Language not found for deletion",
			zap.Int("language_id", languageID),
		)
		return errors.New("language not found")
	}

	if err := u.languageRepo.Delete(languageID); err != nil {
		u.logger.Error("Failed to delete language",
			zap.Error(err),
			zap.Int("language_id", languageID),
		)
		return errors.New("failed to delete language")
	}

	u.logger.Info("Language deleted successfully",
		zap.Int("language_id", languageID),
		zap.Int("deleted_by", adminID),
	)

	return nil
}

// ActivateLanguage activates a language
func (u *LanguageUsecase) ActivateLanguage(languageID int, adminID int) error {
	lang, err := u.languageRepo.GetByID(languageID)
	if err != nil {
		return errors.New("language not found")
	}

	if lang.IsActive {
		return errors.New("language is already active")
	}

	lang.IsActive = true
	if err := u.languageRepo.Update(lang); err != nil {
		u.logger.Error("Failed to activate language",
			zap.Error(err),
			zap.Int("language_id", languageID),
		)
		return errors.New("failed to activate language")
	}

	u.logger.Info("Language activated successfully",
		zap.Int("language_id", languageID),
		zap.Int("activated_by", adminID),
	)

	return nil
}

// DeactivateLanguage deactivates a language
func (u *LanguageUsecase) DeactivateLanguage(languageID int, adminID int) error {
	lang, err := u.languageRepo.GetByID(languageID)
	if err != nil {
		return errors.New("language not found")
	}

	if !lang.IsActive {
		return errors.New("language is already inactive")
	}

	lang.IsActive = false
	if err := u.languageRepo.Update(lang); err != nil {
		u.logger.Error("Failed to deactivate language",
			zap.Error(err),
			zap.Int("language_id", languageID),
		)
		return errors.New("failed to deactivate language")
	}

	u.logger.Info("Language deactivated successfully",
		zap.Int("language_id", languageID),
		zap.Int("deactivated_by", adminID),
	)

	return nil
}

// ========== USER OPERATIONS ==========

// GetLanguage retrieves a language by ID or language_id
func (u *LanguageUsecase) GetLanguage(identifier string) (*domain.Language, error) {
	var lang *domain.Language
	var err error

	// Try to get by ID first, then by language_id
	if id, parseErr := utils.ParseInt(identifier); parseErr == nil {
		lang, err = u.languageRepo.GetByID(id)
	} else {
		lang, err = u.languageRepo.GetBySlug(identifier)
	}

	if err != nil {
		u.logger.Warn("Language not found",
			zap.String("identifier", identifier),
		)
		return nil, errors.New("language not found")
	}

	return lang, nil
}

// ListActiveLanguages retrieves all active languages (for code editor)
func (u *LanguageUsecase) ListActiveLanguages() ([]*domain.Language, error) {
	languages, err := u.languageRepo.ListActive()
	if err != nil {
		u.logger.Error("Failed to list active languages",
			zap.Error(err),
		)
		return nil, errors.New("failed to retrieve languages")
	}

	result := make([]*domain.Language, len(languages))
	for i := range languages {
		result[i] = &languages[i]
	}

	return result, nil
}

func (u *LanguageUsecase) ListLanguages() ([]*domain.Language, error) {
	languages, err := u.languageRepo.GetAll()
	if err != nil {
		u.logger.Error("Failed to list languages",
			zap.Error(err),
		)
		return nil, errors.New("failed to retrieve languages")
	}

	result := make([]*domain.Language, len(languages))
	for i := range languages {
		result[i] = &languages[i]
	}

	return result, nil
}
