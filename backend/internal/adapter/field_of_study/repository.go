package field_of_study_adapter

import (
	"backend/internal/domain/domain"
	"context"

	"gorm.io/gorm"
)

type fieldOfStudyRepository struct {
	db *gorm.DB
}

func NewFieldOfStudyRepository(db *gorm.DB) *fieldOfStudyRepository {
	return &fieldOfStudyRepository{db: db}
}

func (r *fieldOfStudyRepository) CreateFieldOfStudy(ctx context.Context, field *domain.FieldOfStudy) (*domain.FieldOfStudy, error) {
	if err := r.db.WithContext(ctx).Create(field).Error; err != nil {
		return nil, err
	}
	return field, nil
}

func (r *fieldOfStudyRepository) SearchFieldOfStudies(ctx context.Context, query string) ([]domain.FieldOfStudy, error) {
	var fields []domain.FieldOfStudy
	err := r.db.WithContext(ctx).
		Where("name ILIKE ?", "%"+query+"%").
		Limit(10).
		Order("name ASC").
		Find(&fields).Error
	return fields, err
}

func (r *fieldOfStudyRepository) FindFieldOfStudyById(ctx context.Context, id string) (*domain.FieldOfStudy, error) {
	var field domain.FieldOfStudy
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&field).Error
	return &field, err
}

func (r *fieldOfStudyRepository) FindFieldOfStudyByName(ctx context.Context, name string) (*domain.FieldOfStudy, error) {
	var field domain.FieldOfStudy
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&field).Error
	return &field, err
}

func (r *fieldOfStudyRepository) FindOrCreateFieldOfStudy(ctx context.Context, name string) (*domain.FieldOfStudy, error) {
	field, err := r.FindFieldOfStudyByName(ctx, name)
	if err == nil {
		return field, nil
	}

	field = &domain.FieldOfStudy{Name: name}
	return r.CreateFieldOfStudy(ctx, field)
}
