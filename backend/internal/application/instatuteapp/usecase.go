package instatuteapp

import (
	"backend/internal/domain/domain"
	port "backend/internal/port/institue"
	"context"
)

type Usecase struct {
	repo port.InstitueRepository
}

func NewUsecase(repo port.InstitueRepository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) CreateInstitute(ctx context.Context, institute *domain.Institute) (*domain.Institute, error) {
	return u.repo.CreateInstitute(ctx, institute)
}

func (u *Usecase) SearchInstitutes(ctx context.Context, query string) ([]domain.Institute, error) {
	return u.repo.SearchInstitutes(ctx, query)
}