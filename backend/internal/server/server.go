// controller in express
package server

import (
	"backend/internal/database"
	"backend/internal/user"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)


func NewServer() *fiber.App{
	err := godotenv.Load()
	if err != nil {
		log.Print("No .env was found")
	}

	app := fiber.New(
		fiber.Config{
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				log.Printf("Error occur %v", err)

				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message":"Something went wrong",
				})
			},
			ReduceMemoryUsage: true,
			StrictRouting: true,
			CaseSensitive: true,
		},
	)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	// app.Use(logger.New())
	Auth(app)
	database.ConnectDB()
	user.InitUserModel()
	return  app
}