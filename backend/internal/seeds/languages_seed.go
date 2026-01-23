package seeds

import (
	"log"

	"github.com/prabalesh/loco/backend/internal/domain"
	"gorm.io/gorm"
)

// SeedLanguages seeds the languages table
func SeedLanguages(db *gorm.DB) error {
	languages := []domain.Language{
		{Name: "Python", Slug: "python"},
		{Name: "JavaScript", Slug: "javascript"},
		{Name: "Java", Slug: "java"},
		{Name: "C++", Slug: "c++"},
		{Name: "C", Slug: "c"},
		{Name: "Go", Slug: "go"},
	}

	for _, lang := range languages {
		// Check if exists
		var existing domain.Language
		if err := db.Where("slug = ?", lang.Slug).First(&existing).Error; err == nil {
			log.Printf("Language %s already exists", lang.Name)
			continue // Already exists
		}

		// Create
		if err := db.Create(&lang).Error; err != nil {
			return err
		}
		log.Printf("Created Language: %s", lang.Name)
	}

	return nil
}
