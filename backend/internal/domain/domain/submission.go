package domain

import (
	"time"

	"github.com/google/uuid"
)

type ApplicationStatus string

const (
	AppStatusApplied     ApplicationStatus = "applied"
	AppStatusUnderReview ApplicationStatus = "under_review"
	AppStatusInterview   ApplicationStatus = "interview"
	AppStatusOffered     ApplicationStatus = "offered"
	AppStatusRejected    ApplicationStatus = "rejected"
	AppStatusHired       ApplicationStatus = "hired"
)

type Application struct {
	ID          uuid.UUID         `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProjectID   uuid.UUID         `gorm:"type:uuid;index"`
	CandidateID uuid.UUID         `gorm:"type:uuid;index"` // UserID (candidate)
	CoverLetter *string           `gorm:"type:text"`
	Status      ApplicationStatus `gorm:"type:varchar(30);default:'applied'"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SubmissionStatus string

const (
    StatusSubmitted SubmissionStatus = "submitted"
    StatusApproved  SubmissionStatus = "approved"
    StatusRejected  SubmissionStatus = "rejected"
)

type Submission struct {
    ID            uuid.UUID        `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    ApplicationID uuid.UUID        `gorm:"type:uuid;index"`
    FileURL       string
    ExtraLink     *string          // GitHub, Video, etc.
    Description   *string          `gorm:"type:text"`
    Status        SubmissionStatus `gorm:"type:varchar(20)"`
    SubmittedAt   time.Time
    ReviewedAt    *time.Time

    Application Application `gorm:"foreignKey:ApplicationID"`
}