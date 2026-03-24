package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	
	"portfolio-website-backend/internal/database"
)

func main() {
	// Connect to Database and run migrations
	database.ConnectDB()
	database.RunAutomigrations()

	// Initialize Gin router
	r := gin.Default()

	// Define health check endpoint
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	log.Println("Server is starting on port 8000...")
	
	// Start server
	if err := r.Run(":8000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
