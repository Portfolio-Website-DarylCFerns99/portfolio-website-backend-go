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
		// Fallback to the primary dev database if a test specific URL isn't explicitly configured
		dsn = os.Getenv("DATABASE_URL")
	}

	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=portfolio port=5432 sslmode=disable"
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
