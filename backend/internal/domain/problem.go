package domain

import "time"

// Problem entity
type Problem struct {
	ID               int       `json:"id" gorm:"primaryKey"`
	Title            string    `json:"title" gorm:"size:255;not null;index"`
	Slug             string    `json:"slug" gorm:"size:255;not null;uniqueIndex"`
	Description      string    `json:"description" gorm:"type:text"`
	Difficulty       string    `json:"difficulty" gorm:"size:50;default:'medium'"`
	TimeLimit        int       `json:"time_limit" gorm:"default:1000"`
	MemoryLimit      int       `json:"memory_limit" gorm:"default:256"`
	CurrentStep      int       `json:"current_step" gorm:"default:1"`
	ValidatorType    string    `json:"validator_type" gorm:"size:50;default:'exact'"`
	InputFormat      string    `json:"input_format" gorm:"type:text"`
	OutputFormat     string    `json:"output_format" gorm:"type:text"`
	Constraints      string    `json:"constraints" gorm:"type:text"`
	Status           string    `json:"status" gorm:"size:50;default:'draft'"`
	Visibility       string    `json:"visibility" gorm:"size:50;default:'private'"`
	IsActive         bool      `json:"is_active" gorm:"default:true"`
	AcceptanceRate   float64   `json:"acceptance_rate" gorm:"default:0.0"`
	TotalSubmissions int       `json:"total_submissions" gorm:"column:total_submissions;default:0"`
	TotalAccepted    int       `json:"total_accepted" gorm:"column:total_accepted;default:0"`
	UserStatus       string    `json:"user_status,omitempty" gorm:"-"` // solved, attempted, unsolved
	CreatedBy        *int      `json:"created_by" gorm:"index"`
	Creator          *User     `json:"creator,omitempty" gorm:"foreignKey:CreatedBy;references:ID;constraint:OnDelete:SET NULL"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
