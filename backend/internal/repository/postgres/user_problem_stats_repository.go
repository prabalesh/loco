package postgres

import (
	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/pkg/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userProblemStatsRepository struct {
	db *database.Database
}

func NewUserProblemStatsRepository(db *database.Database) *userProblemStatsRepository {
	return &userProblemStatsRepository{db: db}
}

func (r *userProblemStatsRepository) Create(stats *domain.UserProblemStats) error {
	return r.db.DB.Create(stats).Error
}

func (r *userProblemStatsRepository) Update(stats *domain.UserProblemStats) error {
	return r.db.DB.Model(&domain.UserProblemStats{}).
		Where("user_id = ? AND problem_id = ?", stats.UserID, stats.ProblemID).
		Updates(stats).Error
}

func (r *userProblemStatsRepository) Get(userID, problemID int) (*domain.UserProblemStats, error) {
	stats := &domain.UserProblemStats{}
	err := r.db.DB.Where("user_id = ? AND problem_id = ?", userID, problemID).Limit(1).Find(stats).Error
	if err != nil {
		return nil, err
	}
	if stats.UserID == 0 { // Check if no record found (assuming UserID is non-zero)
		return nil, nil
	}
	return stats, nil
}

func (r *userProblemStatsRepository) Upsert(stats *domain.UserProblemStats) error {
	return r.db.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}, {Name: "problem_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"attempts":           gorm.Expr("user_problem_stats.attempts + ?", 1),
			"status":             clause.Column{Table: "excluded", Name: "status"},
			"first_solved_at":    gorm.Expr("COALESCE(user_problem_stats.first_solved_at, excluded.first_solved_at)"),
			"best_submission_id": gorm.Expr("COALESCE(excluded.best_submission_id, user_problem_stats.best_submission_id)"),
			"updated_at":         clause.Column{Table: "excluded", Name: "updated_at"},
		}),
	}).Create(stats).Error
}

func (r *userProblemStatsRepository) GetStatuses(userID int, problemIDs []int) (map[int]string, error) {
	var results []struct {
		ProblemID int
		Status    string
	}
	err := r.db.DB.Model(&domain.UserProblemStats{}).
		Select("problem_id, status").
		Where("user_id = ? AND problem_id IN ?", userID, problemIDs).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	statusMap := make(map[int]string)
	for _, res := range results {
		statusMap[res.ProblemID] = res.Status
	}
	return statusMap, nil
}
