package migrations

import (
	"github.com/prabalesh/loco/backend/internal/domain"
	"gorm.io/gorm"
)

func MigrateV2Schema(db *gorm.DB) error {
	// Auto-migrate new tables (GORM handles this safely)
	err := db.AutoMigrate(
		&domain.ProblemBoilerplate{},
		&domain.CustomType{},
		&domain.TypeImplementation{},
		&domain.ProblemReferenceSolution{},
	)
	if err != nil {
		return err
	}

	// Execute raw SQL for ALTER TABLE (safer than AutoMigrate for existing tables)
	migrations := []string{
		// Problems table
		`ALTER TABLE problems ADD COLUMN IF NOT EXISTS function_name VARCHAR(255)`,
		`ALTER TABLE problems ADD COLUMN IF NOT EXISTS return_type VARCHAR(100)`,
		`ALTER TABLE problems ADD COLUMN IF NOT EXISTS parameters JSONB`,
		`ALTER TABLE problems ADD COLUMN IF NOT EXISTS validation_type VARCHAR(50) DEFAULT 'EXACT'`,
		`ALTER TABLE problems ADD COLUMN IF NOT EXISTS validation_status VARCHAR(50) DEFAULT 'draft'`,
		`ALTER TABLE problems ADD COLUMN IF NOT EXISTS expected_time_complexity VARCHAR(50)`,
		`ALTER TABLE problems ADD COLUMN IF NOT EXISTS expected_space_complexity VARCHAR(50)`,
		`ALTER TABLE problems ADD COLUMN IF NOT EXISTS has_reference_solution BOOLEAN DEFAULT false`,

		// Test cases table
		`ALTER TABLE test_cases ADD COLUMN IF NOT EXISTS expected_outputs JSONB`,
		`ALTER TABLE test_cases ADD COLUMN IF NOT EXISTS input_size INTEGER`,
		`ALTER TABLE test_cases ADD COLUMN IF NOT EXISTS time_limit_ms INTEGER`,
		`ALTER TABLE test_cases ADD COLUMN IF NOT EXISTS memory_limit_mb INTEGER`,

		// Users table
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS is_bot BOOLEAN DEFAULT false`,

		// Add unique constraints for new tables
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_problem_boilerplate_unique ON problem_boilerplates(problem_id, language_id)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_type_implementation_unique ON type_implementations(custom_type_id, language_id)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_reference_solution_unique ON problem_reference_solutions(problem_id, language_id)`,
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, sql := range migrations {
			if err := tx.Exec(sql).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
