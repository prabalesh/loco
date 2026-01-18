package postgres

import (
	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/pkg/database"
)

type categoryRepository struct {
	db *database.Database
}

func NewCategoryRepository(db *database.Database) domain.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *domain.Category) error {
	return r.db.DB.Create(category).Error
}

func (r *categoryRepository) GetByID(id int) (*domain.Category, error) {
	var category domain.Category
	err := r.db.DB.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetBySlug(slug string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.DB.Where("slug = ?", slug).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Update(category *domain.Category) error {
	return r.db.DB.Save(category).Error
}

func (r *categoryRepository) Delete(id int) error {
	return r.db.DB.Delete(&domain.Category{}, id).Error
}

func (r *categoryRepository) List() ([]domain.Category, error) {
	var categories []domain.Category
	err := r.db.DB.Order("name asc").Find(&categories).Error
	return categories, err
}
