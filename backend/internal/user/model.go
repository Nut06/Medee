package user

import "time"

type User struct {
	ID           string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name         string  `json:"name"`
	Email        string  `json:"email" gorm:"unique"`
	Password     *string `gorm:"uniqueIndex"`
	PhoneNumber  *string
	AvatarURL    *string
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
}