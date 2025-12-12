package institute_adapter

import (
	"backend/internal/domain/domain"
	"context"

	"gorm.io/gorm"
)

type instituteRepository struct {
	db *gorm.DB
}

func NewInstituteRepository(db *gorm.DB) *instituteRepository {
	return &instituteRepository{db: db}
}

// CreateInstitute creates a new institute
func (r *instituteRepository) CreateInstitute(ctx context.Context, institute *domain.Institute) (*domain.Institute, error) {
	if err := r.db.WithContext(ctx).Create(institute).Error; err != nil {
		return nil, err
	}
	return institute, nil
}

// SearchInstitutes searches institutes by name (for autocomplete)
func (r *instituteRepository) SearchInstitutes(ctx context.Context, query string) ([]domain.Institute, error) {
	var institutes []domain.Institute
	err := r.db.WithContext(ctx).
		Where("name ILIKE ?", "%"+query+"%").
		Limit(10).
		Order("name ASC").
		Find(&institutes).Error
	return institutes, err
}

// FindInstituteById finds institute by ID
func (r *instituteRepository) FindInstituteById(ctx context.Context, id string) (*domain.Institute, error) {
	var institute domain.Institute
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&institute).Error
	return &institute, err
}

// FindInstituteByName finds institute by exact name
func (r *instituteRepository) FindInstituteByName(ctx context.Context, name string) (*domain.Institute, error) {
	var institute domain.Institute
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&institute).Error
	return &institute, err
}

// FindOrCreateInstitute finds existing or creates new institute
func (r *instituteRepository) FindOrCreateInstitute(ctx context.Context, name string) (*domain.Institute, error) {
	// Try to find first
	institute, err := r.FindInstituteByName(ctx, name)
	if err == nil {
		return institute, nil
	}

	// Not found, create new
	institute = &domain.Institute{Name: name}
	return r.CreateInstitute(ctx, institute)
}
