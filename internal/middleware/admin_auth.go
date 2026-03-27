package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// RequireAdminAuth ensures the request has the correct X-Admin-Api-Key header
func RequireAdminAuth() gin.HandlerFunc {
	// Cache the expected key at startup. Wait, os.Getenv every time is fine and safer if testing sets env dynamically.
	return func(c *gin.Context) {
		expectedKey := os.Getenv("ADMIN_API_KEY")
		if expectedKey == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": "Admin API key not configured on server"})
			return
		}

		providedKey := c.GetHeader("X-Admin-Api-Key")
		if providedKey == "" || providedKey != expectedKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "Invalid or missing Admin API Key"})
			return
		}

		c.Next()
	}
}
