package domain

import (
	"time"

	"github.com/google/uuid"
)

type ApplicantProfile struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
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
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;index"`
	Position    string
	CompanyName string
	StartDate   *time.Time
	EndDate     *time.Time
	Description *string `gorm:"type:text"`
}

type Education struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID         uuid.UUID `gorm:"type:uuid;index"`
	InstituteID    uuid.UUID `gorm:"type:uuid;index"`
	FieldOfStudyID uuid.UUID `gorm:"type:uuid;index"`
	StartDate      *time.Time
	EndDate        *time.Time
	GPA            *float64
	Description    *string `gorm:"type:text"`
	Degree         string
	GraduationYear *int

	// ✅ Transient fields (read-only) - populated by repository SELECT queries
	InstituteName    string `gorm:"->"`
	FieldOfStudyName string `gorm:"->"`
}

type Institute struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"not null;index:idx_institute_name"`
	Country   string    `gorm:"index:idx_institute_country"`
	Domain    *string   `gorm:"uniqueIndex"`
	WebPages  []string  `gorm:"type:jsonb"`
	Source    string    `gorm:"default:'manual'"` // "manual", "hipo_api", "cache"
	CreatedAt time.Time
	UpdatedAt time.Time
}

type FieldOfStudy struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string
	Code      string `gorm:"uniqueIndex"` // CIP code
	Title     string
	Category  string // Broad category
	Level     string // "broad", "intermediate", "detailed"
	IsStem    bool   `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
