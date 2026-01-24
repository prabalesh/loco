package postgres

import (
	"fmt"
	"time"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/pkg/database"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type submissionRepository struct {
	db *database.Database
}

func NewSubmissionRepository(db *database.Database) domain.SubmissionRepository {
	return &submissionRepository{db: db}
}

func (r *submissionRepository) Create(submission *domain.Submission) error {
	return r.db.DB.Create(submission).Error
}

func (r *submissionRepository) Update(submission *domain.Submission) error {
	return r.db.DB.Save(submission).Error
}

func (r *submissionRepository) GetByID(id int) (*domain.Submission, error) {
	var submission domain.Submission
	err := r.db.DB.
		Where("id = ?", id).
		Preload("User").
		Preload("Problem").
		Preload("Language").
		First(&submission).
		Error

	if err != nil {
		if err != gorm.ErrRecordNotFound {
			r.db.Logger.Error("Database error in GetByID", zap.Int("id", id), zap.Error(err))
		}
		return nil, err
	}
	return &submission, nil
}

func (r *submissionRepository) ListByProblem(problemID int, limit, offset int) ([]domain.Submission, error) {
	var submissions []domain.Submission
	err := r.db.DB.Where("problem_id = ?", problemID).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&submissions).Error
	return submissions, err
}

func (r *submissionRepository) ListByUser(userID int, limit, offset int) ([]domain.Submission, error) {
	var submissions []domain.Submission
	err := r.db.DB.Where("user_id = ? AND is_admin_submission = false", userID).
		Preload("Problem").
		Preload("Language").
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&submissions).Error
	return submissions, err
}

func (r *submissionRepository) ListByUserProblem(userID int, problemID int, limit, offset int) ([]domain.Submission, error) {
	fmt.Println("Reached hereeeeeeee...***********8")
	fmt.Println("User ID: ", userID)
	fmt.Println("Problem ID: ", problemID)
	fmt.Println("Limit: ", limit)
	fmt.Println("Offset: ", offset)
	var submissions []domain.Submission
	err := r.db.DB.Where("user_id = ? AND problem_id = ? AND is_admin_submission = false AND is_run_only = false", userID, problemID).
		Preload("Language").
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&submissions).Error
	return submissions, err
}

func (r *submissionRepository) ListByAdminUser(userID int, limit, offset int) ([]domain.Submission, error) {
	var submissions []domain.Submission
	err := r.db.DB.Where("user_id = ? AND is_admin_submission = true", userID).
		Preload("Problem").
		Preload("Language").
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&submissions).Error
	return submissions, err
}

func (r *submissionRepository) CountByUser(userID int) (int64, error) {
	var count int64
	err := r.db.DB.Model(&domain.Submission{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *submissionRepository) CountByUserProblem(userID int, problemID int) (int64, error) {
	var count int64
	err := r.db.DB.Model(&domain.Submission{}).Where("user_id = ? AND problem_id = ? AND is_admin_submission = false", userID, problemID).Count(&count).Error
	return count, err
}

func (r *submissionRepository) FindSolvedProblemsByUser(userID int, limit int) ([]domain.Problem, error) {
	var problems []domain.Problem
	err := r.db.DB.Raw(`
		SELECT p.* 
		FROM problems p
		JOIN (
			SELECT DISTINCT ON (problem_id) problem_id, created_at
			FROM submissions
			WHERE user_id = ? AND status = ?
			ORDER BY problem_id, created_at DESC
		) s ON p.id = s.problem_id
		ORDER BY s.created_at DESC
		LIMIT ?
	`, userID, domain.SubmissionStatusAccepted, limit).Scan(&problems).Error
	return problems, err
}

// Stats implementations
func (r *submissionRepository) CountTotal() (int64, error) {
	var count int64
	err := r.db.DB.Model(&domain.Submission{}).Count(&count).Error
	return count, err
}

func (r *submissionRepository) CountPending() (int64, error) {
	var count int64
	err := r.db.DB.Model(&domain.Submission{}).Where("status = ?", domain.SubmissionStatusPending).Count(&count).Error
	return count, err
}

func (r *submissionRepository) GetOldestPending(limit int) ([]domain.Submission, error) {
	var submissions []domain.Submission
	err := r.db.DB.Where("status = ?", domain.SubmissionStatusPending).
		Order("created_at ASC").
		Limit(limit).
		Find(&submissions).Error
	return submissions, err
}

func (r *submissionRepository) CountPendingBefore(createdAt time.Time) (int64, error) {
	var count int64
	err := r.db.DB.Model(&domain.Submission{}).
		Where("status = ? AND created_at < ?", domain.SubmissionStatusPending, createdAt).
		Count(&count).Error
	return count, err
}

func (r *submissionRepository) CountCombinedStatus(status domain.SubmissionStatus) (int64, error) {
	var count int64
	err := r.db.DB.Model(&domain.Submission{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

func (r *submissionRepository) CountAcceptedByUser(userID int) (int64, error) {
	var count int64
	err := r.db.DB.Model(&domain.Submission{}).Where("user_id = ? AND status = ?", userID, domain.SubmissionStatusAccepted).Count(&count).Error
	return count, err
}

func (r *submissionRepository) CountProblemsSolvedByUser(userID int) (int64, error) {
	var count int64
	err := r.db.DB.Model(&domain.Submission{}).
		Where("user_id = ? AND status = ?", userID, domain.SubmissionStatusAccepted).
		Distinct("problem_id").
		Count(&count).Error
	return count, err
}

func (r *submissionRepository) CountSubmissionsLast24h() (int64, error) {
	var count int64
	// Use raw SQL or database specific syntax for time interval. Postgres
	err := r.db.DB.Model(&domain.Submission{}).
		Where("created_at >= NOW() - INTERVAL '24 hours'").
		Count(&count).Error
	return count, err
}

func (r *submissionRepository) GetDailyStats(days int) ([]domain.DailySubmissionStat, error) {
	var stats []domain.DailySubmissionStat
	// Postgres SQL for daily aggregation
	query := `
		SELECT 
			TO_CHAR(created_at, 'YYYY-MM-DD') as date, 
			COUNT(*) as count 
		FROM submissions 
		WHERE created_at >= NOW() - make_interval(days => ?)
		GROUP BY date 
		ORDER BY date ASC
	`
	err := r.db.DB.Raw(query, days).Scan(&stats).Error
	return stats, err
}
func (r *submissionRepository) GetTrendingProblems(limit int, days int) ([]domain.TrendingProblem, error) {
	var results []domain.TrendingProblem
	query := `
		SELECT 
			p.id, 
			p.title, 
			p.slug, 
			COUNT(s.id) as submission_count
		FROM problems p
		JOIN submissions s ON p.id = s.problem_id
		WHERE s.created_at >= NOW() - make_interval(days => ?)
		  AND s.is_admin_submission = false
		GROUP BY p.id, p.title, p.slug
		ORDER BY submission_count DESC
		LIMIT ?
	`
	err := r.db.DB.Raw(query, days, limit).Scan(&results).Error
	return results, err
}

func (r *submissionRepository) GetLanguageStats() ([]domain.LanguageStat, error) {
	var stats []domain.LanguageStat
	query := `
		SELECT 
			l.name as language_name, 
			COUNT(s.id) as count
		FROM languages l
		JOIN submissions s ON l.id = s.language_id
		WHERE s.is_admin_submission = false
		GROUP BY l.name
		ORDER BY count DESC
	`
	err := r.db.DB.Raw(query).Scan(&stats).Error
	return stats, err
}
func (r *submissionRepository) GetSolvedDistribution(userID int) ([]domain.DifficultyStat, error) {
	var stats []domain.DifficultyStat
	query := `
		SELECT 
			p.difficulty, 
			COUNT(DISTINCT p.id) as count
		FROM problems p
		JOIN submissions s ON p.id = s.problem_id
		WHERE s.user_id = ? AND s.status = 'Accepted' AND s.is_admin_submission = false
		GROUP BY p.difficulty
	`
	err := r.db.DB.Raw(query, userID).Scan(&stats).Error
	return stats, err
}

func (r *submissionRepository) GetSubmissionHeatmap(userID int) ([]domain.HeatmapEntry, error) {
	var heatmap []domain.HeatmapEntry
	query := `
		SELECT 
			TO_CHAR(created_at, 'YYYY-MM-DD') as date, 
			COUNT(*) as count
		FROM submissions
		WHERE user_id = ? AND is_admin_submission = false AND created_at >= NOW() - INTERVAL '1 year'
		GROUP BY date
		ORDER BY date ASC
	`
	err := r.db.DB.Raw(query, userID).Scan(&heatmap).Error
	return heatmap, err
}

func (r *submissionRepository) GetCurrentStreak(userID int) (int, error) {
	var dates []time.Time
	query := `
		SELECT DISTINCT DATE_TRUNC('day', created_at) as day
		FROM submissions
		WHERE user_id = ? AND is_admin_submission = false
		ORDER BY day DESC
	`
	err := r.db.DB.Raw(query, userID).Scan(&dates).Error
	if err != nil {
		return 0, err
	}

	if len(dates) == 0 {
		return 0, nil
	}

	streak := 0
	now := time.Now().Truncate(24 * time.Hour)
	lastDate := dates[0].Truncate(24 * time.Hour)

	// Check if the last submission was today or yesterday
	if lastDate.Equal(now) || lastDate.Equal(now.AddDate(0, 0, -1)) {
		streak = 1
		for i := 1; i < len(dates); i++ {
			prevDate := dates[i-1].Truncate(24 * time.Hour)
			currDate := dates[i].Truncate(24 * time.Hour)

			if prevDate.AddDate(0, 0, -1).Equal(currDate) {
				streak++
			} else {
				break
			}
		}
	}

	return streak, nil
}
