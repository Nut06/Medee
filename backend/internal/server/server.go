// controller in express
package server

import (
	"backend/internal/database"
	"backend/internal/user"
	"log"
	"os"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("CORS"),
		AllowHeaders: "Origin, Content-type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowCredentials: true,
		MaxAge: 3600,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Use(logger.New())
	Auth(app)
	database.ConnectDB()
	user.InitUserModel()
	return  app
}