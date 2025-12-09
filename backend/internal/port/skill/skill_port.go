package skill_port

import (
	"backend/internal/domain/domain"
	"context"
)

type SkillService interface {
	CreateSkill(ctx context.Context, skill *domain.Skill) (*domain.Skill, error)
	SearchSkills(ctx context.Context, query string) ([]domain.Skill, error)
}

type SkillRepository interface {
	CreateSkill(ctx context.Context, skill *domain.Skill) (*domain.Skill, error)
	SearchSkills(ctx context.Context, query string) ([]domain.Skill, error)
}