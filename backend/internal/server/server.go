// controller in express
package server

import (
	"backend/internal/database"
	"backend/internal/domain/auth"
	"backend/internal/utils"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/idempotency"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func NewServer() *fiber.App {

	app := fiber.New(
		fiber.Config{
			ErrorHandler: func(c fiber.Ctx, err error) error {
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

	originsEnv := os.Getenv("CORS")
	allowedOrigins := strings.Split(originsEnv, ",")

	redis := database.NewRedis()
	sessionStore := session.NewStore(session.Config{
		Storage:         redis,
		CookieSecure:    os.Getenv("HTTPS") == "true", // HTTPS only
		CookieHTTPOnly:  true,                         // Prevent XSS
		CookieSameSite:  "Lax",                        // CSRF protection
		IdleTimeout:     utils.ThirtyMin,
		AbsoluteTimeout: utils.Oneday,
		Extractor:       extractors.FromCookie("__Host-session_id"),
	})

	locker := &database.RedisLocker{Redis: redis}
	app.Use(idempotency.New(idempotency.Config{
		Lifetime:  1 * time.Hour,
		KeyHeader: "X-Idempotency-Key",
		KeyHeaderValidate: func(k string) error {
			if len(k) != 36 {
				return fmt.Errorf("%w: invalid length", auth.ErrInvalidIdempotencyKey)
			}
			return nil
		},
		Storage:               redis,
		Lock:                  locker,
		KeepResponseHeaders:   []string{"Content-Type", "Location"},
		DisableValueRedaction: false,
	}))

	app.Use(limiter.New(limiter.Config{
		Next: func(c fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max: 20,
		MaxFunc: func(c fiber.Ctx) int {
			return 20
		},
		Expiration: utils.ThirtySec,
		ExpirationFunc: func(c fiber.Ctx) time.Duration {
			// Use longer expiration for sensitive endpoints
			if c.Path() == "/auth/login" {
				return utils.SixtySec
			}
			return utils.ThirtySec
		},
		// KeyGenerator:          func(c fiber.Ctx) string {
		// 	return c.Get("x-forwarded-for")
		// },

		// custom later
		// LimitReached: func(c fiber.Ctx) error {
		// 	return c.SendFile("./toofast.html")
		// },
		Storage: redis,
	}))

	csrfCookieName := "XSRF-TOKEN"

	if os.Getenv("ENV") == "production" {
		csrfCookieName = "__Host-csrf_"
	}

	app.Use(csrf.New(csrf.Config{
		TrustedOrigins:    allowedOrigins,
		CookieName:        csrfCookieName,
		CookieSecure:      os.Getenv("HTTPS") == "true",
		CookieHTTPOnly:    true, // false for SPAs
		CookieSameSite:    "Lax",
		CookieSessionOnly: true,
		Extractor:         extractors.FromHeader(csrfCookieName),
		Session:           sessionStore,
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowHeaders:     []string{"Origin, Content-type, Accept, Authorization"},
		AllowMethods:     []string{"GET, POST, PUT, PATCH, DELETE, OPTIONS"},
		AllowCredentials: true,
		MaxAge:           3600,
	}))

	app.Get("/", func(c fiber.Ctx) error {
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
