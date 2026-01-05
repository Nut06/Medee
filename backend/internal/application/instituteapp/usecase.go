package instituteapp

import (
	"backend/internal/domain/domain"
	instituteport "backend/internal/port/institute"
	"context"
	"fmt"
	"time"
	"backend/internal/utils"
)

type Usecase struct {
	repo        instituteport.InstituteRepository
	externalSvc instituteport.InstituteExternalService
	cacheTTL    time.Duration
}

func NewUsecase(repo instituteport.InstituteRepository, externalSvc instituteport.InstituteExternalService) *Usecase {
	return &Usecase{
		repo:        repo,
		externalSvc: externalSvc,
		cacheTTL:    utils.TTL,
	}
}

// SearchInstitutes - Main business logic with caching strategy
func (u *Usecase) SearchInstitutes(ctx context.Context, query string, country string) ([]domain.Institute, error) {
	// 1. Try cache first
	cached, found, err := u.repo.SearchInstitutesFromCache(ctx, query, country)
	if err == nil && found {
		return cached, nil
	}

	// 2. Call external API
	results, err := u.externalSvc.SearchHipoAPI(ctx, query, country)
	if err != nil {
		// Fallback: search in DB
		return u.repo.SearchInstitutesWithCountry(ctx, query, country)
	}

	// 3. Cache results
	cacheKey := fmt.Sprintf("universities:search:%s:%s", query, country)
	_ = u.repo.CacheInstitutes(ctx, cacheKey, results, u.cacheTTL)

	return results, nil
}

// FindOrCreateInstitute - Business logic for creating/finding institute
func (u *Usecase) FindOrCreateInstitute(ctx context.Context, name string) (*domain.Institute, error) {
	return u.repo.FindOrCreateInstitute(ctx, name)
}

// GetInstituteById - Get institute by ID
func (u *Usecase) GetInstituteById(ctx context.Context, id string) (*domain.Institute, error) {
	return u.repo.FindInstituteById(ctx, id)
}
