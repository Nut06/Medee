package skillapp

import (
	"backend/internal/domain/domain"
	port "backend/internal/port/skill"
	"context"
)

type Usecase struct {
	repo port.SkillRepository
}

func NewUsecase(repo port.SkillRepository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) CreateSkill(ctx context.Context, skill *domain.Skill) (*domain.Skill, error) {
	return u.repo.CreateSkill(ctx, skill)
}

func (u *Usecase) SearchSkills(ctx context.Context, query string) ([]domain.Skill, error) {
	return u.repo.SearchSkills(ctx, query)
}
