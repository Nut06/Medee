package domain

import (
	"time"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRole string

const (
	RoleCandidate UserRole = "candidate"
	RoleCompany   UserRole = "company"
)

type User struct {
	// ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ID               uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	Email            string    `json:"email" gorm:"unique"`
	Password         *string   `gorm:"uniqueIndex"`
	PhoneNumber      *string
	Bio              *string
	AvatarURL        *string           `json:"avatar_url"`
	LinkedInURL      *string           `json:"linkedin_url"`
	GitHubURL        *string           `json:"github_url"`
	WebsiteURL       *string           `json:"website_url"`
	UserSkills       []UserSkill       `gorm:"foreignKey:UserID"`
	ApplicantProfile *ApplicantProfile `gorm:"foreignKey:UserID"`
	WorkExperiences  []WorkExperience  `gorm:"foreignKey:UserID"`
	Educations       []Education       `gorm:"foreignKey:UserID"`
	PortfolioItems   []PortfolioItem   `gorm:"foreignKey:UserID"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	RefreshToken     []RefreshToken `gorm:"foreignKey:UserID"`
}

type RefreshToken struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Token     string `gorm:"uniqueIndex"`
	ExpiresAt time.Time
	UserID    uuid.UUID `gorm:"type:uuid"`
	CreatedAt time.Time

	User User `gorm:"foreignKey:UserID"`
}

func (u *User) CheckPassword(pw string) error {
	userpw := *u.Password
	return bcrypt.CompareHashAndPassword([]byte(userpw), []byte(pw))
}
