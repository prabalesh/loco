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

	// 1. Get Languages
	var languages []domain.Language
	db.DB.Find(&languages)
	log.Printf("Found %d languages", len(languages))

	pythonID := 0
	jsID := 0
	for _, l := range languages {
		if l.Slug == "python" {
			pythonID = l.ID
		}
		if l.Slug == "javascript" {
			jsID = l.ID
		}
	}

	if pythonID == 0 || jsID == 0 {
		log.Fatal("Python or JavaScript not found in languages table")
	}

	// 2. Get All Problems
	var problems []domain.Problem
	db.DB.Find(&problems)
	log.Printf("Found %d problems", len(problems))

	for _, p := range problems {
		log.Printf("Seeding languages for Problem: ID=%d, Title=%s", p.ID, p.Title)

		// Python
		plPy := domain.ProblemLanguage{
			ProblemID:    p.ID,
			LanguageID:   pythonID,
			FunctionCode: "def twoSum(nums, target):\n    # Write your code here\n    pass",
			MainCode:     "if __name__ == '__main__':\n    import sys\n    import json\n    nums = json.loads(sys.stdin.readline())\n    target = int(sys.stdin.readline())\n    print(json.dumps(twoSum(nums, target)))",
		}
		db.DB.FirstOrCreate(&plPy, domain.ProblemLanguage{ProblemID: p.ID, LanguageID: pythonID})

		// JavaScript
		plJs := domain.ProblemLanguage{
			ProblemID:    p.ID,
			LanguageID:   jsID,
			FunctionCode: "/**\n * @param {number[]} nums\n * @param {number} target\n * @return {number[]}\n */\nvar twoSum = function(nums, target) {\n    \n};",
			MainCode:     "const fs = require('fs');\nconst input = fs.readFileSync(0, 'utf8').split('\\n');\nconst nums = JSON.parse(input[0]);\nconst target = parseInt(input[1]);\nconsole.log(JSON.stringify(twoSum(nums, target)));",
		}
		db.DB.FirstOrCreate(&plJs, domain.ProblemLanguage{ProblemID: p.ID, LanguageID: jsID})

		// Sample Test Case
		var tcCount int64
		db.DB.Model(&domain.TestCase{}).Where("problem_id = ?", p.ID).Count(&tcCount)
		if tcCount == 0 {
			tc := domain.TestCase{
				ProblemID:      p.ID,
				Input:          "[2,7,11,15]\n9",
				ExpectedOutput: "[0,1]",
				IsSample:       true,
			}
			db.DB.Create(&tc)
		}
	}

	log.Println("Seeding completed successfully for all problems")
}
