package mocks

import (
	"backend/internal/domain/domain"
	"context"

	"github.com/stretchr/testify/mock"
)

// MockInstituteExternalService is a mock implementation of institute.InstituteExternalService
type MockInstituteExternalService struct {
	mock.Mock
}

func (m *MockInstituteExternalService) SearchHipoAPI(ctx context.Context, query string, country string) ([]domain.Institute, error) {
	args := m.Called(ctx, query, country)
	return args.Get(0).([]domain.Institute), args.Error(1)
}
