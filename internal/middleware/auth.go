package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"portfolio-website-backend/internal/database"
	"portfolio-website-backend/internal/repository"
	"portfolio-website-backend/internal/security"
)

// Configured dynamically during server setup
var UserRepo repository.UserRepository

// RequireAuth middleware verifies the JWT token and injects the user into the context
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "Not authenticated"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "Invalid authorization header format"})
			return
		}

		tokenString := parts[1]
		claims, err := security.DecodeToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "Could not validate credentials"})
			return
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "Could not validate credentials"})
			return
		}

		userID, err := uuid.Parse(sub)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "Could not validate credentials"})
			return
		}

		// Ensure UserRepo is initialized (fallback)
		if UserRepo == nil {
			db, _ := database.GetDB()
			UserRepo = repository.NewUserRepository(db)
		}

		user, err := UserRepo.GetByID(userID)
		if err != nil || user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "User not found"})
			return
		}

		// Set current_user and user_id in Gin context
		c.Set("current_user", user)
		c.Set("user_id", userID)
		c.Next()
	}
}
