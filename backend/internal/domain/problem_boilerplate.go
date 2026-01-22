package domain

import (
	"time"

	"gorm.io/datatypes"
)

type ProblemBoilerplate struct {
	ID                  int            `json:"id" gorm:"primaryKey"`
	ProblemID           int            `json:"problem_id" gorm:"not null;index"`
	LanguageID          int            `json:"language_id" gorm:"not null;index"`
	StubCode            string         `json:"stub_code" gorm:"type:text;not null"`              // User-facing starter code
	TestHarnessTemplate datatypes.JSON `json:"test_harness_template" gorm:"type:jsonb;not null"` // Wrapper with {USER_CODE} placeholder
	CreatedAt           time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time      `json:"updated_at" gorm:"autoUpdateTime"`

	Problem  Problem  `json:"problem,omitempty" gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE"`
	Language Language `json:"language,omitempty" gorm:"foreignKey:LanguageID;constraint:OnDelete:CASCADE"`
}
