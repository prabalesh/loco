package postgres

import (
	"fmt"
	"time"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/pkg/database"
)

type achievementRepository struct {
	db *database.Database
}

func NewAchievementRepository(db *database.Database) domain.AchievementRepository {
	return &achievementRepository{db: db}
}

func (r *achievementRepository) GetAll() ([]domain.Achievement, error) {
	var achievements []domain.Achievement
	err := r.db.DB.Find(&achievements).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list achievements: %w", err)
	}
	return achievements, nil
}

func (r *achievementRepository) GetBySlug(slug string) (*domain.Achievement, error) {
	var achievement domain.Achievement
	err := r.db.DB.Where("slug = ?", slug).First(&achievement).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get achievement by slug: %w", err)
	}
	return &achievement, nil
}

func (r *achievementRepository) GetUnlockedByUser(userID int) ([]domain.UserAchievement, error) {
	var userAchievements []domain.UserAchievement
	err := r.db.DB.Preload("Achievement").Where("user_id = ?", userID).Find(&userAchievements).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list unlocked achievements: %w", err)
	}
	return userAchievements, nil
}

func (r *achievementRepository) Unlock(userID int, achievementID int) error {
	userAchievement := domain.UserAchievement{
		UserID:        userID,
		AchievementID: achievementID,
		UnlockedAt:    time.Now(),
	}
	err := r.db.DB.Create(&userAchievement).Error
	if err != nil {
		return fmt.Errorf("failed to unlock achievement: %w", err)
	}
	return nil
}

func (r *achievementRepository) HasUnlocked(userID int, achievementID int) (bool, error) {
	var count int64
	err := r.db.DB.Model(&domain.UserAchievement{}).
		Where("user_id = ? AND achievement_id = ?", userID, achievementID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check achievement status: %w", err)
	}
	return count > 0, nil
}

func (r *achievementRepository) Create(achievement *domain.Achievement) error {
	return r.db.DB.Create(achievement).Error
}
