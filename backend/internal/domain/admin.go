package domain

type AdminAnalytics struct {
	TotalUsers         int                   `json:"total_users"`
	ActiveUsers        int                   `json:"active_users"`
	InactiveUsers      int                   `json:"inactive_users"`
	VerifiedUsers      int                   `json:"verified_users"`
	TotalSubmissions   int                   `json:"total_submissions"`
	PendingSubmissions int                   `json:"pending_submissions"`
	ActiveWorkers      int                   `json:"active_workers"`
	QueueSize          int64                 `json:"queue_size"`
	OldestPendingAge   int64                 `json:"oldest_pending_age_seconds"`
	QueueHealthStatus  string                `json:"queue_health_status"`
	SubmissionHistory  []DailySubmissionStat `json:"submission_history"`
}

type UpdateRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=user admin moderator"`
}

type UpdateStatusRequest struct {
	IsActive bool `json:"is_active"`
}
