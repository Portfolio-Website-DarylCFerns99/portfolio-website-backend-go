package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/services"
)

type SkillHandler struct {
	service services.SkillService
}

func NewSkillHandler(service services.SkillService) *SkillHandler {
	return &SkillHandler{service: service}
}

func (h *SkillHandler) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	skills := r.Group("/skills")

	// Public routes
	skills.GET("/public/:user_id", h.GetPublicSkillGroups)

	// Protected routes
	protected := skills.Group("")
	protected.Use(authMiddleware)
	{
		// Skill Groups
		protected.POST("/groups", h.CreateSkillGroup)
		protected.GET("/groups", h.GetSkillGroups)
		protected.GET("/groups/:group_id", h.GetSkillGroup)
		protected.PUT("/groups/:group_id", h.UpdateSkillGroup)
		protected.PATCH("/groups/:group_id/visibility", h.UpdateSkillGroupVisibility)
		protected.DELETE("/groups/:group_id", h.DeleteSkillGroup)

		// Skills
		protected.POST("", h.CreateSkill)
		protected.GET("", h.GetSkills)
		protected.GET("/:skill_id", h.GetSkill)
		protected.PUT("/:skill_id", h.UpdateSkill)
		protected.PATCH("/:skill_id/visibility", h.UpdateSkillVisibility)
		protected.DELETE("/:skill_id", h.DeleteSkill)
	}
}

func (h *SkillHandler) parseQueryArgs(c *gin.Context) (skip, limit int) {
	skip, _ = strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "100"))
	return
}

func (h *SkillHandler) getUserID(c *gin.Context) (uuid.UUID, bool) {
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

// --- Skill Group Handlers ---

// CreateSkillGroup
// @Summary      Create skill group
// @Description  Create a new skill group
// @Tags         skills
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      models.SkillGroup  true  "Skill Group Data"
// @Success      201      {object}  models.SkillGroup
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /skills/groups [post]
func (h *SkillHandler) CreateSkillGroup(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	var group models.SkillGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	group.UserID = userID
	for i := range group.Skills {
		group.Skills[i].UserID = userID
	}

	created, err := h.service.CreateSkillGroup(&group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to create skill group"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetSkillGroups
// @Summary      Get skill groups
// @Description  Get a paginated list of skill groups (requires auth)
// @Tags         skills
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        skip   query    int  false  "Skip count"
// @Param        limit  query    int  false  "Limit count"
// @Success      200    {object} map[string]interface{}
// @Failure      400    {object} map[string]interface{}
// @Failure      500    {object} map[string]interface{}
// @Router       /skills/groups [get]
func (h *SkillHandler) GetSkillGroups(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	skip, limit := h.parseQueryArgs(c)

	items, total, err := h.service.GetSkillGroups(userID, skip, limit, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to fetch skill groups"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"skill_groups": items,
		"total":        total,
	})
}

// GetPublicSkillGroups
// @Summary      Get public skill groups
// @Description  Get a paginated list of public skill groups
// @Tags         skills
// @Accept       json
// @Produce      json
// @Param        user_id path     string  true "User ID"
// @Param        skip   query    int  false  "Skip count"
// @Param        limit  query    int  false  "Limit count"
// @Success      200    {object} map[string]interface{}
// @Failure      400    {object} map[string]interface{}
// @Failure      500    {object} map[string]interface{}
// @Router       /skills/public/{user_id} [get]
func (h *SkillHandler) GetPublicSkillGroups(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid user ID"})
		return
	}

	skip, limit := h.parseQueryArgs(c)

	items, total, err := h.service.GetSkillGroups(userID, skip, limit, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to fetch public skill groups"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"skill_groups": items,
		"total":        total,
	})
}

// GetSkillGroup
// @Summary      Get skill group by ID
// @Description  Get a single skill group by its ID
// @Tags         skills
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        group_id  path      string  true  "Group ID"
// @Success      200    {object} models.SkillGroup
// @Failure      400    {object} map[string]interface{}
// @Failure      404    {object} map[string]interface{}
// @Router       /skills/groups/{group_id} [get]
func (h *SkillHandler) GetSkillGroup(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	group, err := h.service.GetSkillGroupByID(userID, id, false)
	if err != nil || group == nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Skill group not found"})
		return
	}

	c.JSON(http.StatusOK, group)
}

// UpdateSkillGroup
// @Summary      Update skill group
// @Description  Update an existing skill group
// @Tags         skills
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        group_id  path      string  true  "Group ID"
// @Param        request   body      map[string]interface{} true "Update Data"
// @Success      200    {object} models.SkillGroup
// @Failure      400    {object} map[string]interface{}
// @Failure      404    {object} map[string]interface{}
// @Router       /skills/groups/{group_id} [put]
func (h *SkillHandler) UpdateSkillGroup(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid request body"})
		return
	}

	updated, err := h.service.UpdateSkillGroup(userID, id, updateData)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Skill group not found or update failed"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

type skillVisibilityReq struct {
	IsVisible bool `json:"is_visible"`
}

// UpdateSkillGroupVisibility
// @Summary      Update skill group visibility
// @Description  Change the visibility of a skill group
// @Tags         skills
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        group_id  path      string  true  "Group ID"
// @Param        request   body      skillVisibilityReq true "Visibility Update"
// @Success      200    {object} models.SkillGroup
// @Failure      400    {object} map[string]interface{}
// @Failure      404    {object} map[string]interface{}
// @Router       /skills/groups/{group_id}/visibility [patch]
func (h *SkillHandler) UpdateSkillGroupVisibility(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	var req skillVisibilityReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid body"})
		return
	}

	updated, err := h.service.UpdateSkillGroupVisibility(userID, id, req.IsVisible)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Skill group not found"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteSkillGroup
// @Summary      Delete skill group
// @Description  Delete a skill group
// @Tags         skills
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        group_id  path      string  true  "Group ID"
// @Success      204
// @Failure      400    {object} map[string]interface{}
// @Failure      404    {object} map[string]interface{}
// @Router       /skills/groups/{group_id} [delete]
func (h *SkillHandler) DeleteSkillGroup(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	err = h.service.DeleteSkillGroup(userID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Skill group not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

// --- Skill Handlers ---

// CreateSkill
// @Summary      Create skill
// @Description  Create a new skill in a skill group
// @Tags         skills
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      models.Skill  true  "Skill Data"
// @Success      201      {object}  models.Skill
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /skills [post]
func (h *SkillHandler) CreateSkill(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	var skill models.Skill
	if err := c.ShouldBindJSON(&skill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	skill.UserID = userID

	created, err := h.service.CreateSkill(&skill)
	if err != nil {
		// Possibly group not found or unauthorized
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetSkills
// @Summary      Get skills
// @Description  Get a paginated list of skills
// @Tags         skills
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        skip   query    int  false  "Skip count"
// @Param        limit  query    int  false  "Limit count"
// @Success      200    {object} map[string]interface{}
// @Failure      400    {object} map[string]interface{}
// @Failure      500    {object} map[string]interface{}
// @Router       /skills [get]
func (h *SkillHandler) GetSkills(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	skip, limit := h.parseQueryArgs(c)

	items, total, err := h.service.GetSkills(userID, skip, limit, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to fetch skills"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"skills": items,
		"total":  total,
	})
}

// GetSkill
// @Summary      Get skill by ID
// @Description  Get a single skill by its ID
// @Tags         skills
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        skill_id  path      string  true  "Skill ID"
// @Success      200    {object} models.Skill
// @Failure      400    {object} map[string]interface{}
// @Failure      404    {object} map[string]interface{}
// @Router       /skills/{skill_id} [get]
func (h *SkillHandler) GetSkill(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("skill_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	skill, err := h.service.GetSkillByID(userID, id)
	if err != nil || skill == nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Skill not found"})
		return
	}

	c.JSON(http.StatusOK, skill)
}

// UpdateSkill
// @Summary      Update skill
// @Description  Update an existing skill
// @Tags         skills
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        skill_id  path      string  true  "Skill ID"
// @Param        request   body      map[string]interface{} true "Update Data"
// @Success      200    {object} models.Skill
// @Failure      400    {object} map[string]interface{}
// @Failure      404    {object} map[string]interface{}
// @Router       /skills/{skill_id} [put]
func (h *SkillHandler) UpdateSkill(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("skill_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid request body"})
		return
	}

	updated, err := h.service.UpdateSkill(userID, id, updateData)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Skill not found or update failed"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// UpdateSkillVisibility
// @Summary      Update skill visibility
// @Description  Change the visibility of a skill
// @Tags         skills
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        skill_id  path      string  true  "Skill ID"
// @Param        request   body      skillVisibilityReq true "Visibility Update"
// @Success      200    {object} models.Skill
// @Failure      400    {object} map[string]interface{}
// @Failure      404    {object} map[string]interface{}
// @Router       /skills/{skill_id}/visibility [patch]
func (h *SkillHandler) UpdateSkillVisibility(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("skill_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	var req skillVisibilityReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid body"})
		return
	}

	updated, err := h.service.UpdateSkillVisibility(userID, id, req.IsVisible)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Skill not found"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteSkill
// @Summary      Delete skill
// @Description  Delete a skill
// @Tags         skills
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        skill_id  path      string  true  "Skill ID"
// @Success      204
// @Failure      400    {object} map[string]interface{}
// @Failure      404    {object} map[string]interface{}
// @Router       /skills/{skill_id} [delete]
func (h *SkillHandler) DeleteSkill(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("skill_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	err = h.service.DeleteSkill(userID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Skill not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
