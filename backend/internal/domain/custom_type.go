package domain

import "time"

type CustomType struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:100;not null;uniqueIndex"` // TreeNode, ListNode, GraphNode
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}
