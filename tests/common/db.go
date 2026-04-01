package common

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"portfolio-website-backend/internal/database"
)

// SetupTestDB initializes a connection to a real Postgres instance specifically for testing.
// It uses TEST_DATABASE_URL if available, otherwise defaults to a local fallback container constraint.
func SetupTestDB() *gorm.DB {
	// Dynamically try to load .env depending on the directory the test was executed from
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load("../.env")
	_ = godotenv.Load(".env")

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}

	if dsn == "" {
		dsn = "postgresql://postgres:postgres@localhost:5432/portfolio_test?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to test database. Are you sure Postgres is running? Error: %v", err)
	}

	// Override global DB block and migrate fresh schema
	database.DB = db
	database.RunAutomigrations()

	return db
}
