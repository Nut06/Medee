package authapp

import (
	"backend/internal/domain/auth"
	"backend/internal/mocks"
	authport "backend/internal/port/auth"
	"context"
	"errors"
	"testing"
	"time"
	"backend/internal/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================
// Test Setup Helpers
// ============================================================

type authTestSuite struct {
	usecase *Usecase
	repo    *mocks.MockAuthRepo
	hasher  *mocks.MockPasswordHasher
	tokens  *mocks.MockTokenService
}

func setupAuthTest() *authTestSuite {
	repo := new(mocks.MockAuthRepo)
	hasher := new(mocks.MockPasswordHasher)
	tokens := new(mocks.MockTokenService)

	uc := NewUsecase(repo, hasher, tokens)

	return &authTestSuite{
		usecase: uc,
		repo:    repo,
		hasher:  hasher,
		tokens:  tokens,
		refresh: refresh,
	}
}

// ============================================================
// Register Tests
// ============================================================

func TestRegister_Success(t *testing.T) {
	s := setupAuthTest()
	ctx := context.Background()
	userID := uuid.New()
	refreshExp := time.Now().Add(utils.SevenDays)

	cmd := authport.RegisterCommand{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "password123",
	}

	// Email doesn't exist
	s.repo.On("FindByEmail", ctx, cmd.Email).Return(auth.User{}, auth.ErrUserNotFound)
	// Hash password
	s.hasher.On("Hash", ctx, cmd.Password).Return("hashed_password", nil)
	// Create user
	s.repo.On("Create", ctx, mock.AnythingOfType("auth.User")).Return(auth.User{
		ID:        userID,
		FirstName: cmd.FirstName,
		LastName:  cmd.LastName,
		Email:     cmd.Email,
	}, nil)
	// Generate tokens
	s.tokens.On("GenerateAccess", ctx, userID).Return("access_token", nil)
	s.tokens.On("GenerateRefresh", ctx, userID).Return("refresh_token", refreshExp, nil)

	result, tokenPair, err := s.usecase.Register(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, tokenPair)
	assert.Equal(t, cmd.FirstName, result.FirstName)
	assert.Equal(t, cmd.LastName, result.LastName)
	assert.Equal(t, cmd.Email, result.Email)
	assert.Equal(t, "access_token", tokenPair.AccessToken)
	assert.Equal(t, "refresh_token", tokenPair.RefreshToken)
	s.repo.AssertExpectations(t)
	s.hasher.AssertExpectations(t)
	s.tokens.AssertExpectations(t)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	s := setupAuthTest()
	ctx := context.Background()

	cmd := authport.RegisterCommand{
		Email:    "existing@example.com",
		Password: "password123",
	}

	// Email already exists (FindByEmail returns a user, no error)
	s.repo.On("FindByEmail", ctx, cmd.Email).Return(auth.User{
		Email: cmd.Email,
	}, nil)

	result, tokenPair, err := s.usecase.Register(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, tokenPair)
	assert.True(t, errors.Is(err, auth.ErrEmailAlreadyUsed))
}

func TestRegister_HashError(t *testing.T) {
	s := setupAuthTest()
	ctx := context.Background()

	cmd := authport.RegisterCommand{
		Email:    "john@example.com",
		Password: "password123",
	}

	s.repo.On("FindByEmail", ctx, cmd.Email).Return(auth.User{}, auth.ErrUserNotFound)
	s.hasher.On("Hash", ctx, cmd.Password).Return("", errors.New("hash failed"))

	result, tokenPair, err := s.usecase.Register(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, tokenPair)
}

// ============================================================
// Login Tests
// ============================================================

func TestLogin_Success(t *testing.T) {
	s := setupAuthTest()
	ctx := context.Background()
	userID := uuid.New()
	refreshExp := time.Now().Add(7 * 24 * time.Hour)

	cmd := authport.LoginCommand{
		Email:    "john@example.com",
		Password: "password123",
	}

	s.repo.On("FindByEmail", ctx, cmd.Email).Return(auth.User{
		ID:           userID,
		FirstName:    "John",
		LastName:     "Doe",
		Email:        cmd.Email,
		PasswordHash: "hashed_password",
	}, nil)
	s.hasher.On("Compare", ctx, "hashed_password", cmd.Password).Return(nil)
	s.tokens.On("GenerateAccess", ctx, userID).Return("access_token", nil)
	s.tokens.On("GenerateRefresh", ctx, userID).Return("refresh_token", refreshExp, nil)
	s.repo.On("GetUserCompanies", ctx, userID.String()).Return([]auth.Company{}, nil)

	result, tokenPair, err := s.usecase.Login(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, tokenPair)
	assert.Equal(t, "John", result.FirstName)
	assert.Equal(t, "access_token", tokenPair.AccessToken)
}

func TestLogin_UserNotFound(t *testing.T) {
	s := setupAuthTest()
	ctx := context.Background()

	cmd := authport.LoginCommand{
		Email:    "unknown@example.com",
		Password: "password123",
	}

	s.repo.On("FindByEmail", ctx, cmd.Email).Return(auth.User{}, auth.ErrUserNotFound)

	result, tokenPair, err := s.usecase.Login(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, tokenPair)
	assert.True(t, errors.Is(err, auth.ErrInvalidCredential))
}

func TestLogin_WrongPassword(t *testing.T) {
	s := setupAuthTest()
	ctx := context.Background()
	userID := uuid.New()

	cmd := authport.LoginCommand{
		Email:    "john@example.com",
		Password: "wrong_password",
	}

	s.repo.On("FindByEmail", ctx, cmd.Email).Return(auth.User{
		ID:           userID,
		PasswordHash: "hashed_password",
	}, nil)
	s.hasher.On("Compare", ctx, "hashed_password", "wrong_password").Return(auth.ErrInvalidCredential)

	result, tokenPair, err := s.usecase.Login(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, tokenPair)
	assert.True(t, errors.Is(err, auth.ErrInvalidCredential))
}

// ============================================================
// Logout Tests
// ============================================================

func TestLogout_NilRefreshStore(t *testing.T) {
	repo := new(mocks.MockAuthRepo)
	hasher := new(mocks.MockPasswordHasher)
	tokens := new(mocks.MockTokenService)
	uc := NewUsecase(repo, hasher, tokens, nil)
	ctx := context.Background()

	err := uc.Logout(ctx, "some_refresh_token")

	assert.NoError(t, err)
}

// ============================================================
// Refresh Tests
// ============================================================

func TestRefresh_Success(t *testing.T) {
	s := setupAuthTest()
	ctx := context.Background()
	userID := uuid.New()
	refreshExp := time.Now().Add(7 * 24 * time.Hour)

	cmd := authport.RefreshCommand{
		RefreshToken: "old_refresh_token",
	}

	s.tokens.On("DecodeRefreshToken", ctx, cmd.RefreshToken).Return(userID, nil)
	s.repo.On("FindByID", ctx, userID.String()).Return(auth.User{
		ID:        userID,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}, nil)
	s.tokens.On("GenerateAccess", ctx, userID).Return("new_access_token", nil)
	s.tokens.On("GenerateRefresh", ctx, userID).Return("new_refresh_token", refreshExp, nil)
	s.repo.On("GetUserCompanies", ctx, userID.String()).Return([]auth.Company{}, nil)

	result, tokenPair, err := s.usecase.Refresh(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, tokenPair)
	assert.Equal(t, "new_access_token", tokenPair.AccessToken)
	assert.Equal(t, "new_refresh_token", tokenPair.RefreshToken)
}

func TestRefresh_InvalidToken(t *testing.T) {
	s := setupAuthTest()
	ctx := context.Background()

	cmd := authport.RefreshCommand{
		RefreshToken: "invalid_token",
	}

	s.tokens.On("DecodeRefreshToken", ctx, cmd.RefreshToken).Return(uuid.Nil, errors.New("invalid"))

	result, tokenPair, err := s.usecase.Refresh(ctx, cmd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, tokenPair)
	assert.True(t, errors.Is(err, auth.ErrInvalidRefreshToken))
}
