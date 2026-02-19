package mocks

import (
	"backend/internal/domain/domain"
	"context"

	"github.com/stretchr/testify/mock"
)

// MockFieldOfStudyRepository is a mock implementation of field_of_study.FieldOfStudyRepository
type MockFieldOfStudyRepository struct {
	mock.Mock
}

func (m *MockFieldOfStudyRepository) CreateFieldOfStudy(ctx context.Context, field *domain.FieldOfStudy) (*domain.FieldOfStudy, error) {
	args := m.Called(ctx, field)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FieldOfStudy), args.Error(1)
}

func (m *MockFieldOfStudyRepository) FindFieldOfStudyById(ctx context.Context, id string) (*domain.FieldOfStudy, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FieldOfStudy), args.Error(1)
}

func (m *MockFieldOfStudyRepository) FindFieldOfStudyByName(ctx context.Context, name string) (*domain.FieldOfStudy, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FieldOfStudy), args.Error(1)
}

func (m *MockFieldOfStudyRepository) FindOrCreateFieldOfStudy(ctx context.Context, name string) (*domain.FieldOfStudy, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FieldOfStudy), args.Error(1)
}

func (m *MockFieldOfStudyRepository) SearchFieldOfStudies(ctx context.Context, query string, level string) ([]domain.FieldOfStudy, error) {
	args := m.Called(ctx, query, level)
	return args.Get(0).([]domain.FieldOfStudy), args.Error(1)
}

func (m *MockFieldOfStudyRepository) LoadCIPCodes(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
