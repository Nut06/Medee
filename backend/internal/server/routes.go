// route in express
package server

import (
	authadapter "backend/internal/adapter/auth"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

func Auth(app *fiber.App, db *gorm.DB) {
	api := app.Group("/")

	authGroup := api.Group("/auth")
	authRoute(authGroup, db)
	userGroup := api.Group("/user")
	userRoute(userGroup)
}

func (s *FiberServer) RegisterRoute() {
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CORS"),
		AllowHeaders:     "Origin, Content-type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowCredentials: true,
		MaxAge:           300,
	}))
}

func authRoute(router fiber.Router, db *gorm.DB) {
	h := authadapter.NewHTTPHandler(db)
	router.Post("/register", h.Register)
	router.Post("/login", h.Login)
	router.Post("/logout", h.Logout)
}

func userRoute(router fiber.Router) {
	// router.Get()
}


