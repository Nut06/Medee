// route in express
package server

import (
	authadapter "backend/internal/adapter/auth"
	useradapter "backend/internal/adapter/user"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Auth(app *fiber.App, db *gorm.DB) {
	api := app.Group("/")

	authGroup := api.Group("/auth")
	authRoute(authGroup, db)
	// userGroup := api.Group("/user")
	// userRoute(userGroup)
}

func authRoute(router fiber.Router, db *gorm.DB) {
	h := authadapter.NewHTTPHandler(db)
	router.Post("/register", h.Register)
	router.Post("/login", h.Login)
	router.Post("/logout", h.Logout)
}

func userRoute(router fiber.Router, db *gorm.DB) {
	h := useradapter.NewHTTPHandler(db)
	router.Get("/user", h.GetUser)
	router.Put("/user", h.UpdateProfile)
	router.Put("/user/avatar", h.UploadAvatar)
}


