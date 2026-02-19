package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockPasswordHasher is a mock implementation of authport.PasswordHasher
type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockPasswordHasher) Hash(ctx context.Context, password string) (string, error) {
	args := m.Called(ctx, password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordHasher) Compare(ctx context.Context, hash, password string) error {
	args := m.Called(ctx, hash, password)
	return args.Error(0)
}
