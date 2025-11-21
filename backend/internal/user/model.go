package user

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleCandidate UserRole = "candidate"
	RoleCompany   UserRole = "company"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name         string    `json:"name"`
	Email        string    `json:"email" gorm:"unique"`
	Password     *string   `gorm:"uniqueIndex"`
	PhoneNumber  *string
	AvatarURL    *string
	Role         *string `gorm:"type:varchar(20);default:candidate"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	RefreshToken []RefreshToken `gorm:"foreignKey:UserID"`
}

type RefreshToken struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Token     string `gorm:"uniqueIndex"`
	ExpiresAt time.Time
	UserID    string `gorm:"type:uuid"`
	CreatedAt time.Time

	User User `gorm:"foreignKey:UserID"`
}

type Project struct {
	ID        string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string
	Detail    string `gorm:"type:text"`
	CreatedBy string `gorm:"type:uuid"`
	CreatedAt time.Time

	Creator User `gorm:"foreignKey:CreatedBy"`
}
type SubmissionStatus string

const (
	StatusSubmitted SubmissionStatus = "submitted"
	StatusApproved  SubmissionStatus = "approved"
	StatusRejected  SubmissionStatus = "rejected"
)

type Submission struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProjectID   uuid.UUID `gorm:"type:uuid"`
	UserID      uuid.UUID `gorm:"type:uuid"`
	FileURL     string
	Status      SubmissionStatus `gorm:"type:varchar(20)"`
	SubmittedAt time.Time
	ApprovedAt  *time.Time

	Project Project `gorm:"foreignKey:ProjectID"`
	User    User    `gorm:"foreignKey:UserID"`
}

type Portfolio struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID       uuid.UUID `gorm:"type:uuid"`
	ProjectID    uuid.UUID `gorm:"type:uuid"`
	SubmissionID uuid.UUID `gorm:"type:uuid"`
	AddedAt      time.Time

	User       User       `gorm:"foreignKey:UserID"`
	Project    Project    `gorm:"foreignKey:ProjectID"`
	Submission Submission `gorm:"foreignKey:SubmissionID"`
}
