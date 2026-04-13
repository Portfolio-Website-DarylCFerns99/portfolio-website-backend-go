package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/services"
)

type ReviewHandler struct {
	service services.ReviewService
}

func NewReviewHandler(service services.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: service}
}

func (h *ReviewHandler) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	reviews := r.Group("/reviews")

	// Public routes
	reviews.GET("/public/:user_id", h.GetPublicReviews)

	// Protected routes
	protected := reviews.Group("")
	protected.Use(authMiddleware)
	{
		protected.POST("", h.CreateReview)
		protected.GET("", h.GetReviews)
		protected.GET("/:review_id", h.GetReview)
		protected.PUT("/:review_id", h.UpdateReview)
		protected.PATCH("/:review_id/visibility", h.UpdateVisibility)
		protected.DELETE("/:review_id", h.DeleteReview)
	}
}

func (h *ReviewHandler) parseQueryArgs(c *gin.Context) (skip, limit int) {
	skip, _ = strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "100"))
	return
}

func (h *ReviewHandler) getUserID(c *gin.Context) (uuid.UUID, bool) {
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

// CreateReview
// @Summary      Create review
// @Description  Create a new review
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      models.Review  true  "Review Data"
// @Success      201      {object}  models.Review
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /reviews [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	review.UserID = userID

	created, err := h.service.CreateReview(&review)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetReviews
// @Summary      Get reviews
// @Description  Get a paginated list of reviews (requires auth)
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        skip   query    int  false  "Skip count"
// @Param        limit  query    int  false  "Limit count"
// @Success      200    {object} map[string]interface{}
// @Failure      400    {object} map[string]interface{}
// @Failure      500    {object} map[string]interface{}
// @Router       /reviews [get]
func (h *ReviewHandler) GetReviews(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	skip, limit := h.parseQueryArgs(c)

	items, total, err := h.service.GetReviews(userID, skip, limit, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reviews": items,
		"total":   total,
	})
}

// GetPublicReviews
// @Summary      Get public reviews
// @Description  Get a paginated list of public reviews
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Param        user_id path     string  true "User ID"
// @Param        skip   query    int  false  "Skip count"
// @Param        limit  query    int  false  "Limit count"
// @Success      200    {object} map[string]interface{}
// @Failure      400    {object} map[string]interface{}
// @Failure      500    {object} map[string]interface{}
// @Router       /reviews/public/{user_id} [get]
func (h *ReviewHandler) GetPublicReviews(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid user ID"})
		return
	}

	skip, limit := h.parseQueryArgs(c)

	items, total, err := h.service.GetReviews(userID, skip, limit, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to fetch public reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reviews": items,
		"total":   total,
	})
}

// GetReview
// @Summary      Get review by ID
// @Description  Get a single review by its ID
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        review_id  path      string  true  "Review ID"
// @Success      200    {object} models.Review
// @Failure      400    {object} map[string]interface{}
// @Failure      404    {object} map[string]interface{}
// @Router       /reviews/{review_id} [get]
func (h *ReviewHandler) GetReview(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("review_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	review, err := h.service.GetReviewByID(userID, id, false)
	if err != nil || review == nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Review not found"})
		return
	}

	c.JSON(http.StatusOK, review)
}

// UpdateReview
// @Summary      Update review
// @Description  Update an existing review
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        review_id  path      string  true  "Review ID"
// @Param        request    body      map[string]interface{} true "Update Data"
// @Success      200    {object} models.Review
// @Failure      400    {object} map[string]interface{}
// @Failure      404    {object} map[string]interface{}
// @Router       /reviews/{review_id} [put]
func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("review_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid request body"})
		return
	}

	updated, err := h.service.UpdateReview(userID, id, updateData)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Review not found or update failed"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

type reviewVisibilityReq struct {
	IsVisible bool `json:"is_visible"`
}

// UpdateVisibility
// @Summary      Update review visibility
// @Description  Change the visibility of a review
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        review_id  path      string  true  "Review ID"
// @Param        request    body      reviewVisibilityReq true "Visibility Update"
// @Success      200    {object} models.Review
// @Failure      400    {object} map[string]interface{}
// @Failure      404    {object} map[string]interface{}
// @Router       /reviews/{review_id}/visibility [patch]
func (h *ReviewHandler) UpdateVisibility(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("review_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	var req reviewVisibilityReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid body"})
		return
	}

	updated, err := h.service.UpdateReviewVisibility(userID, id, req.IsVisible)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Review not found"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteReview
// @Summary      Delete review
// @Description  Delete a review
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        review_id  path      string  true  "Review ID"
// @Success      204
// @Failure      400    {object} map[string]interface{}
// @Failure      404    {object} map[string]interface{}
// @Router       /reviews/{review_id} [delete]
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("review_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	err = h.service.DeleteReview(userID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Review not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
