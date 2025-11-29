package userport

import (
	"backend/internal/domain/domain"
	"context"
	"mime/multipart"
)

type UserService interface {
	GetProfile(ctx context.Context, id string) (*domain.User, error)
	UpdateProfile(ctx context.Context, id string, req *domain.User) (*domain.User, error)
	UploadAvatar(ctx context.Context, id string, file *multipart.FileHeader) (*domain.User, error)
}

type UserRepository interface {
	FindById(ctx context.Context, id string) (*domain.User, error)
	Update(ctx context.Context, id string, req *domain.User) (*domain.User, error)
	UploadAvatar(ctx context.Context, id string, file *multipart.FileHeader) (*domain.User, error)
}
