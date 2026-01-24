package postgres

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/pkg/database"
)

type problemRepository struct {
	db *database.Database
}

func NewProblemRepository(db *database.Database) *problemRepository {
	return &problemRepository{db: db}
}

func (r *problemRepository) GetAll(limit, offset int, search string) ([]domain.Problem, int64, error) {
	var problems []domain.Problem
	var total int64

	query := r.db.DB.Model(&domain.Problem{})
	if search != "" {
		pattern := "%" + search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", pattern, pattern)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Creator").Preload("Tags").Preload("Categories").Limit(limit).Offset(offset).Order("created_at desc").Find(&problems).Error
	return problems, total, err
}

func (r *problemRepository) Create(problem *domain.Problem) error {
	ctx, cancel := database.WithLongTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(problem)
	if result.Error != nil {
		if isUniqueViolation(result.Error) {
			if containsField(result.Error, "title") {
				return fmt.Errorf("title already exists")
			}
			if containsField(result.Error, "slug") {
				return fmt.Errorf("slug already exists")
			}
		}
		return fmt.Errorf("failed to create problem: %w", result.Error)
	}

	return nil
}

func (r *problemRepository) Update(problem *domain.Problem) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	return r.db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update core fields
		// Omit associations here to avoid conflicts or partial updates via Updates
		if err := tx.Model(&domain.Problem{}).
			Where("id = ?", problem.ID).
			Omit("Tags", "Categories", "Creator", "TestCases", "Boilerplates"). // Avoid updating associations via Updates
			Updates(problem).Error; err != nil {
			return fmt.Errorf("failed to update problem fields: %w", err)
		}

		if problem.TestCases != nil {
			if err := tx.Model(problem).Association("TestCases").Replace(problem.TestCases); err != nil {
				return fmt.Errorf("failed to update test cases: %w", err)
			}
		}

		if problem.Tags != nil {
			if err := tx.Model(problem).Association("Tags").Replace(problem.Tags); err != nil {
				return fmt.Errorf("failed to update tags: %w", err)
			}
		}

		if problem.Categories != nil {
			if err := tx.Model(problem).Association("Categories").Replace(problem.Categories); err != nil {
				return fmt.Errorf("failed to update categories: %w", err)
			}
		}

		return nil
	})
}

func (r *problemRepository) Delete(id int) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	return r.db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Clear many-to-many associations (problem_tags, problem_categories)
		var problem domain.Problem
		if err := tx.First(&problem, id).Error; err == nil {
			if err := tx.Model(&problem).Association("Tags").Clear(); err != nil {
				return fmt.Errorf("failed to clear tags: %w", err)
			}
			if err := tx.Model(&problem).Association("Categories").Clear(); err != nil {
				return fmt.Errorf("failed to clear categories: %w", err)
			}
		} else if err != gorm.ErrRecordNotFound {
			return fmt.Errorf("failed to find problem before deletion: %w", err)
		} else {
			return fmt.Errorf("problem not found")
		}

		// 2. Delete all dependent records that refer to this problem
		// List of tables that have a foreign key to problems
		if err := tx.Where("problem_id = ?", id).Delete(&domain.TestCase{}).Error; err != nil {
			return fmt.Errorf("failed to delete test cases: %w", err)
		}
		if err := tx.Where("problem_id = ?", id).Delete(&domain.ProblemBoilerplate{}).Error; err != nil {
			return fmt.Errorf("failed to delete boilerplates: %w", err)
		}
		if err := tx.Where("problem_id = ?", id).Delete(&domain.ProblemReferenceSolution{}).Error; err != nil {
			return fmt.Errorf("failed to delete reference solutions: %w", err)
		}
		if err := tx.Where("problem_id = ?", id).Delete(&domain.ProblemLanguage{}).Error; err != nil {
			return fmt.Errorf("failed to delete problem languages: %w", err)
		}
		if err := tx.Where("problem_id = ?", id).Delete(&domain.Submission{}).Error; err != nil {
			return fmt.Errorf("failed to delete submissions: %w", err)
		}
		if err := tx.Where("problem_id = ?", id).Delete(&domain.UserProblemStats{}).Error; err != nil {
			return fmt.Errorf("failed to delete user problem stats: %w", err)
		}
		if err := tx.Where("problem_id = ?", id).Delete(&domain.ProblemExample{}).Error; err != nil {
			// ignore error if table doesn't exist or other issues with examples
			// since it uses db tags it might be legacy
		}

		// 3. Finally delete the problem itself
		result := tx.Delete(&domain.Problem{}, id)
		if result.Error != nil {
			return fmt.Errorf("failed to delete problem: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("problem not found")
		}

		return nil
	})
}

func (r *problemRepository) GetByID(id int) (*domain.Problem, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	problem := &domain.Problem{}
	err := r.db.DB.WithContext(ctx).
		Preload("Creator").
		Preload("Tags").
		Preload("Categories").
		Preload("TestCases").
		Preload("Boilerplates").
		Preload("Boilerplates.Language").
		First(problem, id).Error

	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("problem not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get problem: %w", err)
	}

	return problem, nil
}

func (r *problemRepository) GetBySlug(slug string) (*domain.Problem, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	problem := &domain.Problem{}
	err := r.db.DB.WithContext(ctx).
		Omit(
			"function_name",
			"return_type",
			"parameters",
			"validation_type",
			"validation_status",
			"expected_time_complexity",
			"expected_space_complexity",
			"has_reference_solution",
			"status",
		).
		Preload("Creator", func(db *gorm.DB) *gorm.DB {
			return db.Omit("created_at", "updated_at", "is_bot", "is_active", "email")
		}).
		Preload("Tags").
		Preload("Categories").
		Preload("TestCases", "is_sample = ?", true).
		Preload("Boilerplates", func(db *gorm.DB) *gorm.DB {
			return db.Omit("test_harness_template", "created_at", "updated_at", "problem_id")
		}).
		Preload("Boilerplates.Language").
		Where("slug = ? AND status = 'published'", slug).
		First(problem).Error

	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("problem not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get problem: %w", err)
	}

	return problem, nil
}

func (r *problemRepository) AdminGetBySlug(slug string) (*domain.Problem, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	problem := &domain.Problem{}
	err := r.db.DB.WithContext(ctx).
		Preload("Creator").
		Preload("Tags").
		Preload("Categories").
		Preload("TestCases").
		Preload("Boilerplates").
		Preload("Boilerplates.Language").
		Where("slug = ?", slug).
		First(problem).Error

	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("problem not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get problem: %w", err)
	}

	return problem, nil
}

func (r *problemRepository) List(filters domain.ProblemFilters) ([]*domain.Problem, int, error) {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	var problems []*domain.Problem
	var total int64

	query := r.db.DB.WithContext(ctx).Model(&domain.Problem{})

	if filters.Difficulty != "" {
		query = query.Where("difficulty = ?", filters.Difficulty)
	}

	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}

	if filters.Visibility != "" {
		query = query.Where("visibility = ?", filters.Visibility)
	}

	if filters.Search != "" {
		searchPattern := "%" + filters.Search + "%"
		query = query.Where("(title ILIKE ? OR description ILIKE ?)", searchPattern, searchPattern)
	}

	if filters.CreatedBy != nil {
		query = query.Where("created_by = ?", *filters.CreatedBy)
	}

	if len(filters.Tags) > 0 {
		// Use subquery for tag filtering
		subQuery := r.db.DB.Table("problem_tags").
			Select("problem_id").
			Joins("JOIN tags ON problem_tags.tag_id = tags.id").
			Where("tags.slug IN ?", filters.Tags)

		query = query.Where("id IN (?)", subQuery)
	}

	if len(filters.Categories) > 0 {
		// Use subquery for category filtering
		subQuery := r.db.DB.Table("problem_categories").
			Select("problem_id").
			Joins("JOIN categories ON problem_categories.category_id = categories.id").
			Where("categories.slug IN ?", filters.Categories)

		query = query.Where("id IN (?)", subQuery)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count problems: %w", err)
	}

	// Pagination
	limit := filters.Limit
	if limit == 0 {
		limit = 20
	}
	offset := (filters.Page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	// Fetch
	if err := query.Preload("Creator").
		Preload("Tags").
		Preload("Categories").
		Preload("TestCases").
		Preload("Boilerplates").
		Preload("Boilerplates.Language").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&problems).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list problems: %w", err)
	}

	return problems, int(total), nil
}

func (r *problemRepository) SlugExists(slug string, excludeID int) (bool, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	var count int64
	query := r.db.DB.WithContext(ctx).Model(&domain.Problem{}).Where("slug = ?", slug)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check slug: %w", err)
	}

	return count > 0, nil
}

func (r *problemRepository) TitleExists(title string) (bool, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	var count int64
	err := r.db.DB.WithContext(ctx).Model(&domain.Problem{}).Where("title = ?", title).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check title: %w", err)
	}

	return count > 0, nil
}

func (r *problemRepository) UpdateCurrentStep(id int, newCurrentStep int) error {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Model(&domain.Problem{}).Where("id = ?", id).Update("current_step", newCurrentStep)
	if result.Error != nil {
		return fmt.Errorf("failed to update current step status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("problem not found")
	}

	return nil
}

func (r *problemRepository) UpdateStats(id int, acceptanceRate float64, totalSubmissions, totalAccepted int) error {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Model(&domain.Problem{}).Where("id = ?", id).Updates(map[string]interface{}{
		"acceptance_rate":   acceptanceRate,
		"total_submissions": totalSubmissions,
		"total_accepted":    totalAccepted,
	})

	if result.Error != nil {
		return fmt.Errorf("failed to update stats: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("problem not found")
	}

	return nil
}

func (r *problemRepository) IncrementStats(id int, isAccepted bool) error {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	return r.db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var problem domain.Problem
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&problem, id).Error; err != nil {
			return err
		}

		updates := map[string]interface{}{
			"total_submissions": gorm.Expr("total_submissions + ?", 1),
		}

		if isAccepted {
			updates["total_accepted"] = gorm.Expr("total_accepted + ?", 1)
		}

		if err := tx.Model(&problem).Updates(updates).Error; err != nil {
			return err
		}

		// Re-fetch to calculate new rate accurately, or just do it in one go.
		// For simplicity, let's just do another query or calculation.
		var updatedProblem domain.Problem
		if err := tx.First(&updatedProblem, id).Error; err != nil {
			return err
		}

		if updatedProblem.TotalSubmissions > 0 {
			rate := (float64(updatedProblem.TotalAccepted) / float64(updatedProblem.TotalSubmissions)) * 100
			if err := tx.Model(&updatedProblem).Update("acceptance_rate", rate).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *problemRepository) UpdateStatus(id int, status string) error {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Model(&domain.Problem{}).Where("id = ?", id).Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("failed to update status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("problem not found")
	}

	return nil
}

func (r *problemRepository) UpdateVisibility(id int, visibility string) error {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	result := r.db.DB.WithContext(ctx).Model(&domain.Problem{}).Where("id = ?", id).Update("visibility", visibility)

	if result.Error != nil {
		return fmt.Errorf("failed to update visibility: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("problem not found")
	}

	return nil
}

func (r *problemRepository) CountProblems() (int, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	var count int64
	err := r.db.DB.WithContext(ctx).Model(&domain.Problem{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count problems: %w", err)
	}

	return int(count), nil
}

func (r *problemRepository) CountByStatus(status string) (int, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	var count int64
	err := r.db.DB.WithContext(ctx).Model(&domain.Problem{}).Where("status = ?", status).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count problems by status: %w", err)
	}

	return int(count), nil
}

func (r *problemRepository) CountByDifficulty(difficulty string) (int, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	var count int64
	err := r.db.DB.WithContext(ctx).Model(&domain.Problem{}).Where("difficulty = ?", difficulty).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count problems by difficulty: %w", err)
	}

	return int(count), nil
}
