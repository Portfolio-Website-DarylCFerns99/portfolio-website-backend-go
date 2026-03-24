package database

import (
	"log"

	"portfolio-website-backend/internal/models"
)

// RunAutomigrations uses GORM to automatically migrate schemas
func RunAutomigrations() {
	if DB == nil {
		log.Fatal("Database connection not established. Call ConnectDB first.")
	}

	log.Println("Running AutoMigrations...")

	// Create the vector extension if it doesn't exist so pgvector fields work
	if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS vector;").Error; err != nil {
		log.Printf("Failed to create vector extension: %v", err)
	}

	err := DB.AutoMigrate(
		&models.User{},
		&models.ProjectCategory{},
		&models.Project{},
		&models.Experience{},
		&models.SkillGroup{},
		&models.Skill{},
		&models.Review{},
		&models.ChatSession{},
		&models.ChatMessage{},
		&models.VectorEmbedding{},
	)

	if err != nil {
		log.Fatalf("Failed to run AutoMigrate: %v", err)
	}

	log.Println("AutoMigrations completed successfully!")
}
