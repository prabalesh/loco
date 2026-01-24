package validator

import (
	"strings"

	"github.com/prabalesh/loco/backend/internal/domain/dto"
)

// ValidateCreateLanguageRequest validates language creation request
func ValidateCreateLanguageRequest(req *dto.CreateLanguageRequest) map[string]string {
	errors := make(map[string]string)

	if req.LanguageID == "" {
		errors["language_id"] = "language_id is required"
	}

	if req.Name == "" {
		errors["name"] = "name is required"
	} else if len(req.Name) < 1 || len(req.Name) > 50 {
		errors["name"] = "name must be 1-50 characters"
	}

	if req.Version == "" {
		errors["version"] = "version is required"
	} else if len(req.Version) > 20 {
		errors["version"] = "version must be max 20 characters"
	}

	if req.Extension == "" {
		errors["extension"] = "extension is required"
	} else if !strings.HasPrefix(req.Extension, ".") {
		errors["extension"] = "extension must start with dot (e.g., .py, .cpp)"
	}

	if req.DefaultTemplate != "" && len(req.DefaultTemplate) > 5000 {
		errors["default_template"] = "default_template must be max 5000 characters"
	}

	return errors
}

// ValidateUpdateLanguageRequest validates language update request
func ValidateUpdateLanguageRequest(req *dto.UpdateLanguageRequest) map[string]string {
	errors := make(map[string]string)

	// At least one field should be provided for update
	updated := false
	if req.LanguageID != "" {
		updated = true
	}
	if req.Name != "" {
		updated = true
	}
	if req.Version != "" {
		updated = true
	}
	if req.Extension != "" {
		updated = true
	}
	if req.DefaultTemplate != "" {
		updated = true
	}
	if req.IsActive {
		updated = true
	}
	if req.ExecutorConfig != nil {
		updated = true
	}

	if !updated {
		errors["general"] = "at least one field must be provided for update"
	}

	if req.LanguageID == "" {
		errors["language_id"] = "language id is required"
	}

	if req.Name != "" && (len(req.Name) < 1 || len(req.Name) > 50) {
		errors["name"] = "name must be 2-50 characters"
	}

	if req.Version != "" && len(req.Version) > 20 {
		errors["version"] = "version must be max 20 characters"
	}

	if req.Extension != "" {
		if !strings.HasPrefix(req.Extension, ".") {
			errors["extension"] = "extension must start with dot (e.g., .py, .cpp)"
		}
	}

	if req.DefaultTemplate != "" && len(req.DefaultTemplate) > 5000 {
		errors["default_template"] = "default_template must be max 5000 characters"
	}

	return errors
}
