package usecase

import (
	"errors"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/domain/dto"
	"go.uber.org/zap"
)

type ProblemLanguageUsecase struct {
	problemLanguageRepo domain.ProblemLanguageRepository
	problemRepo         domain.ProblemRepository
	languageRepo        domain.LanguageRepository
	logger              *zap.Logger
}

func NewProblemLanguageUsecase(
	problemLanguageRepo domain.ProblemLanguageRepository,
	problemRepo domain.ProblemRepository,
	languageRepo domain.LanguageRepository,
	logger *zap.Logger,
) *ProblemLanguageUsecase {
	return &ProblemLanguageUsecase{
		problemLanguageRepo: problemLanguageRepo,
		problemRepo:         problemRepo,
		languageRepo:        languageRepo,
		logger:              logger,
	}
}

func (u *ProblemLanguageUsecase) Create(problemID int, req *dto.CreateProblemLanguageRequest) (*domain.ProblemLanguage, error) {
	// Check if problem exists
	_, err := u.problemRepo.GetByID(problemID)
	if err != nil {
		return nil, errors.New("problem not found")
	}

	// Check if language exists
	lang, err := u.languageRepo.GetByID(req.LanguageID)
	if err != nil {
		return nil, errors.New("language not found")
	}

	// Check if already exists
	existing, _ := u.problemLanguageRepo.GetByProblemAndLanguage(problemID, req.LanguageID)
	if existing != nil {
		return nil, errors.New("language configuration already exists for this problem")
	}

	pl := &domain.ProblemLanguage{
		ProblemID:       problemID,
		LanguageID:      req.LanguageID,
		LanguageName:    lang.Name,
		LanguageVersion: lang.Version,
		FunctionCode:    req.FunctionCode,
		MainCode:        req.MainCode,
		SolutionCode:    req.SolutionCode,
	}

	if err := u.problemLanguageRepo.Create(pl); err != nil {
		u.logger.Error("Failed to create problem language", zap.Error(err))
		return nil, errors.New("failed to save language configuration")
	}

	return pl, nil
}

func (u *ProblemLanguageUsecase) ListByProblem(problemID int) ([]domain.ProblemLanguage, error) {
	return u.problemLanguageRepo.GetByProblemID(problemID)
}

func (u *ProblemLanguageUsecase) GetByProblemAndLanguage(problemID int, languageID int) (*domain.ProblemLanguage, error) {
	return u.problemLanguageRepo.GetByProblemAndLanguage(problemID, languageID)
}

func (u *ProblemLanguageUsecase) Update(problemID int, languageID int, req *dto.UpdateProblemLanguageRequest) (*domain.ProblemLanguage, error) {
	pl, err := u.problemLanguageRepo.GetByProblemAndLanguage(problemID, languageID)
	if err != nil {
		return nil, errors.New("language configuration not found")
	}

	if req.FunctionCode != "" {
		pl.FunctionCode = req.FunctionCode
	}
	if req.MainCode != "" {
		pl.MainCode = req.MainCode
	}
	if req.SolutionCode != "" {
		pl.SolutionCode = req.SolutionCode
	}

	if err := u.problemLanguageRepo.Update(pl); err != nil {
		u.logger.Error("Failed to update problem language", zap.Error(err))
		return nil, errors.New("failed to update language configuration")
	}

	return pl, nil
}

func (u *ProblemLanguageUsecase) Delete(problemID int, languageID int) error {
	return u.problemLanguageRepo.Delete(problemID, languageID)
}
