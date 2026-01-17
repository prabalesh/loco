package usecase

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/domain/uerror"
	"github.com/prabalesh/loco/backend/internal/domain/validator" // This import is kept because it's used later in the file.
	"github.com/prabalesh/loco/backend/pkg/config"
	"github.com/prabalesh/loco/backend/pkg/utils"
	"go.uber.org/zap"
)

type ProblemUsecase struct {
	problemRepo  domain.ProblemRepository
	testcaseRepo domain.TestCaseRepository
	cfg          *config.Config
	logger       *zap.Logger
}

func NewProblemUsecase(problemRepo domain.ProblemRepository, testcaseRepo domain.TestCaseRepository, cfg *config.Config, logger *zap.Logger) *ProblemUsecase {
	return &ProblemUsecase{
		problemRepo:  problemRepo,
		testcaseRepo: testcaseRepo,
		cfg:          cfg,
		logger:       logger,
	}
}

// ========== ADMIN OPERATIONS ==========

// CreateProblem creates a new problem (draft by default)
func (u *ProblemUsecase) CreateProblem(req *domain.CreateProblemRequest, adminID int) (*domain.Problem, error) {
	// Validation
	if validationErrors := validator.ValidateCreateProblemRequest(req); len(validationErrors) > 0 {
		u.logger.Warn("Create problem validation failed",
			zap.Any("errors", validationErrors),
		)
		return nil, &uerror.ValidationError{Errors: validationErrors}
	}

	// Generate slug from title if not provided
	slug := req.Slug
	if slug == "" {
		slug = generateSlug(req.Title)
	}

	// Check if slug already exists
	exists, err := u.problemRepo.SlugExists(slug)
	if err != nil {
		u.logger.Error("Failed to check slug existence",
			zap.Error(err),
			zap.String("slug", slug),
		)
		return nil, errors.New("failed to create problem")
	}

	if exists {
		u.logger.Warn("Problem creation failed: slug already exists",
			zap.String("slug", slug),
		)
		return nil, errors.New("problem with similar title already exists")
	}

	// Set defaults
	timeLimit := req.TimeLimit
	if timeLimit == 0 {
		timeLimit = 2000 // Default 2 seconds
	}

	memoryLimit := req.MemoryLimit
	if memoryLimit == 0 {
		memoryLimit = 256 // Default 256 MB
	}

	validatorType := req.ValidatorType
	if validatorType == "" {
		validatorType = "exact_match"
	}

	status := req.Status
	if status == "" {
		status = "draft"
	}

	visibility := req.Visibility
	if visibility == "" {
		visibility = "public"
	}

	problem := &domain.Problem{
		Title:         req.Title,
		Slug:          slug,
		Description:   req.Description,
		Difficulty:    req.Difficulty,
		TimeLimit:     timeLimit,
		MemoryLimit:   memoryLimit,
		ValidatorType: validatorType,
		InputFormat:   req.InputFormat,
		OutputFormat:  req.OutputFormat,
		Constraints:   req.Constraints,
		Status:        status,
		Visibility:    visibility,
		IsActive:      req.IsActive,
		CreatedBy:     &adminID,
	}

	if err := u.problemRepo.Create(problem); err != nil {
		u.logger.Error("Failed to create problem in database",
			zap.Error(err),
			zap.String("title", req.Title),
			zap.Int("admin_id", adminID),
		)
		return nil, errors.New("failed to create problem")
	}

	u.logger.Info("Problem created successfully",
		zap.Int("problem_id", problem.ID),
		zap.String("title", problem.Title),
		zap.String("slug", problem.Slug),
		zap.Int("created_by", adminID),
	)

	return problem, nil
}

// UpdateProblem updates an existing problem
func (u *ProblemUsecase) UpdateProblem(problemID int, req *domain.UpdateProblemRequest, adminID int) (*domain.Problem, error) {
	// Validation
	if validationErrors := validator.ValidateUpdateProblemRequest(req); len(validationErrors) > 0 {
		u.logger.Warn("Update problem validation failed",
			zap.Any("errors", validationErrors),
		)
		return nil, &uerror.ValidationError{Errors: validationErrors}
	}

	// Get existing problem
	problem, err := u.problemRepo.GetByID(problemID)
	if err != nil {
		u.logger.Warn("Problem not found for update",
			zap.Int("problem_id", problemID),
		)
		return nil, errors.New("problem not found")
	}

	// Update fields
	if req.Title != "" {
		problem.Title = req.Title
	}

	if req.Slug != "" {
		problem.Slug = req.Slug
	}

	if req.Description != "" {
		problem.Description = req.Description
	}

	if req.Difficulty != "" {
		problem.Difficulty = req.Difficulty
	}

	if req.TimeLimit > 0 {
		problem.TimeLimit = req.TimeLimit
	}

	if req.MemoryLimit > 0 {
		problem.MemoryLimit = req.MemoryLimit
	}

	if req.ValidatorType != "" {
		problem.ValidatorType = req.ValidatorType
	}

	if req.InputFormat != "" {
		problem.InputFormat = req.InputFormat
	}

	if req.OutputFormat != "" {
		problem.OutputFormat = req.OutputFormat
	}

	if req.Constraints != "" {
		problem.Constraints = req.Constraints
	}

	if req.Status != "" {
		problem.Status = req.Status
	}

	if req.Visibility != "" {
		problem.Visibility = req.Visibility
	}

	if req.IsActive != nil {
		problem.IsActive = *req.IsActive
	}

	if err := u.problemRepo.Update(problem); err != nil {
		u.logger.Error("Failed to update problem",
			zap.Error(err),
			zap.Int("problem_id", problemID),
		)
		return nil, errors.New("failed to update problem")
	}

	u.logger.Info("Problem updated successfully",
		zap.Int("problem_id", problem.ID),
		zap.String("title", problem.Title),
		zap.Int("updated_by", adminID),
	)

	return problem, nil
}

func (u *ProblemUsecase) ValidateTestCases(problemID int, adminID int) error {
	count, err := u.testcaseRepo.CountByProblemID(problemID)
	if err != nil {
		return errors.New("failed to count test cases")
	}

	if count < 2 {
		return errors.New("at least 2 test cases are required")
	}

	// Update problem's current_step to 2
	err = u.problemRepo.UpdateCurrentStep(problemID, 2)
	if err != nil {
		return errors.New("failed to update problem step")
	}

	u.logger.Info("Problem test cases validated and step updated",
		zap.Int("problem_id", problemID),
		zap.Int("admin_id", adminID),
	)

	return nil
}

// DeleteProblem deletes a problem
func (u *ProblemUsecase) DeleteProblem(problemID int, adminID int) error {
	// Check if problem exists
	_, err := u.problemRepo.GetByID(problemID)
	if err != nil {
		u.logger.Warn("Problem not found for deletion",
			zap.Int("problem_id", problemID),
		)
		return errors.New("problem not found")
	}

	if err := u.problemRepo.Delete(problemID); err != nil {
		u.logger.Error("Failed to delete problem",
			zap.Error(err),
			zap.Int("problem_id", problemID),
		)
		return errors.New("failed to delete problem")
	}

	u.logger.Info("Problem deleted successfully",
		zap.Int("problem_id", problemID),
		zap.Int("deleted_by", adminID),
	)

	return nil
}

// PublishProblem changes status from draft to published
func (u *ProblemUsecase) PublishProblem(problemID int, adminID int) error {
	problem, err := u.problemRepo.GetByID(problemID)
	if err != nil {
		return errors.New("problem not found")
	}

	if problem.Status == "published" {
		return errors.New("problem is already published")
	}

	// Validation: Check if problem is ready to publish
	// TODO: Add validation for examples, test cases, languages, etc.

	if err := u.problemRepo.UpdateStatus(problemID, "published"); err != nil {
		u.logger.Error("Failed to publish problem",
			zap.Error(err),
			zap.Int("problem_id", problemID),
		)
		return errors.New("failed to publish problem")
	}

	u.logger.Info("Problem published successfully",
		zap.Int("problem_id", problemID),
		zap.Int("published_by", adminID),
	)

	return nil
}

// ArchiveProblem changes status to archived
func (u *ProblemUsecase) ArchiveProblem(problemID int, adminID int) error {
	if err := u.problemRepo.UpdateStatus(problemID, "archived"); err != nil {
		u.logger.Error("Failed to archive problem",
			zap.Error(err),
			zap.Int("problem_id", problemID),
		)
		return errors.New("failed to archive problem")
	}

	u.logger.Info("Problem archived successfully",
		zap.Int("problem_id", problemID),
		zap.Int("archived_by", adminID),
	)

	return nil
}

// ========== USER OPERATIONS ==========

// GetProblem retrieves a single problem by ID or slug
func (u *ProblemUsecase) GetProblem(identifier string) (*domain.Problem, error) {
	var problem *domain.Problem
	var err error

	// Try to get by ID first, then by slug
	if id, parseErr := utils.ParseInt(identifier); parseErr == nil {
		problem, err = u.problemRepo.GetByID(id)
	} else {
		problem, err = u.problemRepo.GetBySlug(identifier)
	}

	if err != nil {
		u.logger.Warn("Problem not found",
			zap.String("identifier", identifier),
		)
		return nil, errors.New("problem not found")
	}

	return problem, nil
}

// ListProblems retrieves problems with filters (for users - only published & public)
func (u *ProblemUsecase) ListProblems(req *domain.ListProblemsRequest) ([]*domain.Problem, int, error) {
	filters := domain.ProblemFilters{
		Page:       req.Page,
		Limit:      req.Limit,
		Difficulty: req.Difficulty,
		Status:     "published", // Only published problems for users
		Visibility: "public",    // Only public problems
		Search:     req.Search,
		Tags:       req.Tags,
	}

	problems, total, err := u.problemRepo.List(filters)
	if err != nil {
		u.logger.Error("Failed to list problems",
			zap.Error(err),
		)
		return nil, 0, errors.New("failed to retrieve problems")
	}

	return problems, total, nil
}

// ========== ADMIN LIST OPERATIONS ==========

// ListAllProblems retrieves all problems (admin - includes drafts, private)
func (u *ProblemUsecase) ListAllProblems(req *domain.AdminListProblemsRequest) ([]*domain.Problem, int, error) {
	filters := domain.ProblemFilters{
		Page:       req.Page,
		Limit:      req.Limit,
		Difficulty: req.Difficulty,
		Status:     req.Status,
		Visibility: req.Visibility,
		Search:     req.Search,
		Tags:       req.Tags,
	}

	problems, total, err := u.problemRepo.List(filters)
	if err != nil {
		u.logger.Error("Failed to list all problems (admin)",
			zap.Error(err),
		)
		return nil, 0, errors.New("failed to retrieve problems")
	}

	return problems, total, nil
}

// GetProblemStats returns problem statistics
func (u *ProblemUsecase) GetProblemStats() (*domain.ProblemStats, error) {
	totalProblems, err := u.problemRepo.CountProblems()
	if err != nil {
		return nil, err
	}

	publishedCount, err := u.problemRepo.CountByStatus("published")
	if err != nil {
		return nil, err
	}

	draftCount, err := u.problemRepo.CountByStatus("draft")
	if err != nil {
		return nil, err
	}

	easyCount, err := u.problemRepo.CountByDifficulty("easy")
	if err != nil {
		return nil, err
	}

	mediumCount, err := u.problemRepo.CountByDifficulty("medium")
	if err != nil {
		return nil, err
	}

	hardCount, err := u.problemRepo.CountByDifficulty("hard")
	if err != nil {
		return nil, err
	}

	return &domain.ProblemStats{
		Total:     totalProblems,
		Published: publishedCount,
		Draft:     draftCount,
		Easy:      easyCount,
		Medium:    mediumCount,
		Hard:      hardCount,
	}, nil
}

// ========== HELPER FUNCTIONS ==========

// generateSlug creates URL-friendly slug from title
func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove special characters except hyphens
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)
	// Remove consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	// Trim hyphens from start/end
	slug = strings.Trim(slug, "-")

	// Add timestamp suffix for uniqueness
	slug = fmt.Sprintf("%s-%d", slug, time.Now().Unix())

	return slug
}
