package user

import (
	"github.com/gofiber/fiber/v2"
	"backend/internal/database"
)

func InitUserModel(){
	database.DB.AutoMigrate(
		&User{},
		&RefreshToken{},
	)
}

func HandleHome(c *fiber.Ctx) error{
	return c.SendString("Hello from server")
}


func GetUserName(c *fiber.Ctx) error{
	name := c.Params("name")
	return c.SendString("Hello "+name)
}

func GetUsers(c *fiber.Ctx) error {
	var users []User
	database.DB.Find(&users)

	return c.Status(fiber.StatusOK).JSON(users)
}