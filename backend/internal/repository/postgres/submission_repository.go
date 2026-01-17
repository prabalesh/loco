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
	err := r.db.DB.Where("user_id = ? AND problem_id = ? AND is_admin_submission = false", userID, problemID).
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
