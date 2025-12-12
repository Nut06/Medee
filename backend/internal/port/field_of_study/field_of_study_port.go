package field_of_study

import (
	"backend/internal/domain/domain"
	"context"
)

type FieldOfStudyRepository interface {
	CreateFieldOfStudy(ctx context.Context, field *domain.FieldOfStudy) (*domain.FieldOfStudy, error)
	SearchFieldOfStudies(ctx context.Context, query string) ([]domain.FieldOfStudy, error)
	FindFieldOfStudyById(ctx context.Context, id string) (*domain.FieldOfStudy, error)
	FindFieldOfStudyByName(ctx context.Context, name string) (*domain.FieldOfStudy, error)
	FindOrCreateFieldOfStudy(ctx context.Context, name string) (*domain.FieldOfStudy, error)
}
