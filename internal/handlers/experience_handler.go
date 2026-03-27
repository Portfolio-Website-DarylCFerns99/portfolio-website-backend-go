package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/services"
)

type ExperienceHandler struct {
	service services.ExperienceService
}

func NewExperienceHandler(service services.ExperienceService) *ExperienceHandler {
	return &ExperienceHandler{service: service}
}

func (h *ExperienceHandler) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	experiences := r.Group("/experiences")

	// Public routes
	experiences.GET("/public/:user_id", h.GetPublicExperiences)

	// Protected routes
	protected := experiences.Group("")
	protected.Use(authMiddleware)
	{
		protected.POST("", h.CreateExperience)
		protected.GET("", h.GetExperiences)
		protected.GET("/:experience_id", h.GetExperience)
		protected.PUT("/:experience_id", h.UpdateExperience)
		protected.PATCH("/:experience_id/visibility", h.UpdateVisibility)
		protected.DELETE("/:experience_id", h.DeleteExperience)
	}
}

func (h *ExperienceHandler) parseQueryArgs(c *gin.Context) (skip, limit int, expType string) {
	skip, _ = strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "100"))
	expType = c.Query("type")
	return
}

func (h *ExperienceHandler) getUserID(c *gin.Context) (uuid.UUID, bool) {
	userObj, exists := c.Get("current_user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "Not authenticated"})
		return uuid.Nil, false
	}
	user, ok := userObj.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "User context invalid"})
		return uuid.Nil, false
	}
	return user.ID, true
}

// CreateExperience
// @Summary      Create experience
// @Description  Create a new experience or education entry
// @Tags         experiences
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      models.Experience  true  "Experience Data"
// @Success      201      {object}  models.Experience
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /experiences [post]
func (h *ExperienceHandler) CreateExperience(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	var exp models.Experience
	if err := c.ShouldBindJSON(&exp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	exp.UserID = userID

	created, err := h.service.CreateExperience(&exp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to create experience"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetExperiences
// @Summary      Get experiences
// @Description  Get a paginated list of experiences (requires auth)
// @Tags         experiences
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        skip   query    int  false  "Skip count"
// @Param        limit  query    int  false  "Limit count"
// @Param        type   query    string false "Type (experience/education)"
// @Success      200    {object} map[string]interface{}
// @Failure      400    {object} map[string]interface{}
// @Failure      500    {object} map[string]interface{}
// @Router       /experiences [get]
func (h *ExperienceHandler) GetExperiences(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	skip, limit, expType := h.parseQueryArgs(c)

	var items []models.Experience
	var total int64
	var err error

	if expType != "" {
		if expType != "experience" && expType != "education" {
			c.JSON(http.StatusBadRequest, gin.H{"detail": "Type must be either 'experience' or 'education'"})
			return
		}
		items, total, err = h.service.GetExperiencesByType(userID, expType, skip, limit, false)
	} else {
		items, total, err = h.service.GetExperiences(userID, skip, limit, false)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to fetch experiences"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"experiences": items,
		"total":       total,
	})
}

// GetPublicExperiences
// @Summary      Get public experiences
// @Description  Get a paginated list of public experiences
// @Tags         experiences
// @Accept       json
// @Produce      json
// @Param        user_id path     string  true "User ID"
// @Param        skip   query    int  false  "Skip count"
// @Param        limit  query    int  false  "Limit count"
// @Param        type   query    string false "Type (experience/education)"
// @Success      200    {object} map[string]interface{}
// @Failure      400    {object} map[string]interface{}
// @Failure      500    {object} map[string]interface{}
// @Router       /experiences/public/{user_id} [get]
func (h *ExperienceHandler) GetPublicExperiences(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid user ID"})
		return
	}

	skip, limit, expType := h.parseQueryArgs(c)

	var items []models.Experience
	var total int64

	if expType != "" {
		if expType != "experience" && expType != "education" {
			c.JSON(http.StatusBadRequest, gin.H{"detail": "Type must be either 'experience' or 'education'"})
			return
		}
		items, total, err = h.service.GetExperiencesByType(userID, expType, skip, limit, true)
	} else {
		items, total, err = h.service.GetExperiences(userID, skip, limit, true)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to fetch public experiences"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"experiences": items,
		"total":       total,
	})
}

// GetExperience
// @Summary      Get experience by ID
// @Description  Get a single experience by its ID
// @Tags         experiences
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        experience_id  path      string  true  "Experience ID"
// @Success      200    {object} models.Experience
// @Failure      400    {object} map[string]interface{}
// @Failure      404    {object} map[string]interface{}
// @Router       /experiences/{experience_id} [get]
func (h *ExperienceHandler) GetExperience(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("experience_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	exp, err := h.service.GetExperienceByID(userID, id, false)
	if err != nil || exp == nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Experience not found"})
		return
	}

	c.JSON(http.StatusOK, exp)
}

// UpdateExperience
// @Summary      Update experience
// @Description  Update an existing experience
// @Tags         experiences
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        experience_id  path      string  true  "Experience ID"
// @Param        request        body      map[string]interface{} true "Update Data"
// @Success      200    {object} models.Experience
// @Failure      400    {object} map[string]interface{}
// @Failure      404    {object} map[string]interface{}
// @Router       /experiences/{experience_id} [put]
func (h *ExperienceHandler) UpdateExperience(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("experience_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid request body"})
		return
	}

	updated, err := h.service.UpdateExperience(userID, id, updateData)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Experience not found or update failed"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

type visibilityReq struct {
	IsVisible bool `json:"is_visible"`
}

// UpdateVisibility
// @Summary      Update experience visibility
// @Description  Change the visibility of an experience
// @Tags         experiences
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        experience_id  path      string  true  "Experience ID"
// @Param        request        body      visibilityReq true "Visibility Update"
// @Success      200    {object} models.Experience
// @Failure      400    {object} map[string]interface{}
// @Failure      404    {object} map[string]interface{}
// @Router       /experiences/{experience_id}/visibility [patch]
func (h *ExperienceHandler) UpdateVisibility(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("experience_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	var req visibilityReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid body"})
		return
	}

	updated, err := h.service.UpdateExperienceVisibility(userID, id, req.IsVisible)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Experience not found"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteExperience
// @Summary      Delete experience
// @Description  Delete an experience
// @Tags         experiences
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        experience_id  path      string  true  "Experience ID"
// @Success      204
// @Failure      400    {object} map[string]interface{}
// @Failure      404    {object} map[string]interface{}
// @Router       /experiences/{experience_id} [delete]
func (h *ExperienceHandler) DeleteExperience(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("experience_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	err = h.service.DeleteExperience(userID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Experience not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
