package authapp

import (
	"backend/internal/domain/auth"
	authport "backend/internal/port/auth"
	"context"
	"errors"
)

type Usecase struct {
	users   authport.AuthRepo
	hasher  authport.PasswordHasher
	tokens  authport.TokenService
	refresh authport.RefreshTokenStore
}

func NewUsecase(
	users authport.AuthRepo,
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

func (uc *Usecase) Register(ctx context.Context, cmd authport.RegisterCommand) (*authport.RegisterResult, *auth.TokenPair, error) {
	_, err := uc.users.FindByEmail(ctx, cmd.Email)
	if err == nil {
		return nil, nil, auth.ErrEmailAlreadyUsed
	}
	if !errors.Is(err, auth.ErrUserNotFound) {
		return nil, nil, err
	}

	hash, err := uc.hasher.Hash(ctx, cmd.Password)
	if err != nil {
		return nil, nil, err
	}

	newUser := auth.User{
		FirstName:    cmd.FirstName,
		LastName:     cmd.LastName,
		Email:        cmd.Email,
		PasswordHash: hash,
	}

	created, err := uc.users.Create(ctx, newUser)
	if err != nil {
		return nil, nil, err
	}

	// Generate Tokens
	accessToken, err := uc.tokens.GenerateAccess(ctx, created.ID)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, refreshExp, err := uc.tokens.GenerateRefresh(ctx, created.ID)
	if err != nil {
		return nil, nil, err
	}

	if uc.refresh != nil {
		if err := uc.refresh.SaveRefreshToken(ctx, refreshToken, created.ID, refreshExp); err != nil {
			return nil, nil, err
		}
	}

	return &authport.RegisterResult{
			ID:        created.ID.String(),
			FirstName: created.FirstName,
			LastName:  created.LastName,
			Email:     created.Email,
		}, &auth.TokenPair{
			AccessToken:      accessToken,
			RefreshToken:     refreshToken,
			RefreshExpiresAt: refreshExp,
		}, nil
}

func (uc *Usecase) Login(ctx context.Context, cmd authport.LoginCommand) (*authport.LoginResult, *auth.TokenPair, error) {
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

	accessToken, err := uc.tokens.GenerateAccess(ctx, u.ID)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, refreshExp, err := uc.tokens.GenerateRefresh(ctx, u.ID)
	if err != nil {
		return nil, nil, err
	}

	if uc.refresh != nil {
		if err := uc.refresh.SaveRefreshToken(ctx, refreshToken, u.ID, refreshExp); err != nil {
			return nil, nil, err
		}
	}

	companies, err := uc.users.GetUserCompanies(ctx, u.ID.String())
	if err != nil {
		// Log error but don't fail login? Or fail?
		// For now, let's return empty list if error, or maybe fail.
		// Failing is safer if we rely on this for permissions.
		return nil, nil, err
	}

	return &authport.LoginResult{
			ID:        u.ID.String(),
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Companies: companies,
		}, &auth.TokenPair{
			AccessToken:      accessToken,
			RefreshToken:     refreshToken,
			RefreshExpiresAt: refreshExp,
		}, nil
}

func (uc *Usecase) Logout(ctx context.Context, refreshToken string) error {
	if uc.refresh == nil {
		return nil
	}
	return uc.refresh.DeleteRefreshToken(ctx, refreshToken)
}

func (uc *Usecase) Refresh(ctx context.Context, cmd authport.RefreshCommand) (*authport.LoginResult, *auth.TokenPair, error) {

	userID, err := uc.tokens.DecodeRefreshToken(ctx, cmd.RefreshToken)
	if err != nil {
		return nil, nil, auth.ErrInvalidRefreshToken
	}

	u, err := uc.users.FindByID(ctx, userID.String())
	if err != nil {
		return nil, nil, auth.ErrInvalidRefreshToken
	}

	accessToken, err := uc.tokens.GenerateAccess(ctx, u.ID)
	if err != nil {
		return nil, nil, auth.ErrInvalidRefreshToken
	}
	refreshToken, refreshExp, err := uc.tokens.GenerateRefresh(ctx, u.ID)
	if err != nil {
		return nil, nil, auth.ErrInvalidRefreshToken
	}

	if uc.refresh != nil {
		if err := uc.refresh.SaveRefreshToken(ctx, refreshToken, u.ID, refreshExp); err != nil {
			return nil, nil, auth.ErrInvalidRefreshToken
		}
	}

	companies, err := uc.users.GetUserCompanies(ctx, u.ID.String())
	if err != nil {
		return nil, nil, auth.ErrInvalidRefreshToken
	}

	return &authport.LoginResult{
			ID:        u.ID.String(),
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Companies: companies,
		}, &auth.TokenPair{
			AccessToken:      accessToken,
			RefreshToken:     refreshToken,
			RefreshExpiresAt: refreshExp,
		}, nil
}
