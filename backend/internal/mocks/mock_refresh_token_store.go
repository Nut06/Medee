package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockRefreshTokenStore is a mock implementation of authport.RefreshTokenStore
type MockRefreshTokenStore struct {
	mock.Mock
}

func (m *MockRefreshTokenStore) SaveRefreshToken(ctx context.Context, token string, userID uuid.UUID, expiresAt time.Time) error {
	args := m.Called(ctx, token, userID, expiresAt)
	return args.Error(0)
}

func (m *MockRefreshTokenStore) DeleteRefreshToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}
