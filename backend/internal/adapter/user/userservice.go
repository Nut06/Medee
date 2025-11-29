package useradapter

import (
	"backend/internal/domain/domain"
	"context"
	"mime/multipart"
)

type UserService struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) UpdateProfile(ctx context.Context, id string, req *domain.User) (*domain.User, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *UserService) GetProfile(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.FindById(ctx, id)
}

func (s *UserService) UploadAvatar(ctx context.Context, id string, file *multipart.FileHeader) (*domain.User, error) {
	return s.repo.UploadAvatar(ctx, id, file)
}
