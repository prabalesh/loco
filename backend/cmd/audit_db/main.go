package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/pkg/config"
	"github.com/prabalesh/loco/backend/pkg/database"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()
	config.InitConfig()
	cfg := config.GetConfig()
	logger := zap.NewExample()
	db, err := database.NewPostgresDB(cfg.Database, logger)
	if err != nil {
		log.Fatal(err)
	}

	// 1. List All Problems
	var problems []domain.Problem
	db.DB.Find(&problems)
	log.Printf("--- Problems (%d) ---", len(problems))
	for _, p := range problems {
		log.Printf("ID: %d, Slug: %s, Title: %s", p.ID, p.Slug, p.Title)

		// 2. List Languages for this problem
		var pls []domain.ProblemLanguage
		db.DB.Where("problem_id = ?", p.ID).Find(&pls)
		log.Printf("  Languages: %d", len(pls))
		for _, pl := range pls {
			var lang domain.Language
			db.DB.First(&lang, pl.LanguageID)
			log.Printf("    - %s (ID: %d)", lang.Name, lang.ID)
		}
	}

	// 3. List All Languages in table
	var languages []domain.Language
	db.DB.Find(&languages)
	log.Printf("--- All Languages (%d) ---", len(languages))
	for _, l := range languages {
		log.Printf("ID: %d, Name: %s, Slug: %s", l.ID, l.Name, l.Slug)
	}

	log.Println("Database audit completed")
}
