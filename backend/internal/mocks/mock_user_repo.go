package mocks

import (
	"backend/internal/domain/domain"
	dto "backend/internal/domain/user"
	"context"
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of userport.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindById(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, id string, req *domain.User) (*domain.User, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) UploadAvatar(ctx context.Context, id string, url string) (*domain.User, error) {
	args := m.Called(ctx, id, url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) DeleteAvatar(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) UploadResume(ctx context.Context, id string, file *multipart.FileHeader) (*domain.User, error) {
	args := m.Called(ctx, id, file)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) DeleteResume(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUserCompanies(ctx context.Context, userID string) ([]domain.Company, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Company), args.Error(1)
}

func (m *MockUserRepository) GetFullProfile(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) AddExperience(ctx context.Context, experience *domain.WorkExperience) (*domain.WorkExperience, error) {
	args := m.Called(ctx, experience)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WorkExperience), args.Error(1)
}

func (m *MockUserRepository) UpdateExperience(ctx context.Context, experience *domain.WorkExperience) (*domain.WorkExperience, error) {
	args := m.Called(ctx, experience)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WorkExperience), args.Error(1)
}

func (m *MockUserRepository) DeleteExperience(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) AddEducation(ctx context.Context, education *domain.Education) (*domain.Education, error) {
	args := m.Called(ctx, education)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Education), args.Error(1)
}

func (m *MockUserRepository) UpdateEducation(ctx context.Context, education *domain.Education) (*domain.Education, error) {
	args := m.Called(ctx, education)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Education), args.Error(1)
}

func (m *MockUserRepository) DeleteEducation(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateSkills(ctx context.Context, userID string, req *dto.UpdateSkillsCommand) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *MockUserRepository) AddProject(ctx context.Context, project *domain.PortfolioItem) (*domain.PortfolioItem, error) {
	args := m.Called(ctx, project)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PortfolioItem), args.Error(1)
}

func (m *MockUserRepository) UpdateProject(ctx context.Context, project *domain.PortfolioItem) (*domain.PortfolioItem, error) {
	args := m.Called(ctx, project)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PortfolioItem), args.Error(1)
}

func (m *MockUserRepository) DeleteProject(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) AddSkill(ctx context.Context, userID string, req *dto.AddUserSkillCommand) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteSkill(ctx context.Context, userID string, skillID string) error {
	args := m.Called(ctx, userID, skillID)
	return args.Error(0)
}
