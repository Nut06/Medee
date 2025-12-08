// route in express
package server

import (
	authadapter "backend/internal/adapter/auth"
	useradapter "backend/internal/adapter/user"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Auth(app *fiber.App, db *gorm.DB) {
	api := app.Group("/")

	authGroup := api.Group("/auth")
	authRoute(authGroup, db)

	userGroup := api.Group("/user", AuthMiddleware)
	userRoute(userGroup, db)

	// skillGroup := api.Group("/skill", AuthMiddleware)
	// skillRoute(skillGroup, db)
}

// func skillRoute(router fiber.Router, db *gorm.DB) {
// 	h := skilladapter.NewHTTPHandler(db)
// 	router.Get("", h.GetSkills)
// 	router.Post("", h.AddSkill)
// 	router.Put("/:skillId", h.UpdateSkill)
// 	router.Delete("/:skillId", h.DeleteSkill)
// }

func authRoute(router fiber.Router, db *gorm.DB) {
	h := authadapter.NewHTTPHandler(db)
	router.Post("/register", h.Register)
	router.Post("/login", h.Login)
	router.Post("/refresh", h.Refresh)
	router.Post("/logout", h.Logout)
}

func userRoute(router fiber.Router, db *gorm.DB) {
	h := useradapter.NewHTTPHandler(db)
	router.Use(func(c *fiber.Ctx) error {
		fmt.Println("from IP", c.IP())
		fmt.Println("Method", c.Method())
		fmt.Println("Path", c.Path())
		return c.Next()
	})

	router.Get("", h.GetUser)
	router.Put("", h.UpdateProfile)
	router.Put("/avatar", h.UploadAvatar)
	router.Delete("/avatar", h.DeleteAvatar)

	// Candidate Features
	router.Get("/profile", h.GetFullProfile)
	router.Post("/experience", h.AddExperience)
	router.Put("/experience/:experienceId", h.UpdateExperience)
	router.Delete("/experience/:experienceId", h.DeleteExperience)
	router.Post("/education", h.AddEducation)
	router.Put("/education/:educationId", h.UpdateEducation)
	router.Delete("/education/:educationId", h.DeleteEducation)
	router.Put("/skills", h.UpdateSkills)
	router.Post("/project", h.AddProject)
	router.Put("/project/:projectId", h.UpdateProject)
	router.Delete("/project/:projectId", h.DeleteProject)
}
