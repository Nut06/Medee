package userapp

import (
	"backend/internal/domain/domain"
	"backend/internal/port/user"
	"context"
)

type Usecase struct {
	users   userport.UserRepository
	service userport.UserService
}

func NewUsecase(users userport.UserRepository, service userport.UserService) *Usecase {
	if users == nil {
		panic("user repository is nil")
	}
	if service == nil {
		panic("user service is nil")
	}
	return &Usecase{users: users, service: service}
}

func (uc *Usecase) GetProfile(ctx context.Context, id string) (*domain.User, error) {
	return uc.users.FindById(ctx, id)
}



func (uc *Usecase) UpdateProfile(ctx context.Context, id string, req *domain.User) (*domain.User, error) {
	return uc.users.Update(ctx, id, req)
}
