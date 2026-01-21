package codegen

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/prabalesh/loco/backend/internal/domain"
)

type BoilerplateService struct {
	boilerplateRepo domain.BoilerplateRepository
	languageRepo    domain.LanguageRepository
	codeGenService  *CodeGenService
}

func NewBoilerplateService(
	boilerplateRepo domain.BoilerplateRepository,
	languageRepo domain.LanguageRepository,
	codeGenService *CodeGenService,
) *BoilerplateService {
	return &BoilerplateService{
		boilerplateRepo: boilerplateRepo,
		languageRepo:    languageRepo,
		codeGenService:  codeGenService,
	}
}

func (s *BoilerplateService) GenerateAllBoilerplatesForProblem(problem *domain.Problem) error {
	if problem.FunctionName == nil || *problem.FunctionName == "" {
		return fmt.Errorf("problem has no function name")
	}
	if problem.ReturnType == nil || *problem.ReturnType == "" {
		return fmt.Errorf("problem has no return type")
	}
	if problem.Parameters == nil {
		return fmt.Errorf("problem has no parameters")
	}

	var parameters []Parameter
	if err := json.Unmarshal([]byte(*problem.Parameters), &parameters); err != nil {
		return fmt.Errorf("failed to parse parameters: %w", err)
	}

	signature := ProblemSignature{
		FunctionName: *problem.FunctionName,
		ReturnType:   *problem.ReturnType,
		Parameters:   parameters,
	}

	languages, err := s.languageRepo.ListActive()
	if err != nil {
		return fmt.Errorf("failed to list active languages: %w", err)
	}

	for _, lang := range languages {
		err := s.GenerateBoilerplateForLanguage(problem.ID, lang.ID, signature, lang.Slug)
		if err != nil {
			// Log error but continue with other languages
			fmt.Printf("Warning: failed to generate boilerplate for language %s (ID: %d): %v\n", lang.Name, lang.ID, err)
			continue
		}
	}

	return nil
}

func (s *BoilerplateService) GenerateBoilerplateForLanguage(problemID, languageID int, signature ProblemSignature, languageSlug string) error {
	stubCode, err := s.codeGenService.GenerateStubCode(signature, languageSlug)
	if err != nil {
		return fmt.Errorf("failed to generate stub code: %w", err)
	}

	harnessTemplate, err := s.codeGenService.GenerateTestHarness(signature, "{USER_CODE}", languageSlug)
	if err != nil {
		return fmt.Errorf("failed to generate test harness template: %w", err)
	}

	exists, err := s.boilerplateRepo.Exists(problemID, languageID)
	if err != nil {
		return err
	}

	if exists {
		existing, err := s.boilerplateRepo.GetByProblemAndLanguage(problemID, languageID)
		if err != nil {
			return err
		}
		existing.StubCode = stubCode
		existing.TestHarnessTemplate = harnessTemplate
		return s.boilerplateRepo.Update(existing)
	}

	boilerplate := &domain.ProblemBoilerplate{
		ProblemID:           problemID,
		LanguageID:          languageID,
		StubCode:            stubCode,
		TestHarnessTemplate: harnessTemplate,
	}

	return s.boilerplateRepo.Create(boilerplate)
}

func (s *BoilerplateService) GetStubCode(problemID, languageID int) (string, error) {
	bp, err := s.boilerplateRepo.GetByProblemAndLanguage(problemID, languageID)
	if err != nil {
		return "", err
	}
	return bp.StubCode, nil
}

func (s *BoilerplateService) GetTestHarnessTemplate(problemID, languageID int) (string, error) {
	bp, err := s.boilerplateRepo.GetByProblemAndLanguage(problemID, languageID)
	if err != nil {
		return "", err
	}
	return bp.TestHarnessTemplate, nil
}

func (s *BoilerplateService) InjectUserCodeIntoHarness(template, userCode string) string {
	return strings.Replace(template, "{USER_CODE}", userCode, 1)
}

func (s *BoilerplateService) RegenerateBoilerplatesForProblem(problem *domain.Problem) error {
	// We use GenerateAllBoilerplatesForProblem because it handles updates if exists.
	// But if the signature changed, it might be safer to refresh all.
	// The GenerateAllBoilerplatesForProblem already does Update if Exists.
	return s.GenerateAllBoilerplatesForProblem(problem)
}

func (s *BoilerplateService) GetBoilerplateStats(problemID int) (map[string]interface{}, error) {
	boilerplates, err := s.boilerplateRepo.GetByProblemID(problemID)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_languages": len(boilerplates),
		"languages":       []string{},
	}

	languages := []string{}
	for _, bp := range boilerplates {
		if bp.Language.Name != "" {
			languages = append(languages, bp.Language.Name)
		}
	}
	stats["languages"] = languages

	return stats, nil
}
