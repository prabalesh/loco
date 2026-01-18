package domain

import "time"

// UserRepository interface
type UserRepository interface {
	Create(user *User) error
	GetByEmail(email string) (*User, error)
	GetByUsername(username string) (*User, error)
	GetByID(id int) (*User, error)
	UpdatePassword(userID int, hashedPassword string) error
	UpdateRole(userID int, role string) error
	UpdateActiveStatus(userID int, isActive bool) error
	GetAll() ([]*User, error)
	Delete(id int) error

	// Verification
	UpdateVerificationToken(userID int, token string, expiresAt time.Time) error
	UpdateVerificationAttempts(userID int, attempts int) error
	UpdateLastSentAt(userID int, sentAt time.Time) error
	VerifyEmail(userID int) error
	GetByVerificationToken(token string) (*User, error)

	// Password Reset
	SetPasswordResetToken(userID int, token string, expiresAt time.Time) error
	ClearPasswordResetToken(userID int) error
	GetUserByResetToken(token string) (*User, error)
	GetByPasswordResetToken(token string) (*User, error)
	UpdatePasswordResetToken(userID int, token string, expiresAt time.Time, sentAt time.Time) error

	// Start
	CountUsers() (int, error)
	CountActiveUsers() (int, error)
	CountVerifiedUsers() (int, error)
	GetLeaderboard(limit int) ([]LeaderboardEntry, error)
	GetUserRank(userID int) (int, error)
}

// ProblemFilters for listing problems
type ProblemFilters struct {
	Page       int
	Limit      int
	Difficulty string
	Status     string
	Visibility string
	Search     string
	Tags       []string
	CreatedBy  *int
}

// ProblemRepository interface
type ProblemRepository interface {
	Create(problem *Problem) error
	GetByID(id int) (*Problem, error)
	GetBySlug(slug string) (*Problem, error)
	GetAll(limit, offset int, search string) ([]Problem, int64, error)
	List(filters ProblemFilters) ([]*Problem, int, error)
	Update(problem *Problem) error
	Delete(id int) error
	SlugExists(slug string) (bool, error)
	UpdateCurrentStep(id int, newCurrentStep int) error
	UpdateStats(id int, acceptanceRate float64, totalSubmissions, totalAccepted int) error
	IncrementStats(id int, isAccepted bool) error
	UpdateStatus(id int, status string) error
	UpdateVisibility(id int, visibility string) error
	CountProblems() (int, error)
	CountByStatus(status string) (int, error)
	CountByDifficulty(difficulty string) (int, error)
	ListTags() ([]Tag, error)
	ListCategories() ([]Category, error)
}

// LanguageRepository interface
type LanguageRepository interface {
	Create(language *Language) error
	GetByID(id int) (*Language, error)
	GetBySlug(slug string) (*Language, error)
	Update(language *Language) error
	Delete(id int) error
	GetAll() ([]Language, error)
	ListActive() ([]Language, error)
}

// TestCaseFilters for listing test cases
type TestCaseFilters struct {
	IsSample *bool
	Limit    int
	Page     int
}

// TestCaseRepository interface
type TestCaseRepository interface {
	Create(testCase *TestCase) error
	Update(testCase *TestCase) error
	Delete(id int) error
	GetByID(id int) (*TestCase, error)
	GetByProblemID(problemID int) ([]TestCase, error)
	List(problemID int, filters TestCaseFilters) ([]*TestCase, int, error)
	Exists(id int) (bool, error)
	UpdateOrderIndex(id int, orderIndex int) error
	DeleteByProblemID(problemID int) error
	CountByProblemID(problemID int) (int, error)
	GetSamples(problemID int) ([]TestCase, error)
}

// SubmissionRepository interface
type SubmissionRepository interface {
	Create(submission *Submission) error
	Update(submission *Submission) error
	GetByID(id int) (*Submission, error)
	ListByProblem(problemID int, limit, offset int) ([]Submission, error)
	ListByUserProblem(userID int, problemID int, limit, offset int) ([]Submission, error)
	ListByUser(userID int, limit, offset int) ([]Submission, error)
	ListByAdminUser(userID int, limit, offset int) ([]Submission, error)
	CountByUser(userID int) (int64, error)
	CountByUserProblem(userID int, problemID int) (int64, error)

	// Stats
	CountTotal() (int64, error)
	CountPending() (int64, error)
	CountCombinedStatus(status SubmissionStatus) (int64, error)
	CountAcceptedByUser(userID int) (int64, error)
	CountProblemsSolvedByUser(userID int) (int64, error)
	FindSolvedProblemsByUser(userID int, limit int) ([]Problem, error)
	CountSubmissionsLast24h() (int64, error)
	GetDailyStats(days int) ([]DailySubmissionStat, error)

	// Analytics
	GetTrendingProblems(limit int, days int) ([]TrendingProblem, error)
	GetLanguageStats() ([]LanguageStat, error)
	GetSolvedDistribution(userID int) ([]DifficultyStat, error)
	GetSubmissionHeatmap(userID int) ([]HeatmapEntry, error)

	// Queue monitoring
	GetOldestPending(limit int) ([]Submission, error)
	CountPendingBefore(createdAt time.Time) (int64, error)
}

// UserProblemStatsRepository interface
type UserProblemStatsRepository interface {
	Create(stats *UserProblemStats) error
	Update(stats *UserProblemStats) error
	Get(userID, problemID int) (*UserProblemStats, error)
	Upsert(stats *UserProblemStats) error
}

// ProblemLanguageRepository interface
type ProblemLanguageRepository interface {
	Create(problemLanguage *ProblemLanguage) error
	GetByProblemID(problemID int) ([]ProblemLanguage, error)
	GetByProblemAndLanguage(problemID int, languageID int) (*ProblemLanguage, error)
	Update(problemLanguage *ProblemLanguage) error
	Delete(problemID int, languageID int) error
}
