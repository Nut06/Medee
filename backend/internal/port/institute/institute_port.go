package institute

import (
	"backend/internal/domain/domain"
	"context"
)

type InstituteRepository interface {
	CreateInstitute(ctx context.Context, institute *domain.Institute) (*domain.Institute, error)
	SearchInstitutes(ctx context.Context, query string) ([]domain.Institute, error)
	FindInstituteById(ctx context.Context, id string) (*domain.Institute, error)
	FindInstituteByName(ctx context.Context, name string) (*domain.Institute, error)
	FindOrCreateInstitute(ctx context.Context, name string) (*domain.Institute, error)
}
