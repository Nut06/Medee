// route in express
package server

import (
	"backend/internal/auth"

	"github.com/gofiber/fiber/v2"
)

func Auth(app *fiber.App){
	api := app.Group("/")

	authGroup := api.Group("/auth")
	authRoute(authGroup)
	userGroup := api.Group("/user")
	userRoute(userGroup)
}

func authRoute(router fiber.Router){
	router.Post("/register", auth.RegisterHandler)
	router.Post("/login", auth.LoginHandler)
	router.Post("/logout", auth.LogoutHandler)
}

func userRoute(router fiber.Router){
	// router.Get()
}