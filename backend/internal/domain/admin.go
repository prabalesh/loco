package domain

type AdminAnalytics struct {
	TotalUsers    int `json:"total_users"`
	ActiveUsers   int `json:"active_users"`
	InactiveUsers int `json:"inactive_users"`
	VerifiedUsers int `json:"verified_users"`
}

type UpdateRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=user admin moderator"`
}

type UpdateStatusRequest struct {
	IsActive bool `json:"is_active"`
}
