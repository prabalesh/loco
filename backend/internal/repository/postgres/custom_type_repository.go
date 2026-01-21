package postgres

import (
	"github.com/prabalesh/loco/backend/internal/domain"

	"gorm.io/gorm"
)

type customTypeRepository struct {
	db *gorm.DB
}

func NewCustomTypeRepository(db *gorm.DB) domain.CustomTypeRepository {
	return &customTypeRepository{db: db}
}

func (r *customTypeRepository) Create(customType *domain.CustomType) error {
	return r.db.Create(customType).Error
}

func (r *customTypeRepository) GetByName(name string) (*domain.CustomType, error) {
	var ct domain.CustomType
	if err := r.db.Where("name = ?", name).First(&ct).Error; err != nil {
		return nil, err
	}
	return &ct, nil
}

func (r *customTypeRepository) GetAll() ([]domain.CustomType, error) {
	var cts []domain.CustomType
	if err := r.db.Find(&cts).Error; err != nil {
		return nil, err
	}
	return cts, nil
}
