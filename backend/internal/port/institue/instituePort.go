package instituePort

import (
	"backend/internal/domain/domain"
	"context"
)

type InstitueService interface {
	CreateInstitue(ctx context.Context, institue *domain.Institute) (*domain.Institute, error)
	SearchInstitues(ctx context.Context, query string) ([]domain.Institute, error)
	FindInstitueById(ctx context.Context, id string) (*domain.Institute, error)
	FindInstitueByName(ctx context.Context, name string) (*domain.Institute, error)
}

type InstitueRepository interface {
	CreateInstitue(ctx context.Context, institue *domain.Institute) (*domain.Institute, error)
	SearchInstitues(ctx context.Context, query string) ([]domain.Institute, error)
	FindInstitueById(ctx context.Context, id string) (*domain.Institute, error)
	FindInstitueByName(ctx context.Context, name string) (*domain.Institute, error)
}
