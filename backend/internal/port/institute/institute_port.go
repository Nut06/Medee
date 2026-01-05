package institute

import (
	"backend/internal/domain/domain"
	"context"
	"time"
)

type InstituteRepository interface {
	// Existing methods
	CreateInstitute(ctx context.Context, institute *domain.Institute) (*domain.Institute, error)
	SearchInstitutes(ctx context.Context, query string) ([]domain.Institute, error)
	FindInstituteById(ctx context.Context, id string) (*domain.Institute, error)
	FindInstituteByName(ctx context.Context, name string) (*domain.Institute, error)
	FindOrCreateInstitute(ctx context.Context, name string) (*domain.Institute, error)

	// NEW: Cache methods for external API integration
	SearchInstitutesFromCache(ctx context.Context, query string, country string) ([]domain.Institute, bool, error)
	CacheInstitutes(ctx context.Context, key string, institutes []domain.Institute, ttl time.Duration) error
	SearchInstitutesWithCountry(ctx context.Context, query string, country string) ([]domain.Institute, error)
}

// NEW: External API client interface
type InstituteExternalService interface {
	SearchHipoAPI(ctx context.Context, query string, country string) ([]domain.Institute, error)
}
