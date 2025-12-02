package domain

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// TestCase represents a test case (sample or hidden) for a problem
type TestCase struct {
	ID               int              `json:"id" db:"id"`
	ProblemID        int              `json:"problem_id" db:"problem_id"`
	Input            string           `json:"input" db:"input"`
	ExpectedOutput   string           `json:"expected_output" db:"expected_output"`
	IsSample         bool             `json:"is_sample" db:"is_sample"`
	ValidationConfig ValidationConfig `json:"validation_config" db:"validation_config"` // JSONB
	OrderIndex       int              `json:"order_index" db:"order_index"`
	CreatedAt        time.Time        `json:"created_at" db:"created_at"`
}

// ValidationConfig is a custom type for JSONB handling
type ValidationConfig map[string]interface{}

// Value implements driver.Valuer for database writes
func (vc ValidationConfig) Value() (driver.Value, error) {
	if vc == nil {
		return nil, nil
	}
	return json.Marshal(vc)
}

// Scan implements sql.Scanner for database reads
func (vc *ValidationConfig) Scan(value interface{}) error {
	if value == nil {
		*vc = make(ValidationConfig)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, vc)
}
