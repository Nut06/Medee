package userapp

import (
	"backend/internal/domain/domain"
	port "backend/internal/port/user"
	"context"
	"mime/multipart"
)

type UpdateProfileCommand struct {
	FirstName string
	LastName  string
	Email     string
}



type Usecase struct {
	repo port.UserRepository
}

func NewUsecase(repo port.UserRepository) *Usecase {
	return &Usecase{repo: repo}
}

func (s *Usecase) GetProfile(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.FindById(ctx, id)
}

func (s *Usecase) UpdateProfile(ctx context.Context, id string, req *domain.User) (*domain.User, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *Usecase) UploadAvatar(ctx context.Context, id string, file *multipart.FileHeader) (*domain.User, error) {
	return s.repo.UploadAvatar(ctx, id, file)
}

func (s *Usecase) DeleteAvatar(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.DeleteAvatar(ctx, id)
}
