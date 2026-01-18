package domain

import "time"

// Category represents a broad problem category (e.g., Algorithms, SQL, System Design)
type Category struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:255;not null;uniqueIndex"`
	Slug        string    `json:"slug" gorm:"size:255;not null;uniqueIndex"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
