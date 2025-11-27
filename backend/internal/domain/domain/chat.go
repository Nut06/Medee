package domain

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProjectID *uuid.UUID `gorm:"type:uuid;index"` // null ได้ ถ้าเป็น general chat
	// อาจจะผูกกับ application ด้วยก็ได้
	CreatedAt time.Time
}

type ConversationParticipant struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ConversationID uuid.UUID `gorm:"type:uuid;index"`
	UserID         uuid.UUID `gorm:"type:uuid;index"`
}

type Message struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ConversationID uuid.UUID `gorm:"type:uuid;index"`
	SenderID       uuid.UUID `gorm:"type:uuid;index"`
	Content        string    `gorm:"type:text"`
	CreatedAt      time.Time
	IsSystem       bool `gorm:"default:false"` // สำหรับ system message เช่น "Offer created"
}
