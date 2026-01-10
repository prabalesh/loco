package interfaces

import "github.com/prabalesh/loco/backend/internal/domain"

type ProblemLanguageRepository interface {
	Create(domain.ProblemLanguage) error
}
