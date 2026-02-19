package userapp

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
// Test Setup Helpers
// ============================================================

type userTestSuite struct {
	usecase   *Usecase
	userRepo  *mocks.MockUserRepository
	skillRepo *mocks.MockSkillRepository
}

func setupUserTest() *userTestSuite {
	userRepo := new(mocks.MockUserRepository)
	skillRepo := new(mocks.MockSkillRepository)
	uc := NewUsecase(userRepo, skillRepo)

	return &userTestSuite{
		usecase:   uc,
		userRepo:  userRepo,
		skillRepo: skillRepo,
	}
}

// ============================================================
// GetProfile Tests
// ============================================================

func TestGetProfile_Success(t *testing.T) {
	s := setupUserTest()
	ctx := context.Background()
	userID := uuid.New()

	s.userRepo.On("FindById", ctx, userID.String()).Return(&domain.User{
		ID:        userID,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}, nil)

	result, err := s.usecase.GetProfile(ctx, userID.String())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "John", result.FirstName)
	assert.Equal(t, "Doe", result.LastName)
	s.userRepo.AssertExpectations(t)
}

func TestGetProfile_NotFound(t *testing.T) {
	s := setupUserTest()
	ctx := context.Background()

	s.userRepo.On("FindById", ctx, "nonexistent").Return(nil, errors.New("user not found"))

	result, err := s.usecase.GetProfile(ctx, "nonexistent")

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================================
// GetFullProfile Tests
// ============================================================

func TestGetFullProfile_Success(t *testing.T) {
	s := setupUserTest()
	ctx := context.Background()
	userID := uuid.New()
	skillID := uuid.New()

	s.userRepo.On("GetFullProfile", ctx, userID.String()).Return(&domain.User{
		ID:        userID,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		UserSkills: []domain.UserSkill{
			{
				ID:      uuid.New(),
				SkillID: skillID,
				Skill:   &domain.Skill{ID: skillID, Name: "Go"},
			},
		},
		WorkExperiences: []domain.WorkExperience{
			{
				ID:          uuid.New(),
				Position:    "Software Engineer",
				CompanyName: "Tech Corp",
			},
		},
	}, nil)

	result, err := s.usecase.GetFullProfile(ctx, userID.String())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "John", result.FirstName)
	assert.Len(t, result.Skills, 1)
	assert.Len(t, result.Experiences, 1)
}

// ============================================================
// UploadAvatar Tests
// ============================================================

func TestUploadAvatar_Success(t *testing.T) {
	s := setupUserTest()
	ctx := context.Background()
	userID := uuid.New()
	avatarURL := "https://storage.example.com/avatars/user.jpg"

	s.userRepo.On("UploadAvatar", ctx, userID.String(), avatarURL).Return(&domain.User{
		ID:        userID,
		FirstName: "John",
		AvatarURL: &avatarURL,
	}, nil)

	result, err := s.usecase.UploadAvatar(ctx, userID.String(), avatarURL)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, &avatarURL, result.AvatarURL)
}

// ============================================================
// DeleteAvatar Tests
// ============================================================

func TestDeleteAvatar_Success(t *testing.T) {
	s := setupUserTest()
	ctx := context.Background()
	userID := uuid.New()

	s.userRepo.On("DeleteAvatar", ctx, userID.String()).Return(&domain.User{
		ID:        userID,
		FirstName: "John",
		AvatarURL: nil,
	}, nil)

	result, err := s.usecase.DeleteAvatar(ctx, userID.String())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.AvatarURL)
}

// ============================================================
// AddExperience Tests
// ============================================================

func TestAddExperience_Success(t *testing.T) {
	s := setupUserTest()
	ctx := context.Background()
	userID := uuid.New()
	expID := uuid.New()

	req := &domain.WorkExperience{
		Position:    "Software Engineer",
		CompanyName: "Tech Corp",
	}

	s.userRepo.On("AddExperience", ctx, mock.AnythingOfType("*domain.WorkExperience")).Return(&domain.WorkExperience{
		ID:          expID,
		UserID:      userID,
		Position:    "Software Engineer",
		CompanyName: "Tech Corp",
	}, nil)

	result, err := s.usecase.AddExperience(ctx, userID.String(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Software Engineer", result.Position)
	assert.Equal(t, "Tech Corp", result.CompanyName)
}

func TestAddExperience_RepoError(t *testing.T) {
	s := setupUserTest()
	ctx := context.Background()
	userID := uuid.New()

	req := &domain.WorkExperience{
		Position: "Engineer",
	}

	s.userRepo.On("AddExperience", ctx, mock.AnythingOfType("*domain.WorkExperience")).Return(nil, errors.New("db error"))

	result, err := s.usecase.AddExperience(ctx, userID.String(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================================
// AddEducation Tests
// ============================================================

func TestAddEducation_Success(t *testing.T) {
	s := setupUserTest()
	ctx := context.Background()
	userID := uuid.New()
	eduID := uuid.New()

	req := &domain.Education{
		Degree: "Bachelor",
	}

	s.userRepo.On("AddEducation", ctx, mock.AnythingOfType("*domain.Education")).Return(&domain.Education{
		ID:     eduID,
		UserID: userID,
		Degree: "Bachelor",
	}, nil)

	result, err := s.usecase.AddEducation(ctx, userID.String(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Bachelor", result.Degree)
}

// ============================================================
// DeleteExperience / DeleteEducation / DeleteProject Tests
// ============================================================

func TestDeleteExperience_Success(t *testing.T) {
	s := setupUserTest()
	ctx := context.Background()

	s.userRepo.On("DeleteExperience", ctx, "exp-id-123").Return(nil)

	err := s.usecase.DeleteExperience(ctx, "exp-id-123")

	assert.NoError(t, err)
	s.userRepo.AssertExpectations(t)
}

func TestDeleteEducation_Success(t *testing.T) {
	s := setupUserTest()
	ctx := context.Background()

	s.userRepo.On("DeleteEducation", ctx, "edu-id-123").Return(nil)

	err := s.usecase.DeleteEducation(ctx, "edu-id-123")

	assert.NoError(t, err)
}

func TestDeleteProject_Success(t *testing.T) {
	s := setupUserTest()
	ctx := context.Background()

	s.userRepo.On("DeleteProject", ctx, "proj-id-123").Return(nil)

	err := s.usecase.DeleteProject(ctx, "proj-id-123")

	assert.NoError(t, err)
}

// ============================================================
// AddProject Tests
// ============================================================

func TestAddProject_Success(t *testing.T) {
	s := setupUserTest()
	ctx := context.Background()
	userID := uuid.New()
	projID := uuid.New()

	req := &domain.PortfolioItem{
		Title:       "My Project",
		Description: "A cool project",
	}

	s.userRepo.On("AddProject", ctx, mock.AnythingOfType("*domain.PortfolioItem")).Return(&domain.PortfolioItem{
		ID:          projID,
		UserID:      userID,
		Title:       "My Project",
		Description: "A cool project",
	}, nil)

	result, err := s.usecase.AddProject(ctx, userID.String(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "My Project", result.Title)
}

// ============================================================
// AddSkill / DeleteSkill Tests
// ============================================================

func TestDeleteSkill_Success(t *testing.T) {
	s := setupUserTest()
	ctx := context.Background()
	userID := uuid.New().String()
	skillID := uuid.New().String()

	s.userRepo.On("DeleteSkill", ctx, userID, skillID).Return(nil)

	err := s.usecase.DeleteSkill(ctx, userID, skillID)

	assert.NoError(t, err)
	s.userRepo.AssertExpectations(t)
}
