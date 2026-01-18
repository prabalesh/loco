package main

import (
	"fmt"
	"log"
	"os"

	"github.com/prabalesh/loco/backend/pkg/config"
	"github.com/prabalesh/loco/backend/pkg/database"
	"go.uber.org/zap"
)

func main() {
	// Set required env vars
	os.Setenv("DB_PASSWORD", "toortoor")
	os.Setenv("ACCESS_TOKEN_SECRET", "test")
	os.Setenv("REFRESH_TOKEN_SECRET", "test")
	os.Setenv("RESEND_API_KEY", "test")

	config.InitConfig()
	cfg := config.GetConfig()
	logger, _ := zap.NewDevelopment()

	db, err := database.NewPostgresDB(cfg.Database, logger)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	type UserInfo struct {
		ID       int
		Username string
		Role     string
	}

	var users []UserInfo
	db.DB.Raw("SELECT id, username, role FROM users").Scan(&users)

	f, _ := os.Create("/tmp/db_inspect.txt")
	defer f.Close()

	fmt.Fprintf(f, "--- USERS ---\n")
	for _, u := range users {
		fmt.Fprintf(f, "ID: %d, Username: %s, Role: %s\n", u.ID, u.Username, u.Role)

		var solvedCount int64
		db.DB.Table("user_problem_stats").Where("user_id = ? AND status = ?", u.ID, "solved").Count(&solvedCount)
		fmt.Fprintf(f, "  Solved (status=solved): %d\n", solvedCount)

		var solvedCountUpper int64
		db.DB.Table("user_problem_stats").Where("user_id = ? AND status = ?", u.ID, "Solved").Count(&solvedCountUpper)
		fmt.Fprintf(f, "  Solved (status=Solved): %d\n", solvedCountUpper)
	}
	fmt.Println("Done. Results in /tmp/db_inspect.txt")
}
