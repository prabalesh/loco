package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/infrastructure/piston"
	"github.com/prabalesh/loco/backend/internal/infrastructure/queue"
	"github.com/prabalesh/loco/backend/pkg/config"
	"go.uber.org/zap"
)

type SubmissionUsecase struct {
	submissionRepo      domain.SubmissionRepository
	problemRepo         domain.ProblemRepository
	testCaseRepo        domain.TestCaseRepository
	languageRepo        domain.LanguageRepository
	problemLanguageRepo domain.ProblemLanguageRepository
	pistonService       piston.PistonService
	jobQueue            queue.JobQueue
	achievementUsecase  *AchievementUsecase
	cfg                 *config.Config
	logger              *zap.Logger
}

func NewSubmissionUsecase(
	submissionRepo domain.SubmissionRepository,
	problemRepo domain.ProblemRepository,
	testCaseRepo domain.TestCaseRepository,
	languageRepo domain.LanguageRepository,
	problemLanguageRepo domain.ProblemLanguageRepository,
	pistonService piston.PistonService,
	jobQueue queue.JobQueue,
	achievementUsecase *AchievementUsecase,
	cfg *config.Config,
	logger *zap.Logger,
) *SubmissionUsecase {
	return &SubmissionUsecase{
		submissionRepo:      submissionRepo,
		problemRepo:         problemRepo,
		testCaseRepo:        testCaseRepo,
		languageRepo:        languageRepo,
		problemLanguageRepo: problemLanguageRepo,
		pistonService:       pistonService,
		jobQueue:            jobQueue,
		achievementUsecase:  achievementUsecase,
		cfg:                 cfg,
		logger:              logger,
	}
}

func (u *SubmissionUsecase) Submit(userID int, problemID int, req *domain.CreateSubmissionRequest) (*domain.Submission, error) {
	// 1. Validate Problem and Language
	_, err := u.problemRepo.GetByID(problemID)
	if err != nil {
		return nil, fmt.Errorf("problem not found")
	}

	language, err := u.languageRepo.GetByID(req.LanguageID)
	if err != nil {
		return nil, fmt.Errorf("language not found")
	}

	// 2. Get ProblemLanguage to combine code
	pl, err := u.problemLanguageRepo.GetByProblemAndLanguage(problemID, req.LanguageID)
	finalCode := req.Code
	if err == nil && pl != nil {
		finalCode = pl.GetCombinedCode(language.DefaultTemplate, req.Code)
	}

	// 3. Create Pending Submission
	now := time.Now()
	submission := &domain.Submission{
		UserID:     userID,
		ProblemID:  problemID,
		LanguageID: req.LanguageID,
		Code:       finalCode,
		Status:     domain.SubmissionStatusPending,
		QueuedAt:   &now,
	}

	if err := u.submissionRepo.Create(submission); err != nil {
		return nil, fmt.Errorf("failed to create submission: %w", err)
	}

	// 4. Enqueue submission job to Redis queue
	ctx := context.Background()
	if err := u.jobQueue.EnqueueSubmission(ctx, submission.ID); err != nil {
		u.logger.Error("Failed to enqueue submission",
			zap.Error(err),
			zap.Int("submission_id", submission.ID),
		)
		// Update submission status to indicate queue failure
		submission.Status = domain.SubmissionStatusInternalError
		submission.ErrorMessage = "Failed to enqueue submission for processing"
		u.submissionRepo.Update(submission)
		return nil, fmt.Errorf("failed to enqueue submission: %w", err)
	}

	u.logger.Info("Submission enqueued successfully",
		zap.Int("submission_id", submission.ID),
		zap.Int("user_id", userID),
		zap.Int("problem_id", problemID),
	)

	return submission, nil
}

// ProcessSubmission handles the background processing of a submission
func (u *SubmissionUsecase) ProcessSubmission(submissionID int) error {
	u.logger.Info("Processing submission", zap.Int("submission_id", submissionID))

	// 1. Fetch submission
	submission, err := u.submissionRepo.GetByID(submissionID)
	if err != nil {
		return fmt.Errorf("failed to fetch submission: %w", err)
	}

	if submission.Status != domain.SubmissionStatusPending {
		u.logger.Info("Submission already processed or not pending",
			zap.Int("submission_id", submissionID),
			zap.String("status", string(submission.Status)),
		)
		return nil
	}

	// 2. Fetch related data
	problem, err := u.problemRepo.GetByID(submission.ProblemID)
	if err != nil {
		u.updateSubmissionResult(submission, domain.SubmissionStatusInternalError, "Problem not found")
		return fmt.Errorf("problem not found: %w", err)
	}

	language, err := u.languageRepo.GetByID(submission.LanguageID)
	if err != nil {
		u.updateSubmissionResult(submission, domain.SubmissionStatusInternalError, "Language not found")
		return fmt.Errorf("language not found: %w", err)
	}

	// 3. Update status to Processing
	u.updateSubmissionResult(submission, domain.SubmissionStatusProcessing, "")

	// 4. Evaluate
	// Note: evaluateSubmission handles the logic and status updates
	u.evaluateSubmission(submission, problem, language, false)

	return nil
}

func (u *SubmissionUsecase) evaluateSubmission(submission *domain.Submission, problem *domain.Problem, language *domain.Language, runOnlyPublicTests bool) {
	var testCases []domain.TestCase
	var err error

	if runOnlyPublicTests {
		testCases, err = u.testCaseRepo.GetSamples(submission.ProblemID)
	} else {
		testCases, err = u.testCaseRepo.GetByProblemID(submission.ProblemID)
	}
	if err != nil {
		u.logger.Error("Failed to fetch test cases", zap.Error(err))
		u.updateSubmissionResult(submission, domain.SubmissionStatusInternalError, "Failed to fetch test cases")
		return
	}

	// 1. Code is already combined in Submit/Validate
	finalCode := submission.Code

	finalStatus := domain.SubmissionStatusAccepted
	errorMessage := ""
	passCount := 0
	totalCount := len(testCases)

	var results []domain.TestCaseResult

	submission.TotalTestCases = totalCount

	for _, tc := range testCases {
		result, err := u.pistonService.Execute(language.Slug, language.Version, finalCode, tc.Input)

		tcResult := domain.TestCaseResult{
			Input:          tc.Input,
			ExpectedOutput: tc.ExpectedOutput,
			IsSample:       tc.IsSample,
		}

		if err != nil {
			u.logger.Error("Piston execution failed", zap.Error(err))
			finalStatus = domain.SubmissionStatusInternalError
			errorMessage = "Execution system error"
			tcResult.Status = "Failed"
			tcResult.ActualOutput = "System Error"
			results = append(results, tcResult)
			break
		}

		if result.ExitCode != 0 {
			if finalStatus == domain.SubmissionStatusAccepted {
				finalStatus = domain.SubmissionStatusRuntimeError
				errorMessage = result.Error
			}
			tcResult.Status = "Failed"
			tcResult.ActualOutput = result.Error // Or specific error message
		} else {
			// Normalize output (trim whitespace)
			actual := strings.TrimSpace(result.Output)
			expected := strings.TrimSpace(tc.ExpectedOutput)

			tcResult.ActualOutput = actual

			if actual != expected {
				if finalStatus == domain.SubmissionStatusAccepted {
					finalStatus = domain.SubmissionStatusWrongAnswer
					errorMessage = fmt.Sprintf("Failed on input: %s\nExpected: %s\nActual: %s", tc.Input, expected, actual)
				}
				tcResult.Status = "Failed"
			} else {
				tcResult.Status = "Passed"
				passCount++
			}
		}

		results = append(results, tcResult)
	}

	submission.PassedTestCases = passCount
	submission.TestCaseResults = results

	// Update submission status
	u.updateSubmissionResult(submission, finalStatus, errorMessage)

	// If it's a validation submission, update ProblemLanguage status and last error
	if submission.IsValidationSubmission {
		now := time.Now()
		// Re-fetch pl just in case
		pl, err := u.problemLanguageRepo.GetByProblemAndLanguage(submission.ProblemID, submission.LanguageID)
		if err == nil {
			pl.LastValidationStatus = string(finalStatus)
			pl.LastValidationError = errorMessage
			pl.LastPassCount = passCount
			pl.LastTotalCount = totalCount

			if finalStatus == domain.SubmissionStatusAccepted {
				pl.IsValidated = true
				pl.ValidatedAt = &now
			} else {
				pl.IsValidated = false
				pl.ValidatedAt = &now
			}
			u.problemLanguageRepo.Update(pl)
			u.problemLanguageRepo.Update(pl)
		}
	}

	// Note: Achievements are now handled by AchievementWorker via AchievementQueue
}

func (u *SubmissionUsecase) updateSubmissionResult(submission *domain.Submission, status domain.SubmissionStatus, errorMsg string) {
	submission.Status = status
	submission.ErrorMessage = errorMsg
	if err := u.submissionRepo.Update(submission); err != nil {
		u.logger.Error("Failed to update submission status", zap.Error(err))
	}
}

func (u *SubmissionUsecase) GetSubmission(id int) (*domain.Submission, error) {
	return u.submissionRepo.GetByID(id)
}

func (u *SubmissionUsecase) GetUserProblemSubmissions(userID int, problemID int, limit, offset int) ([]domain.Submission, int64, error) {
	submissions, err := u.submissionRepo.ListByUserProblem(userID, problemID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.submissionRepo.CountByUserProblem(userID, problemID)
	if err != nil {
		return nil, 0, err
	}
	return submissions, count, nil
}

func (u *SubmissionUsecase) GetUserSubmissions(userID int, limit, offset int) ([]domain.Submission, int64, error) {
	submissions, err := u.submissionRepo.ListByUser(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.submissionRepo.CountByUser(userID)
	if err != nil {
		return nil, 0, err
	}
	return submissions, count, nil
}

func (u *SubmissionUsecase) GetAdminUserSubmissions(userID int, limit, offset int) ([]domain.Submission, int64, error) {
	submissions, err := u.submissionRepo.ListByAdminUser(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.submissionRepo.CountByUser(userID)
	if err != nil {
		return nil, 0, err
	}
	return submissions, count, nil
}

// AdminSubmit handles admin test submissions with admin context
func (u *SubmissionUsecase) AdminSubmit(adminID int, req *domain.CreateSubmissionRequest) (*domain.Submission, error) {
	// 1. Validate Problem and Language
	problem, err := u.problemRepo.GetByID(req.ProblemID)
	if err != nil {
		return nil, fmt.Errorf("problem not found")
	}

	language, err := u.languageRepo.GetByID(req.LanguageID)
	if err != nil {
		return nil, fmt.Errorf("language not found")
	}

	// 2. Get ProblemLanguage to combine code
	pl, err := u.problemLanguageRepo.GetByProblemAndLanguage(req.ProblemID, req.LanguageID)
	finalCode := req.Code
	if err == nil && pl != nil {
		finalCode = pl.GetCombinedCode(language.DefaultTemplate, req.Code)
	}

	// 3. Create Pending Submission with admin context
	submission := &domain.Submission{
		UserID:            adminID, // Admin is the submitter
		ProblemID:         req.ProblemID,
		LanguageID:        req.LanguageID,
		Code:              finalCode,
		Status:            domain.SubmissionStatusPending,
		IsAdminSubmission: true,
		SubmittedBy:       &adminID,
	}

	if err := u.submissionRepo.Create(submission); err != nil {
		return nil, fmt.Errorf("failed to create submission: %w", err)
	}

	// 3. Execute evaluation asynchronously (always run all test cases for admin submissions)
	go u.evaluateSubmission(submission, problem, language, false)

	return submission, nil
}

// Validate handles solution code validation for a problem-language pair
func (u *SubmissionUsecase) Validate(adminID int, problemID int, languageID int, code string) (*domain.Submission, error) {
	// 1. Validate Problem and Language
	problem, err := u.problemRepo.GetByID(problemID)
	if err != nil {
		return nil, fmt.Errorf("problem not found")
	}

	language, err := u.languageRepo.GetByID(languageID)
	if err != nil {
		return nil, fmt.Errorf("language not found")
	}

	// 2. Get ProblemLanguage to combine code
	pl, err := u.problemLanguageRepo.GetByProblemAndLanguage(problemID, languageID)
	if err != nil {
		return nil, fmt.Errorf("problem language config not found")
	}

	// Use combined code for validation submission so it's stored in DB
	finalCode := pl.GetAdminCombinedCode(language.DefaultTemplate, code)

	// 3. Create Pending Submission with validation flag
	submission := &domain.Submission{
		UserID:                 adminID,
		ProblemID:              problemID,
		LanguageID:             languageID,
		Code:                   finalCode,
		Status:                 domain.SubmissionStatusPending,
		IsAdminSubmission:      true,
		IsValidationSubmission: true,
		SubmittedBy:            &adminID,
	}

	if err := u.submissionRepo.Create(submission); err != nil {
		return nil, fmt.Errorf("failed to create submission: %w", err)
	}

	// 4. Execute evaluation asynchronously (always run all test cases for validation)
	go u.evaluateSubmission(submission, problem, language, false)

	return submission, nil
}

// GetProblemSubmissions retrieves all submissions for a specific problem
func (u *SubmissionUsecase) GetProblemSubmissions(problemID int, limit, offset int) ([]domain.Submission, error) {
	return u.submissionRepo.ListByProblem(problemID, limit, offset)
}

// RunCode executes code against public test cases without creating a submission
func (u *SubmissionUsecase) RunCode(problemID int, req *domain.RunCodeRequest) (*domain.RunCodeResult, error) {
	// 1. Validate Problem and Language
	_, err := u.problemRepo.GetByID(problemID)
	if err != nil {
		return nil, fmt.Errorf("problem not found")
	}

	language, err := u.languageRepo.GetByID(req.LanguageID)
	if err != nil {
		return nil, fmt.Errorf("language not found")
	}

	// 2. Get ProblemLanguage to combine code
	pl, err := u.problemLanguageRepo.GetByProblemAndLanguage(problemID, req.LanguageID)
	finalCode := req.Code
	if err == nil && pl != nil {
		finalCode = pl.GetCombinedCode(language.DefaultTemplate, req.Code)
	}

	// 3. Get only public (sample) test cases
	testCases, err := u.testCaseRepo.GetSamples(problemID)
	if err != nil {
		u.logger.Error("Failed to fetch test cases", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch test cases: %w", err)
	}

	// 4. Execute test cases
	finalStatus := domain.SubmissionStatusAccepted
	errorMessage := ""
	passCount := 0
	totalCount := len(testCases)

	var results []domain.TestCaseResult

	for _, tc := range testCases {
		result, err := u.pistonService.Execute(language.Slug, language.Version, finalCode, tc.Input)

		tcResult := domain.TestCaseResult{
			Input:          tc.Input,
			ExpectedOutput: tc.ExpectedOutput,
			IsSample:       tc.IsSample,
		}

		if err != nil {
			u.logger.Error("Piston execution failed", zap.Error(err))
			finalStatus = domain.SubmissionStatusInternalError
			errorMessage = "Execution system error"
			tcResult.Status = "Failed"
			tcResult.ActualOutput = "System Error"
			results = append(results, tcResult)
			break
		}

		if result.ExitCode != 0 {
			if finalStatus == domain.SubmissionStatusAccepted {
				finalStatus = domain.SubmissionStatusRuntimeError
				errorMessage = result.Error
			}
			tcResult.Status = "Failed"
			tcResult.ActualOutput = result.Error
		} else {
			// Normalize output (trim whitespace)
			actual := strings.TrimSpace(result.Output)
			expected := strings.TrimSpace(tc.ExpectedOutput)

			tcResult.ActualOutput = actual

			if actual != expected {
				if finalStatus == domain.SubmissionStatusAccepted {
					finalStatus = domain.SubmissionStatusWrongAnswer
					errorMessage = fmt.Sprintf("Failed on input: %s\nExpected: %s\nActual: %s", tc.Input, expected, actual)
				}
				tcResult.Status = "Failed"
			} else {
				tcResult.Status = "Passed"
				passCount++
			}
		}

		results = append(results, tcResult)
	}

	return &domain.RunCodeResult{
		Status:          finalStatus,
		ErrorMessage:    errorMessage,
		PassedTestCases: passCount,
		TotalTestCases:  totalCount,
		Results:         results,
	}, nil
}
