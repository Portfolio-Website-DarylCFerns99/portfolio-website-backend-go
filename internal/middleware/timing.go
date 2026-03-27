package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// TimingMiddleware adds X-Process-Time header similar to the FastAPI middleware
func TimingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		processTime := time.Since(startTime).Seconds()
		c.Writer.Header().Set("X-Process-Time", strconv.FormatFloat(processTime, 'f', 6, 64))
	}
}
