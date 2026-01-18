package postgres

import (
	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/pkg/database"
)

type tagRepository struct {
	db *database.Database
}

func NewTagRepository(db *database.Database) domain.TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Create(tag *domain.Tag) error {
	return r.db.DB.Create(tag).Error
}

func (r *tagRepository) GetByID(id int) (*domain.Tag, error) {
	var tag domain.Tag
	err := r.db.DB.First(&tag, id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) GetBySlug(slug string) (*domain.Tag, error) {
	var tag domain.Tag
	err := r.db.DB.Where("slug = ?", slug).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) Update(tag *domain.Tag) error {
	return r.db.DB.Save(tag).Error
}

func (r *tagRepository) Delete(id int) error {
	return r.db.DB.Delete(&domain.Tag{}, id).Error
}

func (r *tagRepository) List() ([]domain.Tag, error) {
	var tags []domain.Tag
	err := r.db.DB.Order("name asc").Find(&tags).Error
	return tags, err
}
