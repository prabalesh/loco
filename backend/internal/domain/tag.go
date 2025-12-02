package domain

import "time"

// Tag represents a problem category/topic
type Tag struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"` // Array, Hash Table
	Slug      string    `json:"slug" db:"slug"` // array, hash-table
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
