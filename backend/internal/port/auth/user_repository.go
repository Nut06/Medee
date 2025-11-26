package authport

import (
	"backend/internal/domain/auth"
	"context"
)

type UserService interface {
	FindByEmail(ctx context.Context, email string) (auth.User, error)
	CreateUser(ctx context.Context, user auth.User) (auth.User, error)
}

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (auth.User, error)
	Create(ctx context.Context, user auth.User) (auth.User, error)
}
