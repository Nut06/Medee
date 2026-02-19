package instituteapp

import (
	"backend/internal/domain/domain"
	"backend/internal/mocks"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================
// Test Setup
// ============================================================

func setupInstituteTest() (*Usecase, *mocks.MockInstituteRepository, *mocks.MockInstituteExternalService) {
	repo := new(mocks.MockInstituteRepository)
	externalSvc := new(mocks.MockInstituteExternalService)
	uc := NewUsecase(repo, externalSvc)
	return uc, repo, externalSvc
}

// ============================================================
// SearchInstitutes Tests
// ============================================================

func TestSearchInstitutes_CacheHit(t *testing.T) {
	uc, repo, _ := setupInstituteTest()
	ctx := context.Background()

	cached := []domain.Institute{
		{ID: uuid.New(), Name: "Chulalongkorn University"},
	}

	repo.On("SearchInstitutesFromCache", ctx, "chula", "TH").Return(cached, true, nil)

	results, err := uc.SearchInstitutes(ctx, "chula", "TH")

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Chulalongkorn University", results[0].Name)
	repo.AssertExpectations(t)
}

func TestSearchInstitutes_CacheMiss_ExternalAPISuccess(t *testing.T) {
	uc, repo, externalSvc := setupInstituteTest()
	ctx := context.Background()

	apiResults := []domain.Institute{
		{ID: uuid.New(), Name: "MIT"},
		{ID: uuid.New(), Name: "Harvard"},
	}

	// Cache miss
	repo.On("SearchInstitutesFromCache", ctx, "mit", "US").Return([]domain.Institute{}, false, nil)
	// External API success
	externalSvc.On("SearchHipoAPI", ctx, "mit", "US").Return(apiResults, nil)
	// Cache the results
	repo.On("CacheInstitutes", ctx, "universities:search:mit:US", apiResults, mock.Anything).Return(nil)

	results, err := uc.SearchInstitutes(ctx, "mit", "US")

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "MIT", results[0].Name)
	repo.AssertExpectations(t)
	externalSvc.AssertExpectations(t)
}

func TestSearchInstitutes_CacheMiss_FallbackToDB(t *testing.T) {
	uc, repo, externalSvc := setupInstituteTest()
	ctx := context.Background()

	dbResults := []domain.Institute{
		{ID: uuid.New(), Name: "Local University"},
	}

	// Cache miss
	repo.On("SearchInstitutesFromCache", ctx, "local", "TH").Return([]domain.Institute{}, false, nil)
	// External API fails
	externalSvc.On("SearchHipoAPI", ctx, "local", "TH").Return([]domain.Institute{}, errors.New("api error"))
	// Fallback to DB
	repo.On("SearchInstitutesWithCountry", ctx, "local", "TH").Return(dbResults, nil)

	results, err := uc.SearchInstitutes(ctx, "local", "TH")

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Local University", results[0].Name)
}

// ============================================================
// GetInstituteById Tests
// ============================================================

func TestGetInstituteById_Success(t *testing.T) {
	uc, repo, _ := setupInstituteTest()
	ctx := context.Background()
	id := uuid.New()

	repo.On("FindInstituteById", ctx, id.String()).Return(&domain.Institute{
		ID:   id,
		Name: "Kasetsart University",
	}, nil)

	result, err := uc.GetInstituteById(ctx, id.String())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Kasetsart University", result.Name)
}

func TestGetInstituteById_NotFound(t *testing.T) {
	uc, repo, _ := setupInstituteTest()
	ctx := context.Background()

	repo.On("FindInstituteById", ctx, "nonexistent").Return(nil, errors.New("not found"))

	result, err := uc.GetInstituteById(ctx, "nonexistent")

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================================
// FindOrCreateInstitute Tests
// ============================================================

func TestFindOrCreateInstitute_Success(t *testing.T) {
	uc, repo, _ := setupInstituteTest()
	ctx := context.Background()

	repo.On("FindOrCreateInstitute", ctx, "New University").Return(&domain.Institute{
		ID:   uuid.New(),
		Name: "New University",
	}, nil)

	result, err := uc.FindOrCreateInstitute(ctx, "New University")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "New University", result.Name)
}
