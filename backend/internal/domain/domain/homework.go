package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProjectTaskStatus string

const (
	TaskTodo       ProjectTaskStatus = "todo"
	TaskInProgress ProjectTaskStatus = "in_progress"
	TaskDone       ProjectTaskStatus = "done"
)

type ProjectTask struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProjectID   uuid.UUID  `gorm:"type:uuid;index"`
	AssigneeID  *uuid.UUID `gorm:"type:uuid;index"` // candidate user
	Title       string
	Description *string           `gorm:"type:text"`
	Status      ProjectTaskStatus `gorm:"type:varchar(20);default:'todo'"`
	DueDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
