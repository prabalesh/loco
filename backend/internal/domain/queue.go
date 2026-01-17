package domain

// QueueHealthStatus represents the health of the queue system
type QueueHealthStatus string

const (
	QueueHealthHealthy  QueueHealthStatus = "healthy"
	QueueHealthWarning  QueueHealthStatus = "warning"
	QueueHealthCritical QueueHealthStatus = "critical"
)

// QueueStatus represents the overall status of the submission queue
type QueueStatus struct {
	QueueSize         int64             `json:"queue_size"`
	ActiveWorkers     int               `json:"active_workers"`
	OldestPendingAge  int64             `json:"oldest_pending_age_seconds"` // Age in seconds
	HealthStatus      QueueHealthStatus `json:"health_status"`
	WarningMessage    string            `json:"warning_message,omitempty"`
	PendingCount      int64             `json:"pending_count"`
	EstimatedWaitTime int64             `json:"estimated_wait_time_seconds,omitempty"` // Estimated wait in seconds
}

// SubmissionQueueInfo represents queue-specific info for a submission
type SubmissionQueueInfo struct {
	QueuePosition     int   `json:"queue_position,omitempty"`
	EstimatedWaitTime int64 `json:"estimated_wait_time_seconds,omitempty"`
	WorkersActive     bool  `json:"workers_active"`
}
