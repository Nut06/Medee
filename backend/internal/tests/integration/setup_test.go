package integration_test

import (
	"log"
	"os"
	"testing"

	"backend/internal/database"
	"backend/internal/server"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var app *fiber.App

func TestMain(m *testing.M) {
	// Try loading .env file if running locally, ignore if it fails (e.g., inside docker-compose)
	_ = godotenv.Load("../../../../backend/.env")
	_ = godotenv.Load("../../../../.env")

	// Ensure tests are running against the test database for safety
	if os.Getenv("DB_NAME") != "medee_test" && os.Getenv("CI") != "true" {
		log.Println("WARNING: DB_NAME is not medee_test. Make sure you are running tests in the correct environment.")
	}

	app = server.NewServer()

	code := m.Run()

	os.Exit(code)
}

func clearUsersTable() {
	if database.DB != nil {
		database.DB.Exec("TRUNCATE TABLE users CASCADE;")
	}
}
