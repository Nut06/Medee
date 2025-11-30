package authport

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TokenService interface {
	GenerateAccess(ctx context.Context, userID uuid.UUID) (token string, expiresAt time.Time, err error)
	GenerateRefresh(ctx context.Context, userID uuid.UUID) (token string, expiresAt time.Time, err error)
	DecodeAccessToken(ctx context.Context, token string) (uuid.UUID, error)
	DecodeRefreshToken(ctx context.Context, token string) (uuid.UUID, error)
}