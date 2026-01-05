// route in express
package server

import (
	authadapter "backend/internal/adapter/auth"
	field_of_study_adapter "backend/internal/adapter/field_of_study"
	institute_adapter "backend/internal/adapter/institute"
	skilladapter "backend/internal/adapter/skill"
	useradapter "backend/internal/adapter/user"
	"backend/internal/application/fieldapp"
	"backend/internal/application/instituteapp"
	"backend/internal/database"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Auth(app *fiber.App, db *gorm.DB) {
	api := app.Group("/")

	// Initialize Redis
	redisClient := database.NewRedisClient()

	authGroup := api.Group("/auth")
	authRoute(authGroup, db)

	userGroup := api.Group("/user", AuthMiddleware)
	userRoute(userGroup, db)

	skillGroup := api.Group("/skill", AuthMiddleware)
	skillRoute(skillGroup, db)

	// Institute routes (public for autocomplete)
	instituteGroup := api.Group("/institutes")
	instituteRoute(instituteGroup, db, redisClient)

	// FieldOfStudy routes (public for autocomplete)
	fieldGroup := api.Group("/field-of-studies")
	fieldOfStudyRoute(fieldGroup, db)
}

func instituteRoute(router fiber.Router, db *gorm.DB, redisClient *redis.Client) {
	repo := institute_adapter.NewInstituteRepository(db, redisClient)
	hipoClient := institute_adapter.NewHipoAPIClient()
	usecase := instituteapp.NewUsecase(repo, hipoClient)
	h := institute_adapter.NewHTTPHandler(usecase)

	router.Get("/search", h.SearchInstitutes) // GET /institutes/search?q=...&country=...
	router.Get("/:id", h.GetInstitute)        // GET /institutes/:id
}

func fieldOfStudyRoute(router fiber.Router, db *gorm.DB) {
	repo := field_of_study_adapter.NewFieldOfStudyRepository(db)
	usecase := fieldapp.NewUsecase(repo)
	h := field_of_study_adapter.NewHTTPHandler(usecase)

	router.Get("/search", h.SearchFieldOfStudies) // GET /field-of-studies/search?q=...
	router.Get("/:id", h.GetFieldOfStudy)         // GET /field-of-studies/:id
}

func skillRoute(router fiber.Router, db *gorm.DB) {
	h := skilladapter.NewHTTPHandler(db)
	router.Get("", h.GetSkills)    // GET /skill?q=...
	router.Post("", h.CreateSkill) // POST /skill (Create Master Skill - Optional)
}

func authRoute(router fiber.Router, db *gorm.DB) {
	h := authadapter.NewHTTPHandler(db)
	router.Post("/register", h.Register)
	router.Post("/login", h.Login)
	router.Post("/refresh", h.Refresh)
	router.Post("/logout", h.Logout)
	// router.Post("/logout", AuthMiddleware, h.LogoutWithUserId)
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
	router.Put("/resume", h.UploadResume)
	router.Delete("/resume", h.DeleteResume)

	// Candidate Features
	router.Get("/profile", h.GetFullProfile)
	router.Post("/experience", h.AddExperience)
	router.Put("/experience/:experienceId", h.UpdateExperience)
	router.Delete("/experience/:experienceId", h.DeleteExperience)
	router.Post("/education", h.AddEducation)
	router.Put("/education/:educationId", h.UpdateEducation)
	router.Delete("/education/:educationId", h.DeleteEducation)

	// Skills (User)
	router.Post("/skills", h.AddSkill)
	router.Delete("/skills/:skillId", h.DeleteSkill)
	// router.Put("/skills", h.UpdateSkills) // Deprecated in favor of Atomic Add/Delete

	router.Post("/project", h.AddProject)
	router.Put("/project/:projectId", h.UpdateProject)
	router.Delete("/project/:projectId", h.DeleteProject)
}
