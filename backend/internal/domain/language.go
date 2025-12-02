package domain

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Language represents a programming language
type Language struct {
	ID              int            `json:"id" db:"id"`
	LanguageID      string         `json:"language_id" db:"language_id"` // python, cpp, javascript
	Name            string         `json:"name" db:"name"`               // Python 3, C++
	Version         string         `json:"version" db:"version"`
	Extension       string         `json:"extension" db:"extension"` // .py, .cpp
	DefaultTemplate string         `json:"default_template" db:"default_template"`
	IsActive        bool           `json:"is_active" db:"is_active"`
	ExecutorConfig  ExecutorConfig `json:"executor_config" db:"executor_config"` // JSONB
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
}

// ExecutorConfig is a custom type for JSONB handling
type ExecutorConfig map[string]interface{}

// Value implements driver.Valuer for database writes
func (ec ExecutorConfig) Value() (driver.Value, error) {
	return json.Marshal(ec)
}

// Scan implements sql.Scanner for database reads
func (ec *ExecutorConfig) Scan(value interface{}) error {
	if value == nil {
		*ec = make(ExecutorConfig)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, ec)
}

// ProblemLanguage represents the many-to-many relationship between problems and languages
type ProblemLanguage struct {
	ProblemID      int        `json:"problem_id" db:"problem_id"`
	LanguageID     int        `json:"language_id" db:"language_id"`
	Language       *Language  `json:"language,omitempty"`                             // populated via JOIN
	CustomTemplate *string    `json:"custom_template,omitempty" db:"custom_template"` // nullable
	IsEnabled      bool       `json:"is_enabled" db:"is_enabled"`
	IsVerified     bool       `json:"is_verified" db:"is_verified"`
	VerifiedAt     *time.Time `json:"verified_at,omitempty" db:"verified_at"` // nullable
	OrderIndex     int        `json:"order_index" db:"order_index"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
}
