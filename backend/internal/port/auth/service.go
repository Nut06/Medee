package authport

import (
	"backend/internal/domain/auth"
	"context"
)

type AuthService interface {
	Register(ctx context.Context, cmd RegisterCommand) (*RegisterResult, *auth.TokenPair, error)
	Login(ctx context.Context, cmd LoginCommand) (*LoginResult, *auth.TokenPair, error)
	Refresh(ctx context.Context, cmd RefreshCommand) (*LoginResult, *auth.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
}