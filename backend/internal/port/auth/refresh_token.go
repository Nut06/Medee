package authport

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RefreshTokenStore interface {
	Save(ctx context.Context, token string, userID uuid.UUID, expiresAt time.Time) error
	Delete(ctx context.Context, token string) error
}
