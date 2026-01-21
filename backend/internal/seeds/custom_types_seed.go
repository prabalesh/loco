package seeds

import (
	"log"

	"github.com/prabalesh/loco/backend/internal/domain"
	"gorm.io/gorm"
)

// SeedCustomTypes seeds the custom_types table
func SeedCustomTypes(db *gorm.DB) error {
	customTypes := []domain.CustomType{
		{
			Name:        "TreeNode",
			Description: "Binary tree node with val, left, and right",
		},
		{
			Name:        "ListNode",
			Description: "Singly linked list node with val and next",
		},
		{
			Name:        "GraphNode",
			Description: "Graph node represented as adjacency list",
		},
		{
			Name:        "Node",
			Description: "N-ary tree node with val and children array",
		},
	}

	for _, ct := range customTypes {
		// Check if exists
		var existing domain.CustomType
		if err := db.Where("name = ?", ct.Name).First(&existing).Error; err == nil {
			log.Printf("CustomType %s already exists", ct.Name)
			continue // Already exists
		}

		// Create
		if err := db.Create(&ct).Error; err != nil {
			return err
		}
		log.Printf("Created CustomType: %s", ct.Name)
	}

	return nil
}
