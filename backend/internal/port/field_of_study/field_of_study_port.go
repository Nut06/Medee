package field_of_study

import (
	"backend/internal/domain/domain"
	"context"
)

type FieldOfStudyRepository interface {
	// Existing methods
	CreateFieldOfStudy(ctx context.Context, field *domain.FieldOfStudy) (*domain.FieldOfStudy, error)
	FindFieldOfStudyById(ctx context.Context, id string) (*domain.FieldOfStudy, error)
	FindFieldOfStudyByName(ctx context.Context, name string) (*domain.FieldOfStudy, error)
	FindOrCreateFieldOfStudy(ctx context.Context, name string) (*domain.FieldOfStudy, error)

	// NEW: Search with level filtering for CIP codes
	SearchFieldOfStudies(ctx context.Context, query string, level string) ([]domain.FieldOfStudy, error)
	LoadCIPCodes(ctx context.Context) error
}
