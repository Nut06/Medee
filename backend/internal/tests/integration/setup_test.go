package integration_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"backend/internal/database"
	"backend/internal/server"

	"github.com/gofiber/fiber/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

var app *fiber.App

func TestMain(m *testing.M) {
	ctx := context.Background()

	// 1. Spin up PostgreSQL Container
	// NOTE: testcontainers modules MUST use the Official Docker Hub image.
	// Custom images (dhi.io/*) have different entrypoints and are incompatible.
	pgContainer, err := postgres.Run(ctx,
		"postgres:18",
		postgres.WithDatabase("medee_test"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(15*time.Second),
		),
	)
	if err != nil {
		log.Fatalf("failed to start postgres container: %s", err)
	}
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate postgres container: %s", err)
		}
	}()

	pgHost, err := pgContainer.Host(ctx)
	if err != nil {
		log.Fatalf("failed to get pg host: %v", err)
	}
	pgPort, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("failed to get pg port: %v", err)
	}

	// 2. Spin up Redis Container
	// NOTE: Same requirement - must use Official Docker Hub image.
	redisContainer, err := redis.Run(ctx,
		"redis:8.6",
	)
	if err != nil {
		log.Fatalf("failed to start redis container: %s", err)
	}
	defer func() {
		if err := redisContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate redis container: %s", err)
		}
	}()

	redisURL, err := redisContainer.Endpoint(ctx, "")
	if err != nil {
		log.Fatalf("failed to get redis endpoint: %v", err)
	}

	// 3. Override Environment Variables for the test execution
	os.Setenv("DB_HOST", pgHost)
	os.Setenv("DB_PORT", pgPort.Port())
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "medee_test")
	os.Setenv("REDIS_URL", redisURL)
	os.Setenv("JWT_SECRET", "supersecretkey_for_tests")
	os.Setenv("CORS", "http://localhost:5173")

	// 4. Initialize the Server with dynamic configurations
	app = server.NewServer()

	// 5. Run all integration tests in the package
	code := m.Run()

	os.Exit(code)
}

func clearUsersTable() {
	if database.DB != nil {
		// Clean both users and their related refresh tokens 
		// Using CASCADE to ensure foreign key dependencies are wiped
		database.DB.Exec("TRUNCATE TABLE refresh_tokens, users CASCADE;")
	}
}

func clearRefreshTokens() {
	if database.DB != nil {
		database.DB.Exec("TRUNCATE TABLE refresh_tokens;")
	}
}
