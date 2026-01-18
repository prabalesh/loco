package main

import (
	"fmt"
	"log"
	"os"

	"github.com/prabalesh/loco/backend/internal/repository/postgres"
	"github.com/prabalesh/loco/backend/pkg/config"
	"github.com/prabalesh/loco/backend/pkg/database"
	"go.uber.org/zap"
)

func main() {
	// Setup env for config
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "loco_admin")
	os.Setenv("DB_PASSWORD", "toortoor")
	os.Setenv("DB_NAME", "loco")
	os.Setenv("JWT_SECRET", "test")
	os.Setenv("ACCESS_TOKEN_SECRET", "test")
	os.Setenv("REFRESH_TOKEN_SECRET", "test")
	os.Setenv("RESEND_API_KEY", "test")

	config.InitConfig()
	cfg := config.GetConfig()

	logger, _ := zap.NewDevelopment()
	db, err := database.NewPostgresDB(cfg.Database, logger)
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	repo := postgres.NewSubmissionRepository(db)

	fmt.Println("Testing GetTrendingProblems (7 days, limit 5)...")
	trending, err := repo.GetTrendingProblems(5, 7)
	if err != nil {
		log.Fatalf("GetTrendingProblems failed: %v", err)
	}
	fmt.Printf("Found %d trending problems:\n", len(trending))
	for _, p := range trending {
		fmt.Printf("- %s (ID: %d, Slug: %s): %d submissions\n", p.Title, p.ID, p.Slug, p.SubmissionCount)
	}

	fmt.Println("\nTesting GetLanguageStats...")
	langStats, err := repo.GetLanguageStats()
	if err != nil {
		log.Fatalf("GetLanguageStats failed: %v", err)
	}
	fmt.Printf("Found %d language stats:\n", len(langStats))
	for _, s := range langStats {
		fmt.Printf("- %s: %d submissions\n", s.LanguageName, s.Count)
	}
}
