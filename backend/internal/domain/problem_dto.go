package domain

// Request DTOs
type CreateProblemRequest struct {
	Title         string `json:"title"`
	Slug          string `json:"slug"`
	Description   string `json:"description"`
	Difficulty    string `json:"difficulty"`
	TimeLimit     int    `json:"time_limit"`
	MemoryLimit   int    `json:"memory_limit"`
	ValidatorType string `json:"validator_type"`
	InputFormat   string `json:"input_format"`
	OutputFormat  string `json:"output_format"`
	Constraints   string `json:"constraints"`
	Status        string `json:"status"`
	Visibility    string `json:"visibility"`
	IsActive      bool   `json:"is_active"`
	TagIDs        []int  `json:"tag_ids"`
	CategoryIDs   []int  `json:"category_ids"`
}

type UpdateProblemRequest struct {
	Title         string `json:"title"`
	Slug          string `json:"slug"`
	Description   string `json:"description"`
	Difficulty    string `json:"difficulty"`
	TimeLimit     int    `json:"time_limit"`
	MemoryLimit   int    `json:"memory_limit"`
	ValidatorType string `json:"validator_type"`
	InputFormat   string `json:"input_format"`
	OutputFormat  string `json:"output_format"`
	Constraints   string `json:"constraints"`
	Status        string `json:"status"`
	Visibility    string `json:"visibility"`
	IsActive      *bool  `json:"is_active"`
	TagIDs        []int  `json:"tag_ids"`
	CategoryIDs   []int  `json:"category_ids"`
}

type ListProblemsRequest struct {
	Page       int      `json:"page"`
	Limit      int      `json:"limit"`
	Difficulty string   `json:"difficulty"`
	Search     string   `json:"search"`
	Tags       []string `json:"tags"`
	Categories []string `json:"categories"`
}

type AdminListProblemsRequest struct {
	Page       int      `json:"page"`
	Limit      int      `json:"limit"`
	Difficulty string   `json:"difficulty"`
	Status     string   `json:"status"`
	Visibility string   `json:"visibility"`
	Search     string   `json:"search"`
	Tags       []string `json:"tags"`
	Categories []string `json:"categories"`
}

type ProblemStats struct {
	Total     int `json:"total"`
	Published int `json:"published"`
	Draft     int `json:"draft"`
	Easy      int `json:"easy"`
	Medium    int `json:"medium"`
	Hard      int `json:"hard"`
}
