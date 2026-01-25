package domain

import (
	"time"

	"gorm.io/datatypes"
)

type PistonExecution struct {
	ID        int            `json:"id" gorm:"primaryKey"`
	ProblemID int            `json:"problem_id" gorm:"index"`
	Language  string         `json:"language"`
	Version   string         `json:"version"`
	Code      string         `json:"code" gorm:"type:text"`
	Stdin     string         `json:"stdin" gorm:"type:text"`
	Response  datatypes.JSON `json:"response" gorm:"type:jsonb"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`

	// Optional link to submission if available
	SubmissionID *int     `json:"submission_id,omitempty" gorm:"index"`
	Problem      *Problem `json:"problem,omitempty" gorm:"foreignKey:ProblemID"`
}
