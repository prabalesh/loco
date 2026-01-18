package seeder

import (
	"fmt"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/pkg/database"
	"go.uber.org/zap"
)

// SeedAchievements seeds the database with predefined achievements
func SeedAchievements(db *database.Database, log *zap.Logger) error {
	log.Info("Seeding achievements...")

	achievements := []domain.Achievement{
		// Getting Started (5)
		{Slug: "hello-world", Name: "Hello World", Description: "Make your first submission", XPReward: 10, Category: "getting-started", ConditionType: "count", ConditionValue: "1"},
		{Slug: "first-blood", Name: "First Blood", Description: "Get your first Accepted solution", XPReward: 50, Category: "getting-started", ConditionType: "status", ConditionValue: "Accepted"},
		{Slug: "bug-hunter", Name: "Bug Hunter", Description: "Get your first Wrong Answer", XPReward: 10, Category: "getting-started", ConditionType: "status", ConditionValue: "Wrong Answer"},
		{Slug: "speed-demon", Name: "Speed Demon", Description: "Get your first Time Limit Exceeded", XPReward: 10, Category: "getting-started", ConditionType: "status", ConditionValue: "Time Limit Exceeded"},
		{Slug: "memory-leak", Name: "Memory Leak", Description: "Get your first Memory Limit Exceeded", XPReward: 10, Category: "getting-started", ConditionType: "status", ConditionValue: "Memory Limit Exceeded"},

		// Problem Solving (10)
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

		// Difficulty Mastery - Easy (5)
		{Slug: "easy-peasy-i", Name: "Easy Peasy I", Description: "Solve 10 Easy problems", XPReward: 50, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Easy","count":10}`},
		{Slug: "easy-peasy-ii", Name: "Easy Peasy II", Description: "Solve 50 Easy problems", XPReward: 150, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Easy","count":50}`},
		{Slug: "easy-peasy-iii", Name: "Easy Peasy III", Description: "Solve 100 Easy problems", XPReward: 300, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Easy","count":100}`},
		{Slug: "easy-peasy-iv", Name: "Easy Peasy IV", Description: "Solve 200 Easy problems", XPReward: 600, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Easy","count":200}`},
		{Slug: "easy-peasy-v", Name: "Easy Peasy V", Description: "Solve 500 Easy problems", XPReward: 1500, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Easy","count":500}`},

		// Difficulty Mastery - Medium (5)
		{Slug: "medium-well-i", Name: "Medium Well I", Description: "Solve 10 Medium problems", XPReward: 100, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Medium","count":10}`},
		{Slug: "medium-well-ii", Name: "Medium Well II", Description: "Solve 50 Medium problems", XPReward: 300, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Medium","count":50}`},
		{Slug: "medium-well-iii", Name: "Medium Well III", Description: "Solve 100 Medium problems", XPReward: 600, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Medium","count":100}`},
		{Slug: "medium-well-iv", Name: "Medium Well IV", Description: "Solve 200 Medium problems", XPReward: 1200, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Medium","count":200}`},
		{Slug: "medium-well-v", Name: "Medium Well V", Description: "Solve 500 Medium problems", XPReward: 3000, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Medium","count":500}`},

		// Difficulty Mastery - Hard (5)
		{Slug: "hard-core-i", Name: "Hard Core I", Description: "Solve 10 Hard problems", XPReward: 200, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Hard","count":10}`},
		{Slug: "hard-core-ii", Name: "Hard Core II", Description: "Solve 50 Hard problems", XPReward: 600, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Hard","count":50}`},
		{Slug: "hard-core-iii", Name: "Hard Core III", Description: "Solve 100 Hard problems", XPReward: 1200, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Hard","count":100}`},
		{Slug: "hard-core-iv", Name: "Hard Core IV", Description: "Solve 200 Hard problems", XPReward: 2400, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Hard","count":200}`},
		{Slug: "hard-core-v", Name: "Hard Core V", Description: "Solve 500 Hard problems", XPReward: 6000, Category: "difficulty", ConditionType: "difficulty_count", ConditionValue: `{"difficulty":"Hard","count":500}`},

		// Streaks (6)
		{Slug: "getting-serious", Name: "Getting Serious", Description: "Maintain a 3 day streak", XPReward: 50, Category: "streak", ConditionType: "streak", ConditionValue: "3"},
		{Slug: "weekly-warrior", Name: "Weekly Warrior", Description: "Maintain a 7 day streak", XPReward: 100, Category: "streak", ConditionType: "streak", ConditionValue: "7"},
		{Slug: "fortnight-fighter", Name: "Fortnight Fighter", Description: "Maintain a 14 day streak", XPReward: 250, Category: "streak", ConditionType: "streak", ConditionValue: "14"},
		{Slug: "monthly-master", Name: "Monthly Master", Description: "Maintain a 30 day streak", XPReward: 500, Category: "streak", ConditionType: "streak", ConditionValue: "30"},
		{Slug: "century-club", Name: "Century Club", Description: "Maintain a 100 day streak", XPReward: 2000, Category: "streak", ConditionType: "streak", ConditionValue: "100"},
		{Slug: "year-of-code", Name: "Year of Code", Description: "Maintain a 365 day streak", XPReward: 5000, Category: "streak", ConditionType: "streak", ConditionValue: "365"},

		// Efficiency & Misc (2)
		{Slug: "one-shot", Name: "One Shot", Description: "Solve a problem on first attempt", XPReward: 20, Category: "misc", ConditionType: "specific", ConditionValue: "one-shot"},
		{Slug: "persistence", Name: "Persistence", Description: "Solve a problem after 10+ failed attempts", XPReward: 20, Category: "misc", ConditionType: "specific", ConditionValue: "persistence"},
	}

	created := 0
	updated := 0

	for _, ach := range achievements {
		var existing domain.Achievement
		err := db.DB.Where("slug = ?", ach.Slug).First(&existing).Error
		if err == nil {
			// Update if exists
			existing.Name = ach.Name
			existing.Description = ach.Description
			existing.XPReward = ach.XPReward
			existing.Category = ach.Category
			existing.ConditionType = ach.ConditionType
			existing.ConditionValue = ach.ConditionValue
			if err := db.DB.Save(&existing).Error; err != nil {
				log.Error("Failed to update achievement", zap.String("slug", ach.Slug), zap.Error(err))
			} else {
				updated++
			}
		} else {
			// Create new
			if err := db.DB.Create(&ach).Error; err != nil {
				log.Error("Failed to create achievement", zap.String("slug", ach.Slug), zap.Error(err))
			} else {
				created++
			}
		}
	}

	log.Info("Achievement seeding complete",
		zap.Int("created", created),
		zap.Int("updated", updated),
		zap.Int("total", len(achievements)),
	)

	return nil
}

// SeedTags seeds the database with predefined tags
func SeedTags(db *database.Database, log *zap.Logger) error {
	log.Info("Seeding tags...")

	tags := []domain.Tag{
		{Name: "Array", Slug: "array"},
		{Name: "String", Slug: "string"},
		{Name: "Hash Table", Slug: "hash-table"},
		{Name: "Dynamic Programming", Slug: "dynamic-programming"},
		{Name: "Math", Slug: "math"},
		{Name: "Sorting", Slug: "sorting"},
		{Name: "Greedy", Slug: "greedy"},
		{Name: "Depth-First Search", Slug: "depth-first-search"},
		{Name: "Binary Search", Slug: "binary-search"},
		{Name: "Database", Slug: "database"},
		{Name: "Breadth-First Search", Slug: "breadth-first-search"},
		{Name: "Tree", Slug: "tree"},
		{Name: "Matrix", Slug: "matrix"},
		{Name: "Two Pointers", Slug: "two-pointers"},
		{Name: "Bit Manipulation", Slug: "bit-manipulation"},
		{Name: "Stack", Slug: "stack"},
		{Name: "Design", Slug: "design"},
		{Name: "Heap (Priority Queue)", Slug: "heap-priority-queue"},
		{Name: "Graph", Slug: "graph"},
		{Name: "Simulation", Slug: "simulation"},
		{Name: "Backtracking", Slug: "backtracking"},
		{Name: "Prefix Sum", Slug: "prefix-sum"},
		{Name: "Counting", Slug: "counting"},
		{Name: "Sliding Window", Slug: "sliding-window"},
		{Name: "Union Find", Slug: "union-find"},
		{Name: "Linked List", Slug: "linked-list"},
		{Name: "Ordered Set", Slug: "ordered-set"},
		{Name: "Monotonic Stack", Slug: "monotonic-stack"},
		{Name: "Enumeration", Slug: "enumeration"},
		{Name: "Recursion", Slug: "recursion"},
		{Name: "Divide and Conquer", Slug: "divide-and-conquer"},
		{Name: "Binary Tree", Slug: "binary-tree"},
		{Name: "Trie", Slug: "trie"},
		{Name: "Bitmask", Slug: "bitmask"},
		{Name: "Queue", Slug: "queue"},
		{Name: "Memoization", Slug: "memoization"},
		{Name: "Segment Tree", Slug: "segment-tree"},
		{Name: "Geometry", Slug: "geometry"},
		{Name: "Topological Sort", Slug: "topological-sort"},
		{Name: "Binary Indexed Tree", Slug: "binary-indexed-tree"},
	}

	created := 0
	for _, tag := range tags {
		var existing domain.Tag
		err := db.DB.Where("slug = ?", tag.Slug).First(&existing).Error
		if err != nil {
			// Tag doesn't exist, create it
			if err := db.DB.Create(&tag).Error; err != nil {
				log.Error("Failed to create tag", zap.String("slug", tag.Slug), zap.Error(err))
			} else {
				created++
			}
		}
	}

	log.Info("Tag seeding complete",
		zap.Int("created", created),
		zap.Int("total", len(tags)),
	)

	return nil
}

// SeedCategories seeds the database with predefined categories
func SeedCategories(db *database.Database, log *zap.Logger) error {
	log.Info("Seeding categories...")

	categories := []domain.Category{
		{Name: "Algorithms", Slug: "algorithms", Description: "Algorithmic problem solving"},
		{Name: "Data Structures", Slug: "data-structures", Description: "Data structure implementation and usage"},
		{Name: "Mathematics", Slug: "mathematics", Description: "Mathematical problems"},
		{Name: "Database", Slug: "database", Description: "SQL and database queries"},
		{Name: "System Design", Slug: "system-design", Description: "System design and architecture"},
		{Name: "Concurrency", Slug: "concurrency", Description: "Concurrent programming"},
		{Name: "String Manipulation", Slug: "string-manipulation", Description: "String processing problems"},
		{Name: "Graph Theory", Slug: "graph-theory", Description: "Graph algorithms and problems"},
	}

	created := 0
	for _, category := range categories {
		var existing domain.Category
		err := db.DB.Where("slug = ?", category.Slug).First(&existing).Error
		if err != nil {
			// Category doesn't exist, create it
			if err := db.DB.Create(&category).Error; err != nil {
				log.Error("Failed to create category", zap.String("slug", category.Slug), zap.Error(err))
			} else {
				created++
			}
		}
	}

	log.Info("Category seeding complete",
		zap.Int("created", created),
		zap.Int("total", len(categories)),
	)

	return nil
}

// SeedAll runs all seeders
func SeedAll(db *database.Database, log *zap.Logger) error {
	log.Info("Running all seeders...")

	if err := SeedTags(db, log); err != nil {
		return fmt.Errorf("failed to seed tags: %w", err)
	}

	if err := SeedCategories(db, log); err != nil {
		return fmt.Errorf("failed to seed categories: %w", err)
	}

	if err := SeedAchievements(db, log); err != nil {
		return fmt.Errorf("failed to seed achievements: %w", err)
	}

	log.Info("All seeders completed successfully")
	return nil
}
