package validation

import (
	"context"
	"errors"
	"fmt"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/infrastructure/queue"
	"github.com/prabalesh/loco/backend/internal/services/execution"
)

type ValidationService struct {
	referenceSolutionRepo domain.ReferenceSolutionRepository
	problemRepo           domain.ProblemRepository
	testCaseRepo          domain.TestCaseRepository
	submissionRepo        domain.SubmissionRepository
	jobQueue              queue.JobQueue
	executionService      *execution.ExecutionService
}

func NewValidationService(
	referenceSolutionRepo domain.ReferenceSolutionRepository,
	problemRepo domain.ProblemRepository,
	testCaseRepo domain.TestCaseRepository,
	submissionRepo domain.SubmissionRepository,
	jobQueue queue.JobQueue,
	executionService *execution.ExecutionService,
) *ValidationService {
	return &ValidationService{
		referenceSolutionRepo: referenceSolutionRepo,
		problemRepo:           problemRepo,
		testCaseRepo:          testCaseRepo,
		submissionRepo:        submissionRepo,
		jobQueue:              jobQueue,
		executionService:      executionService,
	}
}

type ValidateRequest struct {
	ProblemID    int    `json:"problem_id"`
	LanguageSlug string `json:"language_slug"`
	Code         string `json:"code"`
}

type ValidationResult struct {
	IsValid      bool                    `json:"is_valid"`
	PassedTests  int                     `json:"passed_tests"`
	TotalTests   int                     `json:"total_tests"`
	TestResults  []domain.TestCaseResult `json:"test_results"`
	ErrorMessage string                  `json:"error_message,omitempty"`
}

// ValidateReferenceSolution validates a reference solution against all test cases
func (s *ValidationService) ValidateReferenceSolution(req ValidateRequest, languageID int) (*ValidationResult, error) {
	// Get problem
	_, err := s.problemRepo.GetByID(req.ProblemID)
	if err != nil {
		return nil, fmt.Errorf("problem not found: %w", err)
	}

	// Get all test cases
	testCases, err := s.testCaseRepo.GetByProblemID(req.ProblemID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch test cases: %w", err)
	}

	if len(testCases) == 0 {
		return nil, errors.New("no test cases found for this problem")
	}

	// Execute reference solution
	execReq := execution.ExecutionRequest{
		ProblemID:  req.ProblemID,
		LanguageID: languageID,
		UserCode:   req.Code,
		TestCases:  testCases,
	}

	execResult, err := s.executionService.ExecuteSubmission(execReq, req.LanguageSlug)
	if err != nil {
		return &ValidationResult{
			IsValid:      false,
			ErrorMessage: fmt.Sprintf("Execution failed: %v", err),
		}, nil
	}

	// Check if all tests passed
	isValid := execResult.Status == domain.SubmissionStatusAccepted

	result := &ValidationResult{
		IsValid:      isValid,
		PassedTests:  execResult.PassedTests,
		TotalTests:   execResult.TotalTests,
		TestResults:  execResult.TestResults,
		ErrorMessage: execResult.ErrorMessage,
	}

	return result, nil
}

// SaveReferenceSolution saves and validates a reference solution asynchronously
func (s *ValidationService) SaveReferenceSolution(req ValidateRequest, languageID int, adminID int) (*domain.ProblemReferenceSolution, *domain.Submission, error) {
	// 1. Check if problem exists
	_, err := s.problemRepo.GetByID(req.ProblemID)
	if err != nil {
		return nil, nil, fmt.Errorf("problem not found: %w", err)
	}

	// 2. Create or Update ReferenceSolution (Pending state)
	exists, err := s.referenceSolutionRepo.Exists(req.ProblemID, languageID)
	if err != nil {
		return nil, nil, err
	}

	var referenceSolution *domain.ProblemReferenceSolution

	if exists {
		// Update existing
		referenceSolution, err = s.referenceSolutionRepo.GetByProblemAndLanguage(req.ProblemID, languageID)
		if err != nil {
			return nil, nil, err
		}
		referenceSolution.Code = req.Code
		referenceSolution.IsValidated = false // Mark as not validated until worker finishes
		// We don't clear validation results yet, or maybe we should?
		// Let's keep old results until new ones arrive, but IsValidated=false indicates it's stale/pending.

		if err := s.referenceSolutionRepo.Update(referenceSolution); err != nil {
			return nil, nil, err
		}
	} else {
		// Create new
		referenceSolution = &domain.ProblemReferenceSolution{
			ProblemID:   req.ProblemID,
			LanguageID:  languageID,
			Code:        req.Code,
			IsValidated: false,
		}

		if err := s.referenceSolutionRepo.Create(referenceSolution); err != nil {
			return nil, nil, err
		}
	}

	// 3. Create Validation Submission
	submission := &domain.Submission{
		UserID:                 adminID,
		ProblemID:              req.ProblemID,
		LanguageID:             languageID,
		Code:                   req.Code,
		Status:                 domain.SubmissionStatusPending,
		IsAdminSubmission:      true,
		IsValidationSubmission: true,
		SubmittedBy:            &adminID,
	}

	if err := s.submissionRepo.Create(submission); err != nil {
		return nil, nil, fmt.Errorf("failed to create submission: %w", err)
	}

	// 4. Enqueue Submission
	ctx := context.Background()
	if err := s.jobQueue.EnqueueSubmission(ctx, submission.ID); err != nil {
		// If queue fails, mark submission as error
		submission.Status = domain.SubmissionStatusInternalError
		submission.ErrorMessage = "Failed to enqueue validation job"
		s.submissionRepo.Update(submission)
		return nil, nil, fmt.Errorf("failed to enqueue submission: %w", err)
	}

	return referenceSolution, submission, nil
}

// GetValidationStatus returns validation status for a problem
func (s *ValidationService) GetValidationStatus(problemID int) (map[string]interface{}, error) {
	problem, err := s.problemRepo.GetByID(problemID)
	if err != nil {
		return nil, err
	}

	referenceSolutions, err := s.referenceSolutionRepo.GetAllByProblemID(problemID)
	if err != nil {
		return nil, err
	}

	validatedLanguages := []string{}
	for _, sol := range referenceSolutions {
		if sol.IsValidated && sol.Language.Name != "" {
			validatedLanguages = append(validatedLanguages, sol.Language.Name)
		}
	}

	status := map[string]interface{}{
		"problem_id":          problemID,
		"validation_status":   problem.ValidationStatus,
		"has_reference":       problem.HasReferenceSolution,
		"validated_languages": validatedLanguages,
		"total_solutions":     len(referenceSolutions),
		"can_publish":         problem.ValidationStatus == "validated",
	}

	return status, nil
}

// CanPublishProblem checks if problem can be published
func (s *ValidationService) CanPublishProblem(problemID int) (bool, error) {
	problem, err := s.problemRepo.GetByID(problemID)
	if err != nil {
		return false, err
	}

	return problem.ValidationStatus == "validated" && problem.HasReferenceSolution, nil
}
