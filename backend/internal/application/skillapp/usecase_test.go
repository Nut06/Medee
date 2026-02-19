package skillapp

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

func setupSkillTest() (*Usecase, *mocks.MockSkillRepository) {
	repo := new(mocks.MockSkillRepository)
	uc := NewUsecase(repo)
	return uc, repo
}

// ============================================================
// CreateSkill Tests
// ============================================================

func TestCreateSkill_Success(t *testing.T) {
	uc, repo := setupSkillTest()
	ctx := context.Background()
	skillID := uuid.New()

	input := &domain.Skill{Name: "Golang"}
	repo.On("CreateSkill", ctx, input).Return(&domain.Skill{
		ID:   skillID,
		Name: "Golang",
	}, nil)

	result, err := uc.CreateSkill(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Golang", result.Name)
	assert.Equal(t, skillID, result.ID)
	repo.AssertExpectations(t)
}

func TestCreateSkill_Error(t *testing.T) {
	uc, repo := setupSkillTest()
	ctx := context.Background()

	input := &domain.Skill{Name: "Golang"}
	repo.On("CreateSkill", ctx, input).Return(nil, errors.New("duplicate skill"))

	result, err := uc.CreateSkill(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================================
// SearchSkills Tests
// ============================================================

func TestSearchSkills_Success(t *testing.T) {
	uc, repo := setupSkillTest()
	ctx := context.Background()

	repo.On("SearchSkills", ctx, "go").Return([]domain.Skill{
		{ID: uuid.New(), Name: "Golang"},
		{ID: uuid.New(), Name: "Go Fiber"},
	}, nil)

	results, err := uc.SearchSkills(ctx, "go")

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "Golang", results[0].Name)
}

func TestSearchSkills_EmptyResult(t *testing.T) {
	uc, repo := setupSkillTest()
	ctx := context.Background()

	repo.On("SearchSkills", ctx, "nonexistent").Return([]domain.Skill{}, nil)

	results, err := uc.SearchSkills(ctx, "nonexistent")

	assert.NoError(t, err)
	assert.Len(t, results, 0)
}

func TestSearchSkills_Error(t *testing.T) {
	uc, repo := setupSkillTest()
	ctx := context.Background()

	repo.On("SearchSkills", ctx, "go").Return([]domain.Skill{}, errors.New("db error"))

	_, err := uc.SearchSkills(ctx, "go")

	assert.Error(t, err)
}
