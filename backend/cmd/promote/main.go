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

	email := "admin_v2@example.com"
	log.Printf("Promoting user %s...", email)

	var user domain.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		log.Fatal("User not found:", err)
	}

	user.Role = "admin"
	user.EmailVerified = true

	if err := db.DB.Save(&user).Error; err != nil {
		log.Fatal("Failed to update user:", err)
	}

	log.Println("User promoted to admin and verified successfully!")
}
