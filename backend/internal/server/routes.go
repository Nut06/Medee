// route in express
package server

import (
	authadapter "backend/internal/adapter/auth"
	field_of_study_adapter "backend/internal/adapter/field_of_study"
	institute_adapter "backend/internal/adapter/institute"
	skilladapter "backend/internal/adapter/skill"
	useradapter "backend/internal/adapter/user"
	storage "backend/internal/adapter/storage"
	"backend/internal/application/fieldapp"
	"backend/internal/application/instituteapp"
	"backend/internal/database"
	"github.com/gofiber/storage/redis/v3"
	"backend/internal/utils"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Auth(app *fiber.App, db *gorm.DB) {
	api := app.Group("/")

	// Initialize Redis, and supabase
	redisClient := database.NewRedis()
	supabase := storage.NewSupabaseStorage()
	jwtService := authadapter.NewJWTService(os.Getenv("JWT_SECRET"), utils.FifteenMin, utils.SevenDays)

	authGroup := api.Group("/auth")
	authRoute(authGroup, db, jwtService)

	userGroup := api.Group("/user", AuthMiddleware(jwtService))
	userRoute(userGroup, db)

	skillGroup := api.Group("/skill", AuthMiddleware(jwtService))
	skillRoute(skillGroup, db)

	// Institute routes (public for autocomplete)
	instituteGroup := api.Group("/institutes")
	instituteRoute(instituteGroup, db, redisClient)

	// FieldOfStudy routes (public for autocomplete)
	fieldGroup := api.Group("/field-of-studies")
	fieldOfStudyRoute(fieldGroup, db)

	uploadGroup := api.Group("/upload", AuthMiddleware(jwtService))

	uploadRoute(uploadGroup, supabase)
}

func uploadRoute(router fiber.Router, supabase *storage.SupabaseStorage){
	storageHandler := storage.NewStorageHandler(supabase)
	router.Post("/signed-url", storageHandler.GetSignedUploadURL)
}

func instituteRoute(router fiber.Router, db *gorm.DB, redisClient *redis.Storage) {
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

func authRoute(router fiber.Router, db *gorm.DB, jwtService *authadapter.JWTService) {
	h := authadapter.NewHTTPHandler(db, jwtService)
	router.Post("/register", h.Register)
	router.Post("/login", h.Login)
	router.Post("/refresh", h.Refresh)
	router.Post("/logout", h.Logout)
	// router.Post("/logout", AuthMiddleware, h.LogoutWithUserId)
}

func userRoute(user fiber.Router, db *gorm.DB) {
	h := useradapter.NewHTTPHandler(db)
	user.Use(func(c fiber.Ctx) error {
		fmt.Println("from IP", c.IP())
		fmt.Println("Method", c.Method())
		fmt.Println("Path", c.Path())
		return c.Next()
	})

	user.Get("", h.GetUser)
	user.Put("", h.UpdateProfile)
	user.Put("/avatar", h.UploadAvatar)
	user.Delete("/avatar", h.DeleteAvatar)
	user.Put("/resume", h.UploadResume)
	user.Delete("/resume", h.DeleteResume)

	// Candidate Features
	user.Get("/profile", h.GetFullProfile)
	user.Post("/experience", h.AddExperience)
	user.Put("/experience/:experienceId", h.UpdateExperience)
	user.Delete("/experience/:experienceId", h.DeleteExperience)
	user.Post("/education", h.AddEducation)
	user.Put("/education/:educationId", h.UpdateEducation)
	user.Delete("/education/:educationId", h.DeleteEducation)

	// Skills (User)
	user.Put("/skills", h.UpdateSkills) // PUT /user/skills (bulk replace)
	user.Post("/skills", h.AddSkill)    // POST /user/skills (add single)
	user.Delete("/skills/:skillId", h.DeleteSkill)

	user.Post("/project", h.AddProject)
	user.Put("/project/:projectId", h.UpdateProject)
	user.Delete("/project/:projectId", h.DeleteProject)
}
