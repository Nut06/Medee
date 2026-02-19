package mocks

import (
	"backend/internal/domain/domain"
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

// MockInstituteRepository is a mock implementation of institute.InstituteRepository
type MockInstituteRepository struct {
	mock.Mock
}

func (m *MockInstituteRepository) CreateInstitute(ctx context.Context, institute *domain.Institute) (*domain.Institute, error) {
	args := m.Called(ctx, institute)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Institute), args.Error(1)
}

func (m *MockInstituteRepository) SearchInstitutes(ctx context.Context, query string) ([]domain.Institute, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]domain.Institute), args.Error(1)
}

func (m *MockInstituteRepository) FindInstituteById(ctx context.Context, id string) (*domain.Institute, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Institute), args.Error(1)
}

func (m *MockInstituteRepository) FindInstituteByName(ctx context.Context, name string) (*domain.Institute, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Institute), args.Error(1)
}

func (m *MockInstituteRepository) FindOrCreateInstitute(ctx context.Context, name string) (*domain.Institute, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Institute), args.Error(1)
}

func (m *MockInstituteRepository) SearchInstitutesFromCache(ctx context.Context, query string, country string) ([]domain.Institute, bool, error) {
	args := m.Called(ctx, query, country)
	return args.Get(0).([]domain.Institute), args.Bool(1), args.Error(2)
}

func (m *MockInstituteRepository) CacheInstitutes(ctx context.Context, key string, institutes []domain.Institute, ttl time.Duration) error {
	args := m.Called(ctx, key, institutes, ttl)
	return args.Error(0)
}

func (m *MockInstituteRepository) SearchInstitutesWithCountry(ctx context.Context, query string, country string) ([]domain.Institute, error) {
	args := m.Called(ctx, query, country)
	return args.Get(0).([]domain.Institute), args.Error(1)
}
