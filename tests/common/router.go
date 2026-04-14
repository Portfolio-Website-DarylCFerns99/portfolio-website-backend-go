package common

import (
	"portfolio-website-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// SetupRouter creates a flexible testing router for standard endpoints
func SetupRouter(testUser *models.User, registerRoutes func(r *gin.RouterGroup, auth gin.HandlerFunc)) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	authMiddleware := func(c *gin.Context) {
		if testUser != nil {
			c.Set("current_user", testUser)
			c.Set("user_id", testUser.ID)
		}
		c.Next()
	}

	registerRoutes(r.Group("/"), authMiddleware)
	return r
}

// SetupRouterWithAdmin creates a flexible testing router for handlers requiring admin middleware
func SetupRouterWithAdmin(testUser *models.User, customAuthMiddleware gin.HandlerFunc, customAdminAuthMiddleware gin.HandlerFunc, registerRoutes func(r *gin.RouterGroup, auth gin.HandlerFunc, admin gin.HandlerFunc)) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	if customAuthMiddleware == nil {
		customAuthMiddleware = func(c *gin.Context) {
			if testUser != nil {
				c.Set("current_user", testUser)
				c.Set("user_id", testUser.ID)
			}
			c.Next()
		}
	}

	if customAdminAuthMiddleware == nil {
		customAdminAuthMiddleware = func(c *gin.Context) {
			c.Next()
		}
	}

	registerRoutes(r.Group("/"), customAuthMiddleware, customAdminAuthMiddleware)
	return r
}
