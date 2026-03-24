package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"portfolio-website-backend/internal/config"
)

var DB *gorm.DB

func ConnectDB() {
	if config.Envs == nil {
		config.LoadConfig()
	}

	dbUrl := config.Envs.DatabaseURL
	if dbUrl == "" {
		log.Fatal("DATABASE_URL is not set in environment")
	}

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Example config for tuning
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
	}

	DB = db
	log.Println("Database connection established")
}
