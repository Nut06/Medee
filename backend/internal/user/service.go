package user

import (
	"backend/internal/database"
)
func CreateUser( user *User) error{
	return database.DB.Create(user).Error
}