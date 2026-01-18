package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/pkg/config"
	"github.com/prabalesh/loco/backend/pkg/database"
	"github.com/prabalesh/loco/backend/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()
	config.InitConfig()
	cfg := config.GetConfig()
	logger.InitLogger("info")
	log := logger.GetLogger()

	db, err := database.NewPostgresDB(cfg.Database, log)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	achievements := []domain.Achievement{
		// Getting Started
		{Slug: "hello-world", Name: "Hello World", Description: "Make your first submission", XPReward: 10, Category: "getting-started", ConditionType: "count", ConditionValue: "1"},
		{Slug: "first-blood", Name: "First Blood", Description: "Get your first Accepted solution", XPReward: 50, Category: "getting-started", ConditionType: "status", ConditionValue: "Accepted"},
		{Slug: "bug-hunter", Name: "Bug Hunter", Description: "Get your first Wrong Answer", XPReward: 10, Category: "getting-started", ConditionType: "status", ConditionValue: "Wrong Answer"},
		{Slug: "speed-demon", Name: "Speed Demon", Description: "Get your first Time Limit Exceeded", XPReward: 10, Category: "getting-started", ConditionType: "status", ConditionValue: "Time Limit Exceeded"},
		{Slug: "memory-leak", Name: "Memory Leak", Description: "Get your first Memory Limit Exceeded", XPReward: 10, Category: "getting-started", ConditionType: "status", ConditionValue: "Memory Limit Exceeded"},

		// Problem Solving (Count)
		{Slug: "solver-i", Name: "Solver I", Description: "Solve 1 problem", XPReward: 10, Category: "solving", ConditionType: "count", ConditionValue: "1"},
		{Slug: "solver-ii", Name: "Solver II", Description: "Solve 10 problems", XPReward: 50, Category: "solving", ConditionType: "count", ConditionValue: "10"},
		{Slug: "solver-iii", Name: "Solver III", Description: "Solve 25 problems", XPReward: 100, Category: "solving", ConditionType: "count", ConditionValue: "25"},
		{Slug: "solver-iv", Name: "Solver IV", Description: "Solve 50 problems", XPReward: 250, Category: "solving", ConditionType: "count", ConditionValue: "50"},
		{Slug: "solver-v", Name: "Solver V", Description: "Solve 100 problems", XPReward: 500, Category: "solving", ConditionType: "count", ConditionValue: "100"},
		{Slug: "solver-vi", Name: "Solver VI", Description: "Solve 200 problems", XPReward: 1000, Category: "solving", ConditionType: "count", ConditionValue: "200"},
		{Slug: "solver-vii", Name: "Solver VII", Description: "Solve 300 problems", XPReward: 1500, Category: "solving", ConditionType: "count", ConditionValue: "300"},
		{Slug: "solver-viii", Name: "Solver VIII", Description: "Solve 400 problems", XPReward: 2000, Category: "solving", ConditionType: "count", ConditionValue: "400"},
		{Slug: "solver-ix", Name: "Solver IX", Description: "Solve 500 problems", XPReward: 2500, Category: "solving", ConditionType: "count", ConditionValue: "500"},
		{Slug: "solver-x", Name: "Solver X", Description: "Solve 1000 problems", XPReward: 5000, Category: "solving", ConditionType: "count", ConditionValue: "1000"},

		// Difficulty Mastery - Easy
		{Slug: "easy-peasy-i", Name: "Easy Peasy I", Description: "Solve 10 Easy problems", XPReward: 50, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Easy","count":10}`},
		{Slug: "easy-peasy-ii", Name: "Easy Peasy II", Description: "Solve 50 Easy problems", XPReward: 150, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Easy","count":50}`},
		{Slug: "easy-peasy-iii", Name: "Easy Peasy III", Description: "Solve 100 Easy problems", XPReward: 300, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Easy","count":100}`},
		{Slug: "easy-peasy-iv", Name: "Easy Peasy IV", Description: "Solve 200 Easy problems", XPReward: 600, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Easy","count":200}`},
		{Slug: "easy-peasy-v", Name: "Easy Peasy V", Description: "Solve 500 Easy problems", XPReward: 1500, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Easy","count":500}`},

		// Difficulty Mastery - Medium
		{Slug: "medium-well-i", Name: "Medium Well I", Description: "Solve 10 Medium problems", XPReward: 100, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Medium","count":10}`},
		{Slug: "medium-well-ii", Name: "Medium Well II", Description: "Solve 50 Medium problems", XPReward: 300, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Medium","count":50}`},
		{Slug: "medium-well-iii", Name: "Medium Well III", Description: "Solve 100 Medium problems", XPReward: 600, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Medium","count":100}`},
		{Slug: "medium-well-iv", Name: "Medium Well IV", Description: "Solve 200 Medium problems", XPReward: 1200, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Medium","count":200}`},
		{Slug: "medium-well-v", Name: "Medium Well V", Description: "Solve 500 Medium problems", XPReward: 3000, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Medium","count":500}`},

		// Difficulty Mastery - Hard
		{Slug: "hard-core-i", Name: "Hard Core I", Description: "Solve 10 Hard problems", XPReward: 200, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Hard","count":10}`},
		{Slug: "hard-core-ii", Name: "Hard Core II", Description: "Solve 50 Hard problems", XPReward: 600, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Hard","count":50}`},
		{Slug: "hard-core-iii", Name: "Hard Core III", Description: "Solve 100 Hard problems", XPReward: 1200, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Hard","count":100}`},
		{Slug: "hard-core-iv", Name: "Hard Core IV", Description: "Solve 200 Hard problems", XPReward: 2400, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Hard","count":200}`},
		{Slug: "hard-core-v", Name: "Hard Core V", Description: "Solve 500 Hard problems", XPReward: 6000, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Hard","count":500}`},

		// Streaks
		{Slug: "getting-serious", Name: "Getting Serious", Description: "3 Day Streak", XPReward: 50, Category: "streak", ConditionType: "streak", ConditionValue: "3"},
		{Slug: "weekly-warrior", Name: "Weekly Warrior", Description: "7 Day Streak", XPReward: 100, Category: "streak", ConditionType: "streak", ConditionValue: "7"},
		{Slug: "fortnight-fighter", Name: "Fortnight Fighter", Description: "14 Day Streak", XPReward: 250, Category: "streak", ConditionType: "streak", ConditionValue: "14"},
		{Slug: "monthly-master", Name: "Monthly Master", Description: "30 Day Streak", XPReward: 500, Category: "streak", ConditionType: "streak", ConditionValue: "30"},
		{Slug: "century-club", Name: "Century Club", Description: "100 Day Streak", XPReward: 2000, Category: "streak", ConditionType: "streak", ConditionValue: "100"},
		{Slug: "year-of-code", Name: "Year of Code", Description: "365 Day Streak", XPReward: 5000, Category: "streak", ConditionType: "streak", ConditionValue: "365"},

		// Efficiency & Misc
		{Slug: "one-shot", Name: "One Shot", Description: "Solved on first attempt", XPReward: 20, Category: "misc", ConditionType: "specific", ConditionValue: "one-shot"},
		{Slug: "persistence", Name: "Persistence", Description: "Solved after 10+ failed attempts", XPReward: 20, Category: "misc", ConditionType: "specific", ConditionValue: "persistence"},
	}

	for _, ach := range achievements {
		var existing domain.Achievement
		err := db.DB.Where("slug = ?", ach.Slug).First(&existing).Error
		if err == nil {
			// Update if exists (in case design changes)
			existing.Name = ach.Name
			existing.Description = ach.Description
			existing.XPReward = ach.XPReward
			existing.Category = ach.Category
			existing.ConditionType = ach.ConditionType
			existing.ConditionValue = ach.ConditionValue
			db.DB.Save(&existing)
			fmt.Printf("Updated achievement: %s\n", ach.Slug)
		} else {
			db.DB.Create(&ach)
			fmt.Printf("Created achievement: %s\n", ach.Slug)
		}
	}

	fmt.Println("Seeding complete.")
}
