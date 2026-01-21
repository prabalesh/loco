package problem

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/services/codegen"
	"gorm.io/datatypes"
)

type ProblemService struct {
	problemRepo        domain.ProblemRepository
	testCaseRepo       domain.TestCaseRepository
	customTypeRepo     domain.CustomTypeRepository
	referenceRepo      domain.ReferenceSolutionRepository
	boilerplateService *codegen.BoilerplateService
}

func NewProblemService(
	problemRepo domain.ProblemRepository,
	testCaseRepo domain.TestCaseRepository,
	customTypeRepo domain.CustomTypeRepository,
	referenceRepo domain.ReferenceSolutionRepository,
	boilerplateService *codegen.BoilerplateService,
) *ProblemService {
	return &ProblemService{
		problemRepo:        problemRepo,
		testCaseRepo:       testCaseRepo,
		customTypeRepo:     customTypeRepo,
		referenceRepo:      referenceRepo,
		boilerplateService: boilerplateService,
	}
}

type CreateProblemRequest struct {
	Title                   string              `json:"title"`
	Description             string              `json:"description"`
	Difficulty              string              `json:"difficulty"`
	CategoryIDs             []int               `json:"category_ids"`
	TagIDs                  []int               `json:"tag_ids"`
	FunctionName            string              `json:"function_name"`
	ReturnType              string              `json:"return_type"`
	Parameters              []codegen.Parameter `json:"parameters"`
	ValidationType          string              `json:"validation_type"`
	ExpectedTimeComplexity  string              `json:"expected_time_complexity"`
	ExpectedSpaceComplexity string              `json:"expected_space_complexity"`
	TestCases               []TestCaseInput     `json:"test_cases"`
}

type TestCaseInput struct {
	Input          interface{} `json:"input"` // Array of parameter values or single value if 1 param
	ExpectedOutput interface{} `json:"expected_output"`
	IsSample       bool        `json:"is_sample"`
	InputSize      *int        `json:"input_size"`
	TimeLimitMs    *int        `json:"time_limit_ms"`
	MemoryLimitMb  *int        `json:"memory_limit_mb"`
}

// CreateProblem creates a new problem with auto-generation
func (s *ProblemService) CreateProblem(req CreateProblemRequest, createdBy int) (*domain.Problem, error) {
	// Validate request
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Generate slug
	slug := s.GenerateSlug(req.Title)

	// Ensure slug is unique
	slug, err := s.ensureUniqueSlug(slug, 0)
	if err != nil {
		return nil, err
	}

	// Convert parameters to JSON
	paramsJSON, err := json.Marshal(req.Parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal parameters: %w", err)
	}
	paramsData := datatypes.JSON(paramsJSON)

	// Create problem model
	problem := &domain.Problem{
		Title:                   req.Title,
		Slug:                    slug,
		Description:             req.Description,
		Difficulty:              req.Difficulty,
		FunctionName:            &req.FunctionName,
		ReturnType:              &req.ReturnType,
		Parameters:              &paramsData,
		ValidationType:          req.ValidationType,
		ValidationStatus:        "draft",
		ExpectedTimeComplexity:  &req.ExpectedTimeComplexity,
		ExpectedSpaceComplexity: &req.ExpectedSpaceComplexity,
		Status:                  "draft",
		Visibility:              "private",
		CreatedBy:               &createdBy,
	}

	// Save problem
	if err := s.problemRepo.Create(problem); err != nil {
		return nil, fmt.Errorf("failed to create problem: %w", err)
	}

	// Add categories
	if len(req.CategoryIDs) > 0 {
		categories := make([]domain.Category, len(req.CategoryIDs))
		for i, id := range req.CategoryIDs {
			categories[i] = domain.Category{ID: id}
		}
		problem.Categories = categories
	}

	// Add tags
	if len(req.TagIDs) > 0 {
		tags := make([]domain.Tag, len(req.TagIDs))
		for i, id := range req.TagIDs {
			tags[i] = domain.Tag{ID: id}
		}
		problem.Tags = tags
	}

	// Save associations if any
	if len(req.CategoryIDs) > 0 || len(req.TagIDs) > 0 {
		if err := s.problemRepo.Update(problem); err != nil {
			return nil, fmt.Errorf("failed to update problem associations: %w", err)
		}
	}

	// Create test cases
	if err := s.createTestCases(problem.ID, req.TestCases); err != nil {
		return nil, fmt.Errorf("failed to create test cases: %w", err)
	}

	// Auto-generate boilerplates (async recommended but sync for simplicity)
	if err := s.boilerplateService.GenerateAllBoilerplatesForProblem(problem); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: Failed to generate boilerplates: %v\n", err)
	}

	return problem, nil
}

// validateCreateRequest validates the problem creation request
func (s *ProblemService) validateCreateRequest(req CreateProblemRequest) error {
	// Title
	if len(req.Title) < 5 || len(req.Title) > 200 {
		return errors.New("title must be between 5 and 200 characters")
	}

	// Description
	if len(req.Description) < 20 {
		return errors.New("description must be at least 20 characters")
	}

	// Difficulty
	validDifficulties := map[string]bool{"easy": true, "medium": true, "hard": true}
	if !validDifficulties[req.Difficulty] {
		return errors.New("difficulty must be easy, medium, or hard")
	}

	// Function name
	if req.FunctionName == "" {
		return errors.New("function_name is required")
	}
	if !isValidIdentifier(req.FunctionName) {
		return errors.New("function_name must be a valid identifier")
	}

	// Return type
	if req.ReturnType == "" {
		return errors.New("return_type is required")
	}

	// Parameters
	if len(req.Parameters) == 0 {
		return errors.New("at least one parameter is required")
	}
	for _, param := range req.Parameters {
		if param.Name == "" {
			return errors.New("parameter name cannot be empty")
		}
		if param.Type == "" {
			return errors.New("parameter type cannot be empty")
		}
		if param.IsCustom {
			// Validate custom type exists
			if _, err := s.customTypeRepo.GetByName(param.Type); err != nil {
				return fmt.Errorf("invalid custom type: %s", param.Type)
			}
		}
	}

	// Validation type
	validTypes := map[string]bool{
		"EXACT": true, "UNORDERED": true, "SUBSET": true, "ANY_MATCH": true,
	}
	if req.ValidationType == "" {
		req.ValidationType = "EXACT"
	}
	if !validTypes[req.ValidationType] {
		return errors.New("invalid validation_type")
	}

	// Test cases
	if len(req.TestCases) == 0 {
		return errors.New("at least one test case is required")
	}

	// Check at least one public test case
	hasPublic := false
	for _, tc := range req.TestCases {
		if tc.IsSample {
			hasPublic = true
			break
		}
	}
	if !hasPublic {
		return errors.New("at least one public test case (is_sample=true) is required")
	}

	return nil
}

// GenerateSlug generates a URL-friendly slug from title
func (s *ProblemService) GenerateSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces and special chars with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	return slug
}

// ensureUniqueSlug ensures slug is unique by appending number if needed
func (s *ProblemService) ensureUniqueSlug(slug string, excludeID int) (string, error) {
	originalSlug := slug
	counter := 2

	for {
		exists, err := s.problemRepo.SlugExists(slug, excludeID)
		if err != nil {
			return "", err
		}
		if !exists {
			return slug, nil
		}

		// Append counter
		slug = fmt.Sprintf("%s-%d", originalSlug, counter)
		counter++

		if counter > 100 {
			return "", errors.New("failed to generate unique slug")
		}
	}
}

// createTestCases creates test cases for a problem
func (s *ProblemService) createTestCases(problemID int, testCaseInputs []TestCaseInput) error {
	testCases := []domain.TestCase{}

	for i, tcInput := range testCaseInputs {
		// Convert input to JSON
		inputJSON, err := json.Marshal(tcInput.Input)
		if err != nil {
			return fmt.Errorf("invalid test case input at index %d: %w", i, err)
		}

		// Convert expected output to JSON
		outputJSON, err := json.Marshal(tcInput.ExpectedOutput)
		if err != nil {
			return fmt.Errorf("invalid test case output at index %d: %w", i, err)
		}

		testCase := domain.TestCase{
			ProblemID:      problemID,
			Input:          string(inputJSON),
			ExpectedOutput: string(outputJSON),
			IsSample:       tcInput.IsSample,
			InputSize:      tcInput.InputSize,
			TimeLimitMs:    tcInput.TimeLimitMs,
			MemoryLimitMb:  tcInput.MemoryLimitMb,
			OrderIndex:     i,
		}

		testCases = append(testCases, testCase)
	}

	return s.testCaseRepo.CreateMany(testCases)
}

// isValidIdentifier checks if string is valid identifier (alphanumeric + underscore)
func isValidIdentifier(s string) bool {
	match, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, s)
	return match
}

// GetProblemDetail fetches problem with public test cases
func (s *ProblemService) GetProblemDetail(slug string, isAdmin bool) (map[string]interface{}, error) {
	problem, err := s.problemRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	// Get test cases (public only for non-admin)
	var testCases []domain.TestCase
	if isAdmin {
		testCases, err = s.testCaseRepo.GetByProblemID(problem.ID)
	} else {
		testCases, err = s.testCaseRepo.GetSamples(problem.ID)
	}
	if err != nil {
		return nil, err
	}

	// Build response
	response := map[string]interface{}{
		"id":                        problem.ID,
		"title":                     problem.Title,
		"slug":                      problem.Slug,
		"description":               problem.Description,
		"difficulty":                problem.Difficulty,
		"function_name":             problem.FunctionName,
		"return_type":               problem.ReturnType,
		"parameters":                problem.Parameters,
		"validation_type":           problem.ValidationType,
		"expected_time_complexity":  problem.ExpectedTimeComplexity,
		"expected_space_complexity": problem.ExpectedSpaceComplexity,
		"categories":                problem.Categories,
		"tags":                      problem.Tags,
		"test_cases":                testCases,
		"acceptance_rate":           problem.AcceptanceRate,
		"total_submissions":         problem.TotalSubmissions,
	}

	return response, nil
}

func (s *ProblemService) GetAllProblems(filters map[string]interface{}, page, limit int) ([]*domain.Problem, int, error) {
	domainFilters := domain.ProblemFilters{
		Page:  page,
		Limit: limit,
	}

	if d, ok := filters["difficulty"].(string); ok {
		domainFilters.Difficulty = d
	}
	if s, ok := filters["status"].(string); ok {
		domainFilters.Status = s
	}

	return s.problemRepo.List(domainFilters)
}

func (s *ProblemService) GetByID(id int) (*domain.Problem, error) {
	return s.problemRepo.GetByID(id)
}

func (s *ProblemService) GetAdminByID(id int) (*domain.Problem, error) {
	problem, err := s.problemRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Fetch additional admin info
	testCases, _ := s.testCaseRepo.GetByProblemID(id)
	problem.TestCases = testCases

	boilerplates, _ := s.boilerplateService.GetBoilerplatesByProblemID(id)
	problem.Boilerplates = boilerplates

	referenceSolutions, _ := s.referenceRepo.GetAllByProblemID(id)
	problem.ReferenceSolutions = referenceSolutions

	return problem, nil
}

func (s *ProblemService) DeleteProblem(id int) error {
	// Delete related test cases first (optional if DB has cascade, but safer here)
	s.testCaseRepo.DeleteByProblemID(id)
	s.boilerplateService.DeleteBoilerplatesByProblemID(id)
	// Delete problem
	return s.problemRepo.Delete(id)
}

func (s *ProblemService) PublishProblem(id int) error {
	problem, err := s.problemRepo.GetByID(id)
	if err != nil {
		return err
	}

	if problem.ValidationStatus != "validated" {
		return errors.New("problem must be validated before publishing")
	}

	problem.Status = "published"
	problem.Visibility = "public"

	return s.problemRepo.Update(problem)
}

func (s *ProblemService) RegenerateBoilerplates(id int) error {
	problem, err := s.problemRepo.GetByID(id)
	if err != nil {
		return err
	}

	return s.boilerplateService.RegenerateBoilerplatesForProblem(problem)
}

func (s *ProblemService) GetCustomTypes() ([]domain.CustomType, error) {
	return s.customTypeRepo.GetAll()
}

// SlugExists checks if a slug already exists
func (s *ProblemService) SlugExists(slug string, excludeID int) (bool, error) {
	return s.problemRepo.SlugExists(slug, excludeID)
}
