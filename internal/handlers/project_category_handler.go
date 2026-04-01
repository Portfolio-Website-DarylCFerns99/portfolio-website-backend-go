package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/services"
)

type ProjectCategoryHandler struct {
	service services.ProjectCategoryService
}

func NewProjectCategoryHandler(service services.ProjectCategoryService) *ProjectCategoryHandler {
	return &ProjectCategoryHandler{service: service}
}

func (h *ProjectCategoryHandler) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	categories := r.Group("/project-categories")

	// Public routes
	categories.GET("/public/:user_id", h.GetPublicCategories)

	// Protected routes
	protected := categories.Group("")
	protected.Use(authMiddleware)
	{
		protected.POST("", h.CreateCategory)
		protected.GET("", h.GetCategories)
		protected.GET("/:category_id", h.GetCategory)
		protected.PUT("/:category_id", h.UpdateCategory)
		protected.PATCH("/:category_id/visibility", h.UpdateVisibility)
		protected.DELETE("/:category_id", h.DeleteCategory)
	}
}

func (h *ProjectCategoryHandler) parseQueryArgs(c *gin.Context) (skip, limit int) {
	skip, _ = strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "100"))
	return
}

func (h *ProjectCategoryHandler) getUserID(c *gin.Context) (uuid.UUID, bool) {
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

type CategoryCreateReq struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	IsVisible   bool    `json:"is_visible"`
}

// CreateCategory
// @Summary      Create project category
// @Description  Create a new project category
// @Tags         project-categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CategoryCreateReq  true  "Category Data"
// @Success      201      {object}  models.ProjectCategory
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /project-categories [post]
func (h *ProjectCategoryHandler) CreateCategory(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	var req CategoryCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	cat := models.ProjectCategory{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		IsVisible:   req.IsVisible,
	}

	created, err := h.service.CreateCategory(&cat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to create project category"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetCategories
// @Summary      Get project categories
// @Description  Get a paginated list of project categories (requires auth)
// @Tags         project-categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        skip   query    int  false  "Skip count"
// @Param        limit  query    int  false  "Limit count"
// @Success      200    {array}  models.ProjectCategory
// @Failure      400    {object} map[string]interface{}
// @Failure      500    {object} map[string]interface{}
// @Router       /project-categories [get]
func (h *ProjectCategoryHandler) GetCategories(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	skip, limit := h.parseQueryArgs(c)

	items, err := h.service.GetCategories(userID, skip, limit, false)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to fetch project categories"})
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetPublicCategories
// @Summary      Get public project categories
// @Description  Get a paginated list of public project categories
// @Tags         project-categories
// @Accept       json
// @Produce      json
// @Param        user_id path     string  true "User ID"
// @Param        skip   query    int  false  "Skip count"
// @Param        limit  query    int  false  "Limit count"
// @Success      200    {array}  models.ProjectCategory
// @Failure      400    {object} map[string]interface{}
// @Failure      500    {object} map[string]interface{}
// @Router       /project-categories/public/{user_id} [get]
func (h *ProjectCategoryHandler) GetPublicCategories(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid user ID"})
		return
	}

	skip, limit := h.parseQueryArgs(c)

	items, err := h.service.GetCategories(userID, skip, limit, true)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to fetch public project categories"})
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetCategory
// @Summary      Get project category by ID
// @Description  Get a single project category by its ID
// @Tags         project-categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        category_id  path      string  true  "Category ID"
// @Success      200    {object} models.ProjectCategory
// @Failure      400    {object} map[string]interface{}
// @Failure      404    {object} map[string]interface{}
// @Router       /project-categories/{category_id} [get]
func (h *ProjectCategoryHandler) GetCategory(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("category_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	cat, err := h.service.GetCategoryByID(userID, id, false)
	if err != nil || cat == nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Project category not found"})
		return
	}

	c.JSON(http.StatusOK, cat)
}

// UpdateCategory
// @Summary      Update project category
// @Description  Update an existing project category
// @Tags         project-categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        category_id  path      string  true  "Category ID"
// @Param        request        body      map[string]interface{} true "Update Data"
// @Success      200    {object} models.ProjectCategory
// @Failure      400    {object} map[string]interface{}
// @Failure      404    {object} map[string]interface{}
// @Router       /project-categories/{category_id} [put]
func (h *ProjectCategoryHandler) UpdateCategory(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("category_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid request body"})
		return
	}

	updated, err := h.service.UpdateCategory(userID, id, updateData)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Project category not found or update failed"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

type categoryVisibilityReq struct {
	IsVisible bool `json:"is_visible"`
}

// UpdateVisibility
// @Summary      Update project category visibility
// @Description  Change the visibility of a project category
// @Tags         project-categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        category_id  path      string  true  "Category ID"
// @Param        request        body      categoryVisibilityReq true "Visibility Update"
// @Success      200    {object} models.ProjectCategory
// @Failure      400    {object} map[string]interface{}
// @Failure      404    {object} map[string]interface{}
// @Router       /project-categories/{category_id}/visibility [patch]
func (h *ProjectCategoryHandler) UpdateVisibility(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("category_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	var req categoryVisibilityReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid body"})
		return
	}

	updated, err := h.service.UpdateCategory(userID, id, map[string]interface{}{"is_visible": req.IsVisible})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Project category not found"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteCategory
// @Summary      Delete project category
// @Description  Delete a project category
// @Tags         project-categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        category_id  path      string  true  "Category ID"
// @Success      204
// @Failure      400    {object} map[string]interface{}
// @Failure      404    {object} map[string]interface{}
// @Router       /project-categories/{category_id} [delete]
func (h *ProjectCategoryHandler) DeleteCategory(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("category_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	err = h.service.DeleteCategory(userID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Project category not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
