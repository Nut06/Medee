package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProjectType string

const (
	ProjectTypeInternship ProjectType = "internship"
	ProjectTypeFullTime   ProjectType = "full_time"
	ProjectTypeContract   ProjectType = "contract"
	ProjectTypeTestJob    ProjectType = "test_job"
)

type ProjectStatus string

const (
	ProjectStatusDraft  ProjectStatus = "draft"
	ProjectStatusOpen   ProjectStatus = "open"
	ProjectStatusClosed ProjectStatus = "closed"
)

type Project struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CompanyID     uuid.UUID `gorm:"type:uuid;index"`
	CreatedByUser uuid.UUID `gorm:"type:uuid;index"`
	Title         string
	Description   string      `gorm:"type:text"` // rich text
	Type          ProjectType `gorm:"type:varchar(30)"`
	Location      *string
	MinExperience *int // years
	EducationReq  *string
	Status        ProjectStatus `gorm:"type:varchar(20);default:'open'"`
	ApplyDeadline *time.Time
	Duration      *string // e.g. "3 months"
	IsTemplate    bool    `gorm:"default:false"`
	TemplateName  *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ProjectSkill struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProjectID uuid.UUID `gorm:"type:uuid;index"`
	SkillID   uuid.UUID `gorm:"type:uuid;index"`
}
