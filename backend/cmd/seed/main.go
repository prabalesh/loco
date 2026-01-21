package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/seeds"
	"github.com/prabalesh/loco/backend/pkg/config"
	"github.com/prabalesh/loco/backend/pkg/database"
	"go.uber.org/zap"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	config.InitConfig()
	cfg := config.GetConfig()

	logger := zap.NewExample()

	db, err := database.NewPostgresDB(cfg.Database, logger)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migrate for custom types
	log.Println("Running migrations...")
	if err := db.DB.AutoMigrate(&domain.CustomType{}, &domain.TypeImplementation{}, &domain.Language{}); err != nil {
		log.Fatal("Failed to migrate:", err)
	}

	log.Println("Seeding languages...")
	if err := seeds.SeedLanguages(db.DB); err != nil {
		log.Fatal("Failed to seed languages:", err)
	}

	log.Println("Seeding custom types...")
	if err := seeds.SeedCustomTypes(db.DB); err != nil {
		log.Fatal("Failed to seed custom types:", err)
	}

	log.Println("Seeding type implementations...")
	if err := seeds.SeedTypeImplementations(db.DB); err != nil {
		log.Fatal("Failed to seed type implementations:", err)
	}

	log.Println("Seeding completed successfully!")
}
