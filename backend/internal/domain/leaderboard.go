package domain

type LeaderboardEntry struct {
	Rank             int     `json:"rank"`
	UserID           int     `json:"user_id"`
	Username         string  `json:"username"`
	ProblemsSolved   int     `json:"problems_solved"`
	TotalSubmissions int     `json:"total_submissions"`
	AcceptanceRate   float64 `json:"acceptance_rate"`
}
