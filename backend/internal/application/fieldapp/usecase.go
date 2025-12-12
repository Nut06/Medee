package fieldapp

import (
	"backend/internal/domain/domain"
	port "backend/internal/port/field_of_study"
	"context"
)

type Usecase struct {
	repo port.FieldOfStudyRepository
}

func NewUsecase(repo port.FieldOfStudyRepository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) SearchFieldOfStudies(ctx context.Context, query string) ([]domain.FieldOfStudy, error) {
	return u.repo.SearchFieldOfStudies(ctx, query)
}

func (u *Usecase) GetFieldOfStudyById(ctx context.Context, id string) (*domain.FieldOfStudy, error) {
	return u.repo.FindFieldOfStudyById(ctx, id)
}

func (u *Usecase) FindOrCreateFieldOfStudy(ctx context.Context, name string) (*domain.FieldOfStudy, error) {
	return u.repo.FindOrCreateFieldOfStudy(ctx, name)
}
