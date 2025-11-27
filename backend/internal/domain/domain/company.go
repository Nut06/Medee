package domain

import (
	"time"

	"github.com/google/uuid"
)

type Company struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name       string
	LogoURL    *string
	CoverURL   *string
	Industry   *string
	Size       *string // "1-10", "11-50", ...
	Location   *string
	WebsiteURL *string
	About      *string `gorm:"type:text"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CompanyMemberRole string

const (
	RoleRecruiter     CompanyMemberRole = "recruiter"
	RoleHiringManager CompanyMemberRole = "hiring_manager"
	RoleCompanyAdmin  CompanyMemberRole = "admin"
)

type CompanyMember struct {
	ID        uuid.UUID         `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CompanyID uuid.UUID         `gorm:"type:uuid;index"`
	UserID    uuid.UUID         `gorm:"type:uuid;index"`
	Role      CompanyMemberRole `gorm:"type:varchar(30)"`
	CreatedAt time.Time

	Company Company `gorm:"foreignKey:CompanyID"`
	User    User    `gorm:"foreignKey:UserID"`
}
