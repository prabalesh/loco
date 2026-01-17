package domain

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"
)

// Language represents a programming language
type Language struct {
	ID              int            `json:"id" gorm:"primaryKey"`
	Slug            string         `json:"language_id" gorm:"column:slug;size:50;not null;uniqueIndex"`
	Name            string         `json:"name" gorm:"size:100;not null"`
	Version         string         `json:"version" gorm:"size:50"`
	Extension       string         `json:"extension" gorm:"size:10"`
	DefaultTemplate string         `json:"default_template" gorm:"type:text"`
	IsActive        bool           `json:"is_active" gorm:"default:true"`
	ExecutorConfig  ExecutorConfig `json:"executor_config" gorm:"type:jsonb;serializer:json"`
	CreatedAt       time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
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
	ProblemID            int        `json:"problem_id" gorm:"primaryKey"`
	LanguageID           int        `json:"language_id" gorm:"primaryKey"`
	LanguageName         string     `json:"language_name" gorm:"->"`
	LanguageVersion      string     `json:"language_version" gorm:"->"`
	FunctionCode         string     `json:"function_code" gorm:"type:text"`
	MainCode             string     `json:"main_code" gorm:"type:text"`
	SolutionCode         string     `json:"solution_code" gorm:"type:text"`
	IsValidated          bool       `json:"is_validated" gorm:"default:false"`
	ValidatedAt          *time.Time `json:"validated_at,omitempty"`
	LastValidationStatus string     `json:"last_validation_status" gorm:"size:50"`
	LastValidationError  string     `json:"last_validation_error" gorm:"type:text"`
	LastPassCount        int        `json:"last_pass_count" gorm:"default:0"`
	LastTotalCount       int        `json:"last_total_count" gorm:"default:0"`
	CreatedAt            time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (pl *ProblemLanguage) GetCombinedCode(template string, implementationCode string) string {
	if implementationCode == "" {
		implementationCode = pl.SolutionCode
	}

	completeCode := strings.Replace(template, "##funccodegoeshere", implementationCode, 1)
	completeCode = strings.Replace(completeCode, "##maincodegoeshere", pl.MainCode, 1)

	return completeCode
}

func (pl *ProblemLanguage) GetAdminCombinedCode(template string, implementationCode string) string {
	if implementationCode == "" {
		implementationCode = pl.SolutionCode
	}

	funcCode := strings.Replace(pl.FunctionCode, "##codegoeshere", pl.SolutionCode, 1)

	completeCode := strings.Replace(template, "##funccodegoeshere", funcCode, 1)
	completeCode = strings.Replace(completeCode, "##maincodegoeshere", pl.MainCode, 1)

	return completeCode
}
