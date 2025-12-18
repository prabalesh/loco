package interfaces

import "github.com/prabalesh/loco/backend/internal/domain"

type ProblemRepository interface {
	Create(problem *domain.Problem) error
	Update(problem *domain.Problem) error
	Delete(id int) error
	GetByID(id int) (*domain.Problem, error)
	GetBySlug(slug string) (*domain.Problem, error)
	List(filters ProblemFilters) ([]*domain.Problem, int, error)
	SlugExists(slug string) (bool, error)
	UpdateCurrentStep(id int, newCurrentStep int) error
	UpdateStats(id int, acceptanceRate float64, totalSubmissions, totalAccepted int) error
	UpdateStatus(id int, status string) error
	CountProblems() (int, error)
	CountByStatus(status string) (int, error)
	CountByDifficulty(difficulty string) (int, error)
}

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
