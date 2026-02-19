package fieldapp

import (
	"backend/internal/domain/domain"
	"backend/internal/mocks"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// ============================================================
// Test Setup
// ============================================================

func setupFieldTest() (*Usecase, *mocks.MockFieldOfStudyRepository) {
	repo := new(mocks.MockFieldOfStudyRepository)
	uc := NewUsecase(repo)
	return uc, repo
}

// ============================================================
// SearchFieldOfStudies Tests
// ============================================================

func TestSearchFieldOfStudies_Success(t *testing.T) {
	uc, repo := setupFieldTest()
	ctx := context.Background()

	repo.On("SearchFieldOfStudies", ctx, "computer", "").Return([]domain.FieldOfStudy{
		{ID: uuid.New(), Name: "Computer Science"},
		{ID: uuid.New(), Name: "Computer Engineering"},
	}, nil)

	results, err := uc.SearchFieldOfStudies(ctx, "computer", "")

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "Computer Science", results[0].Name)
}

func TestSearchFieldOfStudies_WithLevel(t *testing.T) {
	uc, repo := setupFieldTest()
	ctx := context.Background()

	repo.On("SearchFieldOfStudies", ctx, "engineering", "bachelor").Return([]domain.FieldOfStudy{
		{ID: uuid.New(), Name: "Electrical Engineering"},
	}, nil)

	results, err := uc.SearchFieldOfStudies(ctx, "engineering", "bachelor")

	assert.NoError(t, err)
	assert.Len(t, results, 1)
}

func TestSearchFieldOfStudies_Error(t *testing.T) {
	uc, repo := setupFieldTest()
	ctx := context.Background()

	repo.On("SearchFieldOfStudies", ctx, "test", "").Return([]domain.FieldOfStudy{}, errors.New("db error"))

	_, err := uc.SearchFieldOfStudies(ctx, "test", "")

	assert.Error(t, err)
}

// ============================================================
// GetFieldOfStudyById Tests
// ============================================================

func TestGetFieldOfStudyById_Success(t *testing.T) {
	uc, repo := setupFieldTest()
	ctx := context.Background()
	id := uuid.New()

	repo.On("FindFieldOfStudyById", ctx, id.String()).Return(&domain.FieldOfStudy{
		ID:   id,
		Name: "Computer Science",
	}, nil)

	result, err := uc.GetFieldOfStudyById(ctx, id.String())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Computer Science", result.Name)
}

func TestGetFieldOfStudyById_NotFound(t *testing.T) {
	uc, repo := setupFieldTest()
	ctx := context.Background()

	repo.On("FindFieldOfStudyById", ctx, "nonexistent").Return(nil, errors.New("not found"))

	result, err := uc.GetFieldOfStudyById(ctx, "nonexistent")

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================================
// FindOrCreateFieldOfStudy Tests
// ============================================================

func TestFindOrCreateFieldOfStudy_Success(t *testing.T) {
	uc, repo := setupFieldTest()
	ctx := context.Background()

	repo.On("FindOrCreateFieldOfStudy", ctx, "Data Science").Return(&domain.FieldOfStudy{
		ID:   uuid.New(),
		Name: "Data Science",
	}, nil)

	result, err := uc.FindOrCreateFieldOfStudy(ctx, "Data Science")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Data Science", result.Name)
}

// ============================================================
// LoadCIPCodes Tests
// ============================================================

func TestLoadCIPCodes_Success(t *testing.T) {
	uc, repo := setupFieldTest()
	ctx := context.Background()

	repo.On("LoadCIPCodes", ctx).Return(nil)

	err := uc.LoadCIPCodes(ctx)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLoadCIPCodes_Error(t *testing.T) {
	uc, repo := setupFieldTest()
	ctx := context.Background()

	repo.On("LoadCIPCodes", ctx).Return(errors.New("failed to load"))

	err := uc.LoadCIPCodes(ctx)

	assert.Error(t, err)
}
