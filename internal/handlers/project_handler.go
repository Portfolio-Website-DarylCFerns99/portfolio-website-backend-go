package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/services"
)

type ProjectHandler struct {
	service services.ProjectService
}

func NewProjectHandler(service services.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

func (h *ProjectHandler) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	projects := r.Group("/projects")

	// Public routes
	projects.GET("/public/:user_id", h.GetPublicProjects)

	// Protected routes
	protected := projects.Group("")
	protected.Use(authMiddleware)
	{
		protected.POST("", h.CreateProject)
		protected.GET("", h.GetProjects)
		protected.GET("/:project_id", h.GetProject)
		protected.PUT("/:project_id", h.UpdateProject)
		protected.PATCH("/:project_id/visibility", h.UpdateVisibility)
		protected.DELETE("/:project_id", h.DeleteProject)
	}
}

func (h *ProjectHandler) parseQueryArgs(c *gin.Context) (skip, limit int) {
	skip, _ = strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "100"))
	return
}

func (h *ProjectHandler) getUserID(c *gin.Context) (uuid.UUID, bool) {
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

type ProjectCreateReq struct {
	Title             string     `json:"title" binding:"required"`
	Description       *string    `json:"description"`
	Type              string     `json:"type" binding:"required,oneof=github custom"`
	Image             *string    `json:"image"`
	Tags              []string   `json:"tags"`
	URL               *string    `json:"url"`
	IsVisible         bool       `json:"is_visible"`
	ProjectCategoryID *uuid.UUID `json:"project_category_id"`
}

// CreateProject
// @Summary      Create project
// @Description  Create a new project entry
// @Tags         projects
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      ProjectCreateReq  true  "Project Data"
// @Success      201      {object}  models.Project
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	var req ProjectCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	exp := models.Project{
		UserID:            userID,
		Title:             req.Title,
		Description:       req.Description,
		Type:              req.Type,
		Image:             req.Image,
		Tags:              req.Tags,
		URL:               req.URL,
		IsVisible:         req.IsVisible,
		ProjectCategoryID: req.ProjectCategoryID,
	}

	created, err := h.service.CreateProject(&exp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetProjects
// @Summary      Get projects
// @Description  Get a paginated list of projects (requires auth)
// @Tags         projects
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        skip   query    int  false  "Skip count"
// @Param        limit  query    int  false  "Limit count"
// @Success      200    {object} map[string]interface{}
// @Failure      400    {object} map[string]interface{}
// @Failure      500    {object} map[string]interface{}
// @Router       /projects [get]
func (h *ProjectHandler) GetProjects(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	skip, limit := h.parseQueryArgs(c)

	items, total, err := h.service.GetProjects(userID, skip, limit, false)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to fetch projects"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"projects": items,
		"total":    total,
	})
}

// GetPublicProjects
// @Summary      Get public projects
// @Description  Get a paginated list of public projects
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        user_id path     string  true "User ID"
// @Param        skip   query    int  false  "Skip count"
// @Param        limit  query    int  false  "Limit count"
// @Success      200    {object} map[string]interface{}
// @Failure      400    {object} map[string]interface{}
// @Failure      500    {object} map[string]interface{}
// @Router       /projects/public/{user_id} [get]
func (h *ProjectHandler) GetPublicProjects(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid user ID"})
		return
	}

	skip, limit := h.parseQueryArgs(c)

	items, total, err := h.service.GetProjects(userID, skip, limit, true)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to fetch public projects"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"projects": items,
		"total":    total,
	})
}

// GetProject
// @Summary      Get project by ID
// @Description  Get a single project by its ID
// @Tags         projects
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        project_id  path      string  true  "Project ID"
// @Success      200    {object} models.Project
// @Failure      400    {object} map[string]interface{}
// @Failure      404    {object} map[string]interface{}
// @Router       /projects/{project_id} [get]
func (h *ProjectHandler) GetProject(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	exp, err := h.service.GetProjectByID(userID, id, false)
	if err != nil || exp == nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, exp)
}

// UpdateProject
// @Summary      Update project
// @Description  Update an existing project
// @Tags         projects
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        project_id  path      string  true  "Project ID"
// @Param        request        body      map[string]interface{} true "Update Data"
// @Success      200    {object} models.Project
// @Failure      400    {object} map[string]interface{}
// @Failure      404    {object} map[string]interface{}
// @Router       /projects/{project_id} [put]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid request body"})
		return
	}

	if t, ok := updateData["type"].(string); ok && t != "github" && t != "custom" {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "type must be github or custom"})
		return
	}

	updated, err := h.service.UpdateProject(userID, id, updateData)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Project not found or update failed"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

type projectVisibilityReq struct {
	IsVisible bool `json:"is_visible"`
}

// UpdateVisibility
// @Summary      Update project visibility
// @Description  Change the visibility of a project
// @Tags         projects
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        project_id  path      string  true  "Project ID"
// @Param        request        body      projectVisibilityReq true "Visibility Update"
// @Success      200    {object} models.Project
// @Failure      400    {object} map[string]interface{}
// @Failure      404    {object} map[string]interface{}
// @Router       /projects/{project_id}/visibility [patch]
func (h *ProjectHandler) UpdateVisibility(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	var req projectVisibilityReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid body"})
		return
	}

	updated, err := h.service.UpdateProjectVisibility(userID, id, req.IsVisible)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteProject
// @Summary      Delete project
// @Description  Delete a project
// @Tags         projects
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        project_id  path      string  true  "Project ID"
// @Success      204
// @Failure      400    {object} map[string]interface{}
// @Failure      404    {object} map[string]interface{}
// @Router       /projects/{project_id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	err = h.service.DeleteProject(userID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Project not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
