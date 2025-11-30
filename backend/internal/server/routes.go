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

	userGroup := api.Group("/user", AuthMiddleware)
	userRoute(userGroup, db)
}

func authRoute(router fiber.Router, db *gorm.DB) {
	h := authadapter.NewHTTPHandler(db)
	router.Post("/register", h.Register)
	router.Post("/login", h.Login)
	router.Post("/refresh", h.Refresh)
	router.Post("/logout", h.Logout)
}

func userRoute(router fiber.Router, db *gorm.DB) {
	h := useradapter.NewHTTPHandler(db)
	router.Get("/user", h.GetUser)
	router.Put("/user", h.UpdateProfile)
	router.Put("/user/avatar", h.UploadAvatar)
	router.Delete("/user/avatar", h.DeleteAvatar)

	// Candidate Features
	router.Get("/user/profile", h.GetFullProfile)
	router.Post("/user/experience", h.AddExperience)
	router.Put("/user/experience/:experienceId", h.UpdateExperience)
	router.Delete("/user/experience/:experienceId", h.DeleteExperience)
	router.Post("/user/education", h.AddEducation)
	router.Put("/user/education/:educationId", h.UpdateEducation)
	router.Delete("/user/education/:educationId", h.DeleteEducation)
	router.Put("/user/skills", h.UpdateSkills)
	router.Post("/user/project", h.AddProject)
	router.Put("/user/project/:projectId", h.UpdateProject)
	router.Delete("/user/project/:projectId", h.DeleteProject)
}
