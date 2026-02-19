package mocks

import (
	"backend/internal/domain/domain"
	"context"

	"github.com/stretchr/testify/mock"
)

// MockSkillRepository is a mock implementation of skillPort.SkillRepository
type MockSkillRepository struct {
	mock.Mock
}

func (m *MockSkillRepository) CreateSkill(ctx context.Context, skill *domain.Skill) (*domain.Skill, error) {
	args := m.Called(ctx, skill)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Skill), args.Error(1)
}

func (m *MockSkillRepository) SearchSkills(ctx context.Context, query string) ([]domain.Skill, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]domain.Skill), args.Error(1)
}

func (m *MockSkillRepository) FindSkillById(ctx context.Context, id string) (*domain.Skill, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Skill), args.Error(1)
}

func (m *MockSkillRepository) FindSkillByName(ctx context.Context, name string) (*domain.Skill, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Skill), args.Error(1)
}
