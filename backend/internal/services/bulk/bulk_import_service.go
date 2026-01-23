package bulk

import (
	"fmt"
	"strings"
	"time"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/services/problem"
	"github.com/prabalesh/loco/backend/internal/services/validation"
	"gorm.io/gorm"
)

type BulkImportService struct {
	problemService    *problem.ProblemService
	validationService *validation.ValidationService
	db                *gorm.DB
}

func NewBulkImportService(problemService *problem.ProblemService, validationService *validation.ValidationService, db *gorm.DB) *BulkImportService {
	return &BulkImportService{
		problemService:    problemService,
		validationService: validationService,
		db:                db,
	}
}

type BulkImportRequest struct {
	Problems []ProblemImportData `json:"problems"`
	Options  ImportOptions       `json:"options"`
}

type ProblemImportData struct {
	Title                   string                   `json:"title"`
	Description             string                   `json:"description"`
	Difficulty              string                   `json:"difficulty"`
	CategoryIDs             []int                    `json:"category_ids"`
	TagIDs                  []int                    `json:"tag_ids"`
	FunctionName            string                   `json:"function_name"`
	ReturnType              domain.GenericType       `json:"return_type"`
	Parameters              []domain.SchemaParameter `json:"parameters"`
	ValidationType          string                   `json:"validation_type"`
	ExpectedTimeComplexity  string                   `json:"expected_time_complexity"`
	ExpectedSpaceComplexity string                   `json:"expected_space_complexity"`
	TestCases               []problem.TestCaseInput  `json:"test_cases"`
	ReferenceSolution       *ReferenceSolutionData   `json:"reference_solution,omitempty"`
}

type ReferenceSolutionData struct {
	LanguageSlug string `json:"language_slug"`
	Code         string `json:"code"`
}

type ImportOptions struct {
	ValidateReferences bool `json:"validate_references"` // Auto-validate reference solutions
	SkipDuplicates     bool `json:"skip_duplicates"`     // Skip if slug exists
	StopOnError        bool `json:"stop_on_error"`       // Stop entire import on first error
}

type BulkImportResult struct {
	TotalSubmitted   int                    `json:"total_submitted"`
	TotalCreated     int                    `json:"total_created"`
	TotalFailed      int                    `json:"total_failed"`
	CreatedProblems  []ProblemImportSuccess `json:"created_problems"`
	FailedProblems   []ProblemImportFailure `json:"failed_problems"`
	ProcessingTimeMs int64                  `json:"processing_time_ms"`
}

type ProblemImportSuccess struct {
	Index            int    `json:"index"`
	Title            string `json:"title"`
	Slug             string `json:"slug"`
	ProblemID        int    `json:"problem_id"`
	ValidationStatus string `json:"validation_status"`
}

type ProblemImportFailure struct {
	Index        int      `json:"index"`
	Title        string   `json:"title"`
	Errors       []string `json:"errors"`
	ErrorMessage string   `json:"error_message"`
}

// BulkImport imports multiple problems
func (s *BulkImportService) BulkImport(req BulkImportRequest, createdBy int) (*BulkImportResult, error) {
	startTime := time.Now()

	result := &BulkImportResult{
		TotalSubmitted:  len(req.Problems),
		CreatedProblems: []ProblemImportSuccess{},
		FailedProblems:  []ProblemImportFailure{},
	}

	// Validate all problems first (basic validation)
	validationErrors := s.validateAllProblems(req.Problems)

	// Process each problem
	for i, problemData := range req.Problems {
		// Check if validation failed
		if errors, hasError := validationErrors[i]; hasError {
			result.FailedProblems = append(result.FailedProblems, ProblemImportFailure{
				Index:        i,
				Title:        problemData.Title,
				Errors:       errors,
				ErrorMessage: strings.Join(errors, "; "),
			})
			result.TotalFailed++

			if req.Options.StopOnError {
				break
			}
			continue
		}

		// Check for duplicates
		if req.Options.SkipDuplicates {
			slug := s.problemService.GenerateSlug(problemData.Title)
			exists, _ := s.problemService.SlugExists(slug, 0)
			if exists {
				result.FailedProblems = append(result.FailedProblems, ProblemImportFailure{
					Index:        i,
					Title:        problemData.Title,
					ErrorMessage: "Problem with this title already exists (slug conflict)",
				})
				result.TotalFailed++
				continue
			}
		}

		// Create problem
		createReq := s.convertToCreateRequest(problemData)
		createdProblem, err := s.problemService.CreateProblem(createReq, createdBy)
		if err != nil {
			result.FailedProblems = append(result.FailedProblems, ProblemImportFailure{
				Index:        i,
				Title:        problemData.Title,
				ErrorMessage: err.Error(),
			})
			result.TotalFailed++

			if req.Options.StopOnError {
				break
			}
			continue
		}

		// Validate reference solution if provided
		validationStatus := "draft"
		if problemData.ReferenceSolution != nil && req.Options.ValidateReferences {
			validationStatus = s.validateReferenceSolution(createdProblem.ID, *problemData.ReferenceSolution, createdBy)
		}

		result.CreatedProblems = append(result.CreatedProblems, ProblemImportSuccess{
			Index:            i,
			Title:            createdProblem.Title,
			Slug:             createdProblem.Slug,
			ProblemID:        createdProblem.ID,
			ValidationStatus: validationStatus,
		})
		result.TotalCreated++
	}

	result.ProcessingTimeMs = time.Since(startTime).Milliseconds()

	return result, nil
}

// validateAllProblems validates all problems and returns errors by index
func (s *BulkImportService) validateAllProblems(problems []ProblemImportData) map[int][]string {
	errors := make(map[int][]string)

	for i, p := range problems {
		problemErrors := []string{}

		// Validate required fields
		if len(p.Title) < 5 || len(p.Title) > 200 {
			problemErrors = append(problemErrors, "title must be 5-200 characters")
		}
		if len(p.Description) < 20 {
			problemErrors = append(problemErrors, "description must be at least 20 characters")
		}
		if p.Difficulty != "easy" && p.Difficulty != "medium" && p.Difficulty != "hard" {
			problemErrors = append(problemErrors, "difficulty must be easy, medium, or hard")
		}
		if p.FunctionName == "" {
			problemErrors = append(problemErrors, "function_name is required")
		}
		if p.ReturnType == "" {
			problemErrors = append(problemErrors, "return_type is required")
		}
		if len(p.Parameters) == 0 {
			problemErrors = append(problemErrors, "at least one parameter is required")
		}
		if len(p.TestCases) == 0 {
			problemErrors = append(problemErrors, "at least one test case is required")
		}

		// Check for at least one public test case
		hasPublic := false
		for _, tc := range p.TestCases {
			if tc.IsSample {
				hasPublic = true
				break
			}
		}
		if !hasPublic {
			problemErrors = append(problemErrors, "at least one public test case required")
		}

		if len(problemErrors) > 0 {
			errors[i] = problemErrors
		}
	}

	return errors
}

// convertToCreateRequest converts import data to create request
func (s *BulkImportService) convertToCreateRequest(data ProblemImportData) problem.CreateProblemRequest {
	return problem.CreateProblemRequest{
		Title:                   data.Title,
		Description:             data.Description,
		Difficulty:              data.Difficulty,
		CategoryIDs:             data.CategoryIDs,
		TagIDs:                  data.TagIDs,
		FunctionName:            data.FunctionName,
		ReturnType:              data.ReturnType,
		Parameters:              data.Parameters,
		ValidationType:          data.ValidationType,
		ExpectedTimeComplexity:  data.ExpectedTimeComplexity,
		ExpectedSpaceComplexity: data.ExpectedSpaceComplexity,
		TestCases:               data.TestCases,
	}
}

// validateReferenceSolution validates reference solution if provided
func (s *BulkImportService) validateReferenceSolution(problemID int, refSol ReferenceSolutionData, adminID int) string {
	// Get language ID
	var language domain.Language
	if err := s.db.Where("slug = ?", refSol.LanguageSlug).First(&language).Error; err != nil {
		return "draft" // Language not found, skip validation
	}

	// Validate
	validateReq := validation.ValidateRequest{
		ProblemID:    problemID,
		LanguageSlug: refSol.LanguageSlug,
		Code:         refSol.Code,
	}

	_, _, err := s.validationService.SaveReferenceSolution(validateReq, language.ID, adminID)
	if err != nil {
		return "draft"
	}

	// Return "pending" since validation is asynchronous
	return "pending"
}

// BulkImportAsync processes import asynchronously (for large batches)
func (s *BulkImportService) BulkImportAsync(req BulkImportRequest, createdBy int) (string, error) {
	// Generate job ID
	jobID := fmt.Sprintf("import_%d_%d", createdBy, time.Now().Unix())

	// Process in goroutine
	go func() {
		result, err := s.BulkImport(req, createdBy)
		if err != nil {
			// Log error
			fmt.Printf("Bulk import job %s failed: %v\n", jobID, err)
		} else {
			// Store result (could use Redis or database)
			fmt.Printf("Bulk import job %s completed: %d created, %d failed\n",
				jobID, result.TotalCreated, result.TotalFailed)
		}
	}()

	return jobID, nil
}
