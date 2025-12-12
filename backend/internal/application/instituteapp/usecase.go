package instituteapp

import (
	"backend/internal/domain/domain"
	port "backend/internal/port/institute"
	"context"
)

type Usecase struct {
	repo port.InstituteRepository
}

func NewUsecase(repo port.InstituteRepository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) SearchInstitutes(ctx context.Context, query string) ([]domain.Institute, error) {
	return u.repo.SearchInstitutes(ctx, query)
}

func (u *Usecase) GetInstituteById(ctx context.Context, id string) (*domain.Institute, error) {
	return u.repo.FindInstituteById(ctx, id)
}

func (u *Usecase) FindOrCreateInstitute(ctx context.Context, name string) (*domain.Institute, error) {
	return u.repo.FindOrCreateInstitute(ctx, name)
}
