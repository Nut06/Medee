package authapp

import (
	"backend/internal/domain/auth"
	authport "backend/internal/port/auth"
	"context"
	"errors"
)

type Usecase struct {
	users   authport.UserRepository
	hasher  authport.PasswordHasher
	tokens  authport.TokenService
	refresh authport.RefreshTokenStore
}

func NewUsecase(
	users authport.UserRepository,
	hasher authport.PasswordHasher,
	tokens authport.TokenService,
	refresh authport.RefreshTokenStore,
) *Usecase {
	return &Usecase{
		users:   users,
		hasher:  hasher,
		tokens:  tokens,
		refresh: refresh,
	}
}

type RegisterCommand struct {
	FirstName     string
	LastName     string
	Email    string
	Password string
	Role     string
}

type RegisterResult struct {
	ID    string
	FirstName string
	LastName string
	Email string
	Role  string
}

type LoginCommand struct {
	Email    string
	Password string
}

type LoginResult struct {
	ID    string
	FirstName  string
	LastName string
	Email string
	Role  string
}

func (uc *Usecase) Register(ctx context.Context, cmd RegisterCommand) (*RegisterResult, error) {
	_, err := uc.users.FindByEmail(ctx, cmd.Email)
	if err == nil {
		return nil, auth.ErrEmailAlreadyUsed
	}
	if !errors.Is(err, auth.ErrUserNotFound) {
		return nil, err
	}

	hash, err := uc.hasher.Hash(ctx, cmd.Password)
	if err != nil {
		return nil, err
	}

	role := cmd.Role
	if role == "" {
		role = "candidate"
	}

	newUser := auth.User{
		FirstName:    cmd.FirstName,
		LastName:     cmd.LastName,
		Email:        cmd.Email,
		Role:         role,
		PasswordHash: hash,
	}

	created, err := uc.users.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return &RegisterResult{
		ID:    created.ID.String(),
		FirstName:  created.FirstName,
		LastName: created.LastName,
		Email: created.Email,
		Role:  created.Role,
	}, nil
}

func (uc *Usecase) Login(ctx context.Context, cmd LoginCommand) (*LoginResult, *auth.TokenPair, error) {
	u, err := uc.users.FindByEmail(ctx, cmd.Email)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			return nil, nil, auth.ErrInvalidCredential
		}
		return nil, nil, err
	}

	if err := uc.hasher.Compare(ctx, u.PasswordHash, cmd.Password); err != nil {
		return nil, nil, auth.ErrInvalidCredential
	}

	accessToken, accessExp, err := uc.tokens.GenerateAccess(ctx, u.ID)
	if err != nil {
		return nil, nil, err
	}
	refreshToken, refreshExp, err := uc.tokens.GenerateRefresh(ctx, u.ID)
	if err != nil {
		return nil, nil, err
	}

	if uc.refresh != nil {
		if err := uc.refresh.Save(ctx, refreshToken, u.ID, refreshExp); err != nil {
			return nil, nil, err
		}
	}

	return &LoginResult{
			ID:    u.ID.String(),
			FirstName:  u.FirstName,
			LastName: u.LastName,
			Email: u.Email,
			Role:  u.Role,
		}, &auth.TokenPair{
			AccessToken:      accessToken,
			RefreshToken:     refreshToken,
			AccessExpiresAt:  accessExp,
			RefreshExpiresAt: refreshExp,
		}, nil
}

func (uc *Usecase) Logout(ctx context.Context, refreshToken string) error {
	if uc.refresh == nil {
		return nil
	}
	return uc.refresh.Delete(ctx, refreshToken)
}
