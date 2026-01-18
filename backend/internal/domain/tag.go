package domain

import "time"

// Tag represents a problem category/topic
type Tag struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:255;not null;uniqueIndex"` // Array, Hash Table
	Slug      string    `json:"slug" gorm:"size:255;not null;uniqueIndex"` // array, hash-table
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
