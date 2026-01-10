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
	ProblemID    int        `json:"problem_id" db:"problem_id"`
	LanguageID   int        `json:"language_id" db:"language_id"`
	FunctionCode string     `json:"function_code" db:"function_code"`
	MainCode     string     `json:"main_code" db:"main_code"`
	SolutionCode string     `json:"solution_code" db:"solution_code"`
	IsValidated  bool       `json:"is_validated" db:"is_validated"`
	ValidatedAt  *time.Time `json:"validated_at,omitempty" db:"validated_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}
