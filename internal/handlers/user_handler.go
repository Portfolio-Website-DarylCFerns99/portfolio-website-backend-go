package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/repository"
	"portfolio-website-backend/internal/security"
	"portfolio-website-backend/internal/services"
)

type UserHandler struct {
	userService services.UserService
	userRepo    repository.UserRepository
}

func NewUserHandler(userService services.UserService, userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{
		userService: userService,
		userRepo:    userRepo,
	}
}

// RegisterRoutes registers the user endpoints
func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc, adminAuthMiddleware gin.HandlerFunc) {
	// Public routes
	r.POST("/login", h.Login)
	r.GET("/public-data/:user_id", h.GetPublicData)

	// Admin routes
	admin := r.Group("/admin/users")
	admin.Use(adminAuthMiddleware)
	{
		admin.POST("", h.CreateUser)
		admin.GET("", h.ListUsers)
		admin.GET("/:id", h.GetUser)
	}

	// Protected routes
	protected := r.Group("/")
	protected.Use(authMiddleware)
	{
		protected.GET("/profile", h.GetProfile)
		protected.PUT("/profile", h.UpdateProfile)
	}
}

type loginRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// Login
// @Summary      User login
// @Description  Login and get an access token
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      loginRequest  true  "Login Request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid form data"})
		return
	}

	// Try finding user by username or email
	user, _ := h.userRepo.GetByEmailorUsername(req.Username)

	if user == nil || !security.VerifyPassword(req.Password, user.HashedPassword) {
		c.Header("WWW-Authenticate", "Bearer")
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "Invalid credentials"})
		return
	}

	token, err := security.CreateAccessToken(map[string]interface{}{
		"sub": user.ID.String(),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"token_type":   "bearer",
		"user":         user,
	})
}

// GetProfile
// @Summary      Get current user profile
// @Description  Get the profile of the currently logged in user
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  models.User
// @Failure      401      {object}  map[string]interface{}
// @Router       /profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	user, exists := c.Get("current_user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "Not authenticated"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile
// @Summary      Update user profile
// @Description  Update the profile of the currently logged in user
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      map[string]interface{}  true  "Update Data"
// @Success      200      {object}  models.User
// @Failure      400      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userObj, exists := c.Get("current_user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "Not authenticated"})
		return
	}
	currentUser := userObj.(*models.User)

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid request format"})
		return
	}

	// Remove unmodifiable fields
	delete(updateData, "id")
	delete(updateData, "hashed_password") // Password should be changed via special endpoint if needed

	// Validate unique fields
	if un, ok := updateData["username"].(string); ok && un != currentUser.Username {
		existing, _ := h.userRepo.GetByEmailorUsername(un)
		if existing != nil {
			c.JSON(http.StatusBadRequest, gin.H{"detail": "Username already exists"})
			return
		}
	}
	if em, ok := updateData["email"].(string); ok && em != currentUser.Email {
		existing, _ := h.userRepo.GetByEmailorUsername(em)
		if existing != nil {
			c.JSON(http.StatusBadRequest, gin.H{"detail": "Email already exists"})
			return
		}
	}

	updatedUser, err := h.userService.UpdateUserProfile(currentUser.ID, updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to update profile", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// GetPublicData
// @Summary      Get public data
// @Description  Get public portfolio data for a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user_id  path      string  true  "User ID"
// @Success      200      {object}  dto.PublicDataResponse
// @Failure      400      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /public-data/{user_id} [get]
func (h *UserHandler) GetPublicData(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid user ID format"})
		return
	}

	publicData, err := h.userService.GetPublicPortfolioData(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "User or public data not found", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, publicData)
}

// Admin Handlers

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Username string `json:"username"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid form data"})
		return
	}

	uname := req.Username
	if uname == "" {
		uname = req.Email
	}

	existingUser, _ := h.userRepo.GetByEmailorUsername(uname)
	if existingUser != nil && existingUser.Username == uname {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Username already exists"})
		return
	}
	existingEmail, _ := h.userRepo.GetByEmailorUsername(req.Email)
	if existingEmail != nil && existingEmail.Email == req.Email {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Email already exists"})
		return
	}

	hashedPassword, err := security.GetPasswordHash(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to hash password"})
		return
	}

	user := &models.User{
		Email:          req.Email,
		Username:       uname,
		HashedPassword: hashedPassword,
	}

	createdUser, err := h.userRepo.Create(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.userRepo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to list users", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid user ID"})
		return
	}

	user, err := h.userRepo.GetByID(id)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
