package domain

import (
	"time"

	"github.com/google/uuid"
)

type Provider string

const (
	ProviderGoogle   Provider = "google"
	ProviderGithub   Provider = "github"
	ProviderLinkedIn Provider = "linkedin"
)

type OAuthAccount struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID         uuid.UUID `gorm:"type:uuid;index"`
	Provider       Provider  `gorm:"type:varchar(20)"`
	ProviderUserID string    `gorm:"index"`
	AccessToken    *string   `gorm:"type:text"`
	RefreshToken   *string   `gorm:"type:text"`
	ExpiresAt      *time.Time
}
