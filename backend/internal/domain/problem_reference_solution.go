package domain

import (
	"time"

	"gorm.io/datatypes"
)

type ProblemReferenceSolution struct {
	ID                int            `json:"id" gorm:"primaryKey"`
	ProblemID         int            `json:"problem_id" gorm:"not null;index"`
	LanguageID        int            `json:"language_id" gorm:"not null;index"`
	Code              string         `json:"code" gorm:"type:text;not null"`
	IsValidated       bool           `json:"is_validated" gorm:"default:false"`
	ValidationResults datatypes.JSON `json:"validation_results" gorm:"type:jsonb"` // JSON with test results
	CreatedAt         time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time      `json:"updated_at" gorm:"autoUpdateTime"`

	Problem  Problem  `json:"problem,omitempty" gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE"`
	Language Language `json:"language,omitempty" gorm:"foreignKey:LanguageID;constraint:OnDelete:CASCADE"`
}
