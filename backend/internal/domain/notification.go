package domain

const AchievementEventChannel = "notifications:achievements"

// NotificationEvent types
const (
	EventAchievementUnlocked = "achievement_unlocked"
)

// AchievementUnlockedEvent data
type AchievementUnlockedEvent struct {
	UserID        int    `json:"user_id"`
	AchievementID int    `json:"achievement_id"`
	Slug          string `json:"slug"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	XPReward      int    `json:"xp_reward"`
	IconURL       string `json:"icon_url"`
}

// NotificationEvent represents a real-time event sent to the user
type NotificationEvent struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
