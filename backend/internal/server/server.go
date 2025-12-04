// controller in express
package server

import (
	"backend/internal/database"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func NewServer() *fiber.App {

	app := fiber.New(
		fiber.Config{
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				log.Printf("Error occur %v", err)
				code := fiber.StatusInternalServerError
				if e, ok := err.(*fiber.Error); ok {
					code = e.Code
				}
				return c.Status(code).JSON(fiber.Map{
					"message": err.Error(),
				})
			},
			ReduceMemoryUsage: true,
			StrictRouting:     true,
			CaseSensitive:     true,
		},
	)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CORS"),
		AllowHeaders:     "Origin, Content-type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowCredentials: true,
		MaxAge:           3600,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Use(logger.New())
	db := database.ConnectDB()
	database.AutoMigrate(db)
	Auth(app, db)

	return app
}

type FiberServer struct {
	*fiber.App
	db database.Service
}

// func New() *FiberServer {
// 	server := &FiberServer{
// 		App: fiber.New(
// 				fiber.Config{
// 				ErrorHandler: func(c *fiber.Ctx, err error) error {
// 					log.Printf("Error occur %v", err)

// 					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 						"message":"Something went wrong",
// 					})
// 				},
// 				ReduceMemoryUsage: true,
// 				StrictRouting: true,
// 				CaseSensitive: true,
// 			}),
// 		db: database.New(),
// 	}
// 	return  server
// }
