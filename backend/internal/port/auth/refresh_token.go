package authport

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RefreshTokenStore interface {
	SaveRefreshToken(ctx context.Context, token string, userID uuid.UUID, expiresAt time.Time) error
	DeleteRefreshToken(ctx context.Context, token string) error
}
