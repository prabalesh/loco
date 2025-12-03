package domain

// CreateLanguageRequest for admin language creation
type CreateLanguageRequest struct {
	LanguageID      string                 `json:"language_id" validate:"required,lowercase"` // "python", "cpp"
	Name            string                 `json:"name" validate:"required,min=2,max=50"`
	Version         string                 `json:"version" validate:"required,min=1,max=20"`
	Extension       string                 `json:"extension" validate:"required,oneof=.py .cpp .js .java .c .go .rs"`
	DefaultTemplate string                 `json:"default_template" validate:"omitempty,max=5000"`
	ExecutorConfig  map[string]interface{} `json:"executor_config" validate:"omitempty"`
}

// UpdateLanguageRequest for partial language updates
type UpdateLanguageRequest struct {
	LanguageID      string                 `json:"language_id,omitempty" validate:"omitempty,lowercase"`
	Name            string                 `json:"name,omitempty" validate:"omitempty,min=2,max=50"`
	Version         string                 `json:"version,omitempty" validate:"omitempty,min=1,max=20"`
	Extension       string                 `json:"extension,omitempty" validate:"omitempty,oneof=.py .cpp .js .java .c .go .rs"`
	DefaultTemplate string                 `json:"default_template,omitempty" validate:"omitempty,max=5000"`
	IsActive        bool                   `json:"is_active"`
	ExecutorConfig  map[string]interface{} `json:"executor_config,omitempty"`
}
