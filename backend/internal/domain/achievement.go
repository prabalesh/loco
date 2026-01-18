package domain

import "time"

// Achievement represents a gamification achievement
type Achievement struct {
	ID             int    `json:"id" gorm:"primaryKey"`
	Slug           string `json:"slug" gorm:"unique;not null"`
	Name           string `json:"name" gorm:"not null"`
	Description    string `json:"description"`
	IconURL        string `json:"icon_url"` // Can be a frontend asset path or external URL
	XPReward       int    `json:"xp_reward" gorm:"not null"`
	Category       string `json:"category" gorm:"index"` // 'streak', 'solving', 'language', 'difficulty', 'misc'
	ConditionType  string `json:"condition_type"`        // 'count', 'streak', 'specific'
	ConditionValue string `json:"condition_value"`       // JSON or string value representing the target
}

// UserAchievement represents an unlocked achievement for a user
type UserAchievement struct {
	ID            int       `json:"id" gorm:"primaryKey"`
	UserID        int       `json:"user_id" gorm:"uniqueIndex:idx_user_achievement;not null"`
	AchievementID int       `json:"achievement_id" gorm:"uniqueIndex:idx_user_achievement;not null"`
	UnlockedAt    time.Time `json:"unlocked_at" gorm:"autoCreateTime"`

	Achievement Achievement `json:"achievement" gorm:"foreignKey:AchievementID"`
	User        User        `json:"-" gorm:"foreignKey:UserID"`
}

// AchievementRepository defines the interface for achievement persistence
type AchievementRepository interface {
	GetAll() ([]Achievement, error)
	GetBySlug(slug string) (*Achievement, error)
	GetUnlockedByUser(userID int) ([]UserAchievement, error)
	Unlock(userID int, achievementID int) error
	HasUnlocked(userID int, achievementID int) (bool, error)
	Create(achievement *Achievement) error
}
