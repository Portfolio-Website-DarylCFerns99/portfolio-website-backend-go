package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSMiddleware simple CORS configuration mapping FastAPI's defaults
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// In a real production app, read from settings.CORS_ORIGINS
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
