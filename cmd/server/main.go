package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"portfolio-website-backend/internal/database"
	"portfolio-website-backend/internal/handlers"
	"portfolio-website-backend/internal/middleware"
	"portfolio-website-backend/internal/repository"
	"portfolio-website-backend/internal/services"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Generated docs
	_ "portfolio-website-backend/docs"
)

// @title           Portfolio Website API
// @version         1.0
// @description     This is the backend API for my portfolio website.
// @host            localhost:8000
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Connect to Database
	db, err := database.GetDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Run migrations
	database.RunAutomigrations()

	// Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	experienceRepo := repository.NewExperienceRepository(db)

	// Initialize Services
	userService := services.NewUserService(db, userRepo)
	experienceService := services.NewExperienceService(experienceRepo)

	// Initialize Gin router
	r := gin.Default()

	// Apply Global Middlewares
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.TimingMiddleware())

	apiPrefix := os.Getenv("API_PREFIX")
	if apiPrefix == "" {
		apiPrefix = "/api/v1"
	}
	apiGroup := r.Group(apiPrefix)

	authMiddleware := middleware.RequireAuth()
	adminAuthMiddleware := middleware.RequireAdminAuth()

	// Initialize Handlers
	userHandler := handlers.NewUserHandler(userService, userRepo)
	experienceHandler := handlers.NewExperienceHandler(experienceService)

	// Register Routes
	userHandler.RegisterRoutes(apiGroup, authMiddleware, adminAuthMiddleware)
	experienceHandler.RegisterRoutes(apiGroup, authMiddleware)

	// Swagger API Docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Define health check endpoint
	r.GET("/healthz", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err == nil && sqlDB.Ping() == nil {
			c.Header("Cache-Control", "no-cache")
			c.Status(http.StatusOK)
		} else {
			c.Status(http.StatusServiceUnavailable)
		}
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Portfolio Website API"})
	})

	log.Println("Server is starting on port 8000...")

	// Start server
	if err := r.Run(":8000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
