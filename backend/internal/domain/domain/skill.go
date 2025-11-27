package domain

import (
	"time"
	"github.com/google/uuid"
)

type Skill struct {
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name string    `gorm:"uniqueIndex"`
}

type UserSkill struct {
	ID      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID  uuid.UUID `gorm:"type:uuid;index"`
	SkillID uuid.UUID `gorm:"type:uuid;index"`
	Level   *string   // optional: beginner/intermediate/advanced
}

type PortfolioItem struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;index"`
	Title       string
	Description string `gorm:"type:text"`
	ImageURL    *string
	GithubURL   *string
	DemoURL     *string
	CreatedAt   time.Time
}
