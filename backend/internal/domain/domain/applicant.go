package domain

import (
	"time"

	"github.com/google/uuid"
)

type ApplicantProfile struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;uniqueIndex"` // one-to-one
	Summary     *string   `gorm:"type:text"`             // About Me
	Address     *string
	WebsiteURL  *string
	LinkedInURL *string
	GithubURL   *string
	ResumeURL   *string // CV file path
	CreatedAt   time.Time
	UpdatedAt   time.Time

	User User `gorm:"foreignKey:UserID"`
}

type WorkExperience struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;index"`
	Position    string
	CompanyName string
	StartDate   *time.Time
	EndDate     *time.Time
	Description *string `gorm:"type:text"`
}

type Education struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID         uuid.UUID `gorm:"type:uuid;index"`
	Degree         string
	Institution    string
	FieldOfStudy   string
	GraduationYear *int
}
