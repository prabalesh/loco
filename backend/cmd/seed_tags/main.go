package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/pkg/config"
	"github.com/prabalesh/loco/backend/pkg/database"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()

	// Initialize config
	config.InitConfig()
	cfg := config.GetConfig()

	// Connect to database
	logger, _ := zap.NewDevelopment()
	db, err := database.NewPostgresDB(cfg.Database, logger)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate to create new tables
	log.Println("Migrating database...")
	err = db.DB.AutoMigrate(&domain.Tag{}, &domain.Category{}, &domain.Problem{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	// Seed tags
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
		{Name: "Breadth-First Search", Slug: "breadth-first-search"},
		{Name: "Tree", Slug: "tree"},
		{Name: "Binary Tree", Slug: "binary-tree"},
		{Name: "Graph", Slug: "graph"},
		{Name: "Linked List", Slug: "linked-list"},
		{Name: "Stack", Slug: "stack"},
		{Name: "Heap (Priority Queue)", Slug: "heap-priority-queue"},
		{Name: "Recursion", Slug: "recursion"},
		{Name: "Backtracking", Slug: "backtracking"},
		{Name: "Two Pointers", Slug: "two-pointers"},
		{Name: "Sliding Window", Slug: "sliding-window"},
	}

	log.Println("Seeding tags...")
	for _, tag := range tags {
		if err := db.DB.FirstOrCreate(&tag, domain.Tag{Slug: tag.Slug}).Error; err != nil {
			log.Printf("Failed to seed tag %s: %v", tag.Name, err)
		}
	}

	// Seed categories
	categories := []domain.Category{
		{Name: "Algorithms", Slug: "algorithms"},
		{Name: "Database", Slug: "database"},
		{Name: "Shell", Slug: "shell"},
		{Name: "JavaScript", Slug: "javascript"},
		{Name: "Concurrency", Slug: "concurrency"},
		{Name: "System Design", Slug: "system-design"},
	}

	log.Println("Seeding categories...")
	for _, cat := range categories {
		if err := db.DB.FirstOrCreate(&cat, domain.Category{Slug: cat.Slug}).Error; err != nil {
			log.Printf("Failed to seed category %s: %v", cat.Name, err)
		}
	}

	fmt.Println("Seeding completed successfully!")
}
