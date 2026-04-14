package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"portfolio-website-backend/internal/dto"
	"portfolio-website-backend/internal/services"
)

// ContactHandler handles HTTP requests for the contact form.
type ContactHandler struct {
	service services.ContactService
}

// NewContactHandler creates a new ContactHandler.
func NewContactHandler(service services.ContactService) *ContactHandler {
	return &ContactHandler{service: service}
}

// RegisterRoutes registers the contact routes on the given router group.
// The contact endpoint is intentionally public — no auth middleware is applied.
func (h *ContactHandler) RegisterRoutes(r *gin.RouterGroup) {
	contact := r.Group("/contact")
	{
		contact.POST("/:user_id", h.SendContactEmail)
	}
}

// SendContactEmail godoc
// @Summary      Send contact email
// @Description  Send a contact form email to the portfolio owner identified by user_id.
//
//	Two emails are dispatched: a confirmation to the sender
//	and a notification to the admin.
//
// @Tags         contact
// @Accept       json
// @Produce      json
// @Param        user_id  path      string              true  "Portfolio owner user ID"
// @Param        request  body      dto.ContactRequest  true  "Contact form data"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /contact/{user_id} [post]
func (h *ContactHandler) SendContactEmail(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid user ID"})
		return
	}

	var req dto.ContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	if err := h.service.SendContactEmail(userID, req); err != nil {
		switch err.Error() {
		case "user not found":
			c.JSON(http.StatusNotFound, gin.H{"detail": "User not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Emails sent successfully"})
}
