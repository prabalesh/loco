package migrations

import (
	"gorm.io/gorm"
)

func MigrateV3Indices(db *gorm.DB) error {
	migrations := []string{
		`CREATE INDEX IF NOT EXISTS idx_problems_status_visibility ON problems(status, visibility)`,
		`CREATE INDEX IF NOT EXISTS idx_problems_difficulty ON problems(difficulty)`,
		`CREATE INDEX IF NOT EXISTS idx_submissions_user_problem ON submissions(user_id, problem_id)`,
		`CREATE INDEX IF NOT EXISTS idx_user_problem_stats_user ON user_problem_stats(user_id)`,
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
