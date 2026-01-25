package domain

type TrendingProblem struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	Slug            string `json:"slug"`
	SubmissionCount int    `json:"submission_count"`
}

type LanguageStat struct {
	LanguageName string `json:"language_name"`
	Count        int    `json:"count"`
}

type DailySubmissionStat struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type DifficultyStat struct {
	Difficulty string `json:"difficulty"`
	Count      int    `json:"count"`
}

type HeatmapEntry struct {
	Date  string `json:"date"` // YYYY-MM-DD
	Count int    `json:"count"`
}

type ProblemStats struct {
	Total     int `json:"total"`
	Published int `json:"published"`
	Draft     int `json:"draft"`
	Easy      int `json:"easy"`
	Medium    int `json:"medium"`
	Hard      int `json:"hard"`
}
