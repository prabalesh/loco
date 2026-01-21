package domain

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

// TestCase represents a test case (sample or hidden) for a problem
type TestCase struct {
	ID               int              `json:"id" gorm:"primaryKey"`
	ProblemID        int              `json:"problem_id" gorm:"not null;index"`
	Input            string           `json:"input" gorm:"type:text"`
	ExpectedOutput   string           `json:"expected_output" gorm:"type:text"`
	IsSample         bool             `json:"is_sample" gorm:"default:false"`
	ValidationConfig ValidationConfig `json:"validation_config" gorm:"type:jsonb;serializer:json"`
	OrderIndex       int              `json:"order_index" gorm:"default:0"`

	// V2 additions
	ExpectedOutputs *datatypes.JSON `json:"expected_outputs,omitempty" gorm:"type:jsonb"`
	InputSize       *int            `json:"input_size,omitempty"`
	TimeLimitMs     *int            `json:"time_limit_ms,omitempty"`
	MemoryLimitMb   *int            `json:"memory_limit_mb,omitempty"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
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
