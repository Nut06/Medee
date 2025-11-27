package domain

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotifNewApplication NotificationType = "new_application"
	NotifStatusChange   NotificationType = "status_change"
	NotifNewMessage     NotificationType = "new_message"
	NotifNewProject     NotificationType = "new_project"
)

type Notification struct {
	ID        uuid.UUID        `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID        `gorm:"type:uuid;index"`
	Type      NotificationType `gorm:"type:varchar(30)"`
	Title     string
	Body      string  `gorm:"type:text"`
	IsRead    bool    `gorm:"default:false"`
	LinkURL   *string // deep link ไปหน้าใน frontend
	CreatedAt time.Time
}
