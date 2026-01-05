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

func (r *fieldOfStudyRepository) SearchFieldOfStudies(ctx context.Context, query string, level string) ([]domain.FieldOfStudy, error) {
	var fields []domain.FieldOfStudy

	// Build query with level filter
	db := r.db.WithContext(ctx).Where(
		"name ILIKE ? OR code ILIKE ? OR title ILIKE ? OR category ILIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%",
	)

	// Filter by level if specified
	if level != "all" && level != "" {
		db = db.Where("level = ?", level)
	}

	err := db.Limit(50).Order("name ASC").Find(&fields).Error
	return fields, err
}

// LoadCIPCodes - Placeholder for loading CIP codes from static file
// For now, this is a no-op as CIP codes would be loaded separately
func (r *fieldOfStudyRepository) LoadCIPCodes(ctx context.Context) error {
	// TODO: In production, load CIP codes from embedded JSON file
	// and bulk insert into database if not already present
	return nil
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
