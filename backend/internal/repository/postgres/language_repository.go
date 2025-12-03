package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/usecase/interfaces"
	"github.com/prabalesh/loco/backend/pkg/database"
)

type languageRepository struct {
	db *database.Database
}

func NewLanguageRepository(db *database.Database) interfaces.LanguageRepository {
	return &languageRepository{db: db}
}

func (r *languageRepository) Create(ctx context.Context, lang *domain.Language) error {
	ctx, cancel := database.WithLongTimeout()
	defer cancel()

	query := `
        INSERT INTO languages (language_id, name, version, extension, default_template, is_active, executor_config)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at
    `

	err := r.db.QueryRowContext(ctx, query,
		lang.LanguageID,
		lang.Name,
		lang.Version,
		lang.Extension,
		lang.DefaultTemplate,
		lang.IsActive,
		lang.ExecutorConfig,
	).Scan(&lang.ID, &lang.CreatedAt, &lang.UpdatedAt)

	if err != nil {
		if isUniqueViolation(err) {
			if containsField(err, "language_id") {
				return fmt.Errorf("language ID already exists")
			}
		}
		return fmt.Errorf("failed to create language: %w", err)
	}

	return nil
}

func (r *languageRepository) GetByID(ctx context.Context, id int) (*domain.Language, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	lang := &domain.Language{}

	query := `
        SELECT id, language_id, name, version, extension, default_template, 
               is_active, executor_config, created_at, updated_at
        FROM languages WHERE id = $1
    `

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&lang.ID,
		&lang.LanguageID,
		&lang.Name,
		&lang.Version,
		&lang.Extension,
		&lang.DefaultTemplate,
		&lang.IsActive,
		&lang.ExecutorConfig,
		&lang.CreatedAt,
		&lang.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("language not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get language: %w", err)
	}

	return lang, nil
}

func (r *languageRepository) GetByLanguageID(ctx context.Context, languageID string) (*domain.Language, error) {
	ctx, cancel := database.WithShortTimeout()
	defer cancel()

	lang := &domain.Language{}

	query := `
        SELECT id, language_id, name, version, extension, default_template, 
               is_active, executor_config, created_at, updated_at
        FROM languages WHERE language_id = $1
    `

	err := r.db.QueryRowContext(ctx, query, languageID).Scan(
		&lang.ID,
		&lang.LanguageID,
		&lang.Name,
		&lang.Version,
		&lang.Extension,
		&lang.DefaultTemplate,
		&lang.IsActive,
		&lang.ExecutorConfig,
		&lang.CreatedAt,
		&lang.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("language not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get language by language_id: %w", err)
	}

	return lang, nil
}

func (r *languageRepository) Update(ctx context.Context, lang *domain.Language) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	query := `
        UPDATE languages SET
            language_id = $1,
            name = $2,
            version = $3,
            extension = $4,
            default_template = $5,
            is_active = $6,
            executor_config = $7,
            updated_at = NOW()
        WHERE id = $8
        RETURNING updated_at
    `

	err := r.db.QueryRowContext(ctx, query,
		lang.LanguageID,
		lang.Name,
		lang.Version,
		lang.Extension,
		lang.DefaultTemplate,
		lang.IsActive,
		lang.ExecutorConfig,
		lang.ID,
	).Scan(&lang.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("language not found")
	}

	if err != nil {
		return fmt.Errorf("failed to update language: %w", err)
	}

	return nil
}

func (r *languageRepository) Delete(ctx context.Context, id int) error {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	query := `DELETE FROM languages WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete language: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("language not found")
	}

	return nil
}

func (r *languageRepository) ListActive() ([]*domain.Language, error) {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	query := `
        SELECT id, language_id, name, version, extension, default_template, 
               is_active, executor_config, created_at, updated_at
        FROM languages 
        WHERE is_active = true
        ORDER BY name ASC
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list active languages: %w", err)
	}
	defer rows.Close()

	var languages []*domain.Language
	for rows.Next() {
		var lang domain.Language
		err := rows.Scan(
			&lang.ID,
			&lang.LanguageID,
			&lang.Name,
			&lang.Version,
			&lang.Extension,
			&lang.DefaultTemplate,
			&lang.IsActive,
			&lang.ExecutorConfig,
			&lang.CreatedAt,
			&lang.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan language: %w", err)
		}
		languages = append(languages, &lang)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating languages: %w", err)
	}

	return languages, nil
}

func (r *languageRepository) List() ([]*domain.Language, error) {
	ctx, cancel := database.WithMediumTimeout()
	defer cancel()

	query := `
        SELECT id, language_id, name, version, extension, default_template, 
               is_active, executor_config, created_at, updated_at
        FROM languages 
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list languages: %w", err)
	}
	defer rows.Close()

	var languages []*domain.Language
	for rows.Next() {
		var lang domain.Language
		err := rows.Scan(
			&lang.ID,
			&lang.LanguageID,
			&lang.Name,
			&lang.Version,
			&lang.Extension,
			&lang.DefaultTemplate,
			&lang.IsActive,
			&lang.ExecutorConfig,
			&lang.CreatedAt,
			&lang.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan language: %w", err)
		}
		languages = append(languages, &lang)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating languages: %w", err)
	}

	return languages, nil
}
