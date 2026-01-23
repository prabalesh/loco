package codegen

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/prabalesh/loco/backend/internal/domain"
	"gorm.io/datatypes"
)

type BoilerplateService struct {
	boilerplateRepo domain.BoilerplateRepository
	languageRepo    domain.LanguageRepository
	testCaseRepo    domain.TestCaseRepository
	codeGenService  *CodeGenService
}

func NewBoilerplateService(
	boilerplateRepo domain.BoilerplateRepository,
	languageRepo domain.LanguageRepository,
	testCaseRepo domain.TestCaseRepository,
	codeGenService *CodeGenService,
) *BoilerplateService {
	return &BoilerplateService{
		boilerplateRepo: boilerplateRepo,
		languageRepo:    languageRepo,
		testCaseRepo:    testCaseRepo,
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

	var parameters []domain.SchemaParameter
	if err := json.Unmarshal([]byte(*problem.Parameters), &parameters); err != nil {
		return fmt.Errorf("failed to parse parameters: %w", err)
	}

	signature := domain.ProblemSchema{
		FunctionName: *problem.FunctionName,
		ReturnType:   domain.GenericType(*problem.ReturnType),
		Parameters:   parameters,
	}

	languages, err := s.languageRepo.ListActive()
	if err != nil {
		return fmt.Errorf("failed to list active languages: %w", err)
	}

	testCases, err := s.testCaseRepo.GetByProblemID(problem.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch test cases: %w", err)
	}

	var errors []string
	for _, lang := range languages {
		err := s.GenerateBoilerplateForLanguage(problem.ID, lang.ID, signature, lang.Slug, testCases, problem.ValidationType)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", lang.Name, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to generate boilerplates for some languages: %s", strings.Join(errors, "; "))
	}

	return nil
}

func (s *BoilerplateService) GenerateBoilerplateForLanguage(problemID, languageID int, signature domain.ProblemSchema, languageSlug string, testCases []domain.TestCase, validationType string) error {
	stubCode, err := s.codeGenService.GenerateStubCode(signature, languageSlug)
	if err != nil {
		return fmt.Errorf("failed to generate stub code: %w", err)
	}

	harnessTemplate, err := s.codeGenService.GenerateTestHarness(signature, "{USER_CODE}", languageSlug, testCases, validationType)
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
		harnessJSON, _ := json.Marshal(map[string]string{languageSlug: harnessTemplate})
		existing.TestHarnessTemplate = datatypes.JSON(harnessJSON)
		return s.boilerplateRepo.Update(existing)
	}

	harnessJSON, _ := json.Marshal(map[string]string{languageSlug: harnessTemplate})
	boilerplate := &domain.ProblemBoilerplate{
		ProblemID:           problemID,
		LanguageID:          languageID,
		StubCode:            stubCode,
		TestHarnessTemplate: datatypes.JSON(harnessJSON),
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

	var templates map[string]string
	if err := json.Unmarshal(bp.TestHarnessTemplate, &templates); err != nil {
		return "", fmt.Errorf("failed to parse harness template: %w", err)
	}

	// The map should contain the template for the specific language slug
	// However, we don't have the slug here easily without another query or passing it in.
	// But since it's a map with one key (usually), we can just take the first value or logic it out.
	// Better: the caller should know what they want.
	// For now, let's just return the first string value found in the map.
	for _, template := range templates {
		return template, nil
	}

	return "", fmt.Errorf("no template found in JSON")
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

func (s *BoilerplateService) GetBoilerplatesByProblemID(problemID int) ([]domain.ProblemBoilerplate, error) {
	return s.boilerplateRepo.GetByProblemID(problemID)
}

func (s *BoilerplateService) DeleteBoilerplatesByProblemID(problemID int) error {
	return s.boilerplateRepo.DeleteByProblemID(problemID)
}
