package mocks

import (
	"backend/internal/domain/auth"
	"context"

	"github.com/stretchr/testify/mock"
)

// MockAuthRepo is a mock implementation of authport.AuthRepo
type MockAuthRepo struct {
	mock.Mock
}

func (m *MockAuthRepo) FindByEmail(ctx context.Context, email string) (auth.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(auth.User), args.Error(1)
}

func (m *MockAuthRepo) Create(ctx context.Context, user auth.User) (auth.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(auth.User), args.Error(1)
}

func (m *MockAuthRepo) GetUserCompanies(ctx context.Context, userID string) ([]auth.Company, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]auth.Company), args.Error(1)
}

func (m *MockAuthRepo) FindByID(ctx context.Context, id string) (auth.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(auth.User), args.Error(1)
}
