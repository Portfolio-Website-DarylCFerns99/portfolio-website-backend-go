package users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"portfolio-website-backend/internal/dto"
	"portfolio-website-backend/internal/handlers"
	"portfolio-website-backend/internal/middleware"
	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/security"
	"portfolio-website-backend/tests/common"
)

func TestLogin_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockSvc := new(MockUserService)
	handler := handlers.NewUserHandler(mockSvc, mockRepo)
	r := common.SetupRouterWithAdmin(nil, nil, nil, handler.RegisterRoutes)

	password := "testpass123"
	hashedPassword, _ := security.GetPasswordHash(password)

	testUser := &models.User{
		BaseModel:      models.BaseModel{ID: uuid.New()},
		Username:       "testuser",
		Email:          "test@example.com",
		HashedPassword: hashedPassword,
	}

	mockRepo.On("GetByEmailorUsername", "testuser").Return(testUser, nil)

	loginReq := map[string]string{
		"username": "testuser",
		"password": password,
	}
	jsonBody, _ := json.Marshal(loginReq)

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "access_token")
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockSvc := new(MockUserService)
	handler := handlers.NewUserHandler(mockSvc, mockRepo)
	r := common.SetupRouterWithAdmin(nil, nil, nil, handler.RegisterRoutes)

	password := "testpass123"
	hashedPassword, _ := security.GetPasswordHash(password)

	testUser := &models.User{
		BaseModel:      models.BaseModel{ID: uuid.New()},
		Username:       "testuser",
		Email:          "test@example.com",
		HashedPassword: hashedPassword,
	}

	mockRepo.On("GetByEmailorUsername", "testuser").Return(testUser, nil)

	loginReq := map[string]string{
		"username": "testuser",
		"password": "wrongpassword",
	}
	jsonBody, _ := json.Marshal(loginReq)

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetPublicData_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockSvc := new(MockUserService)
	handler := handlers.NewUserHandler(mockSvc, mockRepo)
	r := common.SetupRouterWithAdmin(nil, nil, nil, handler.RegisterRoutes)

	testID := uuid.New()
	name := "John"
	surname := "Doe"
	publicData := &dto.PublicDataResponse{
		Name:    &name,
		Surname: &surname,
	}

	mockSvc.On("GetPublicPortfolioData", testID).Return(publicData, nil)

	req, _ := http.NewRequest(http.MethodGet, "/public-data/"+testID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response dto.PublicDataResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "John", *response.Name)
}

func TestGetPublicData_InvalidUUID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockSvc := new(MockUserService)
	handler := handlers.NewUserHandler(mockSvc, mockRepo)
	r := common.SetupRouterWithAdmin(nil, nil, nil, handler.RegisterRoutes)

	req, _ := http.NewRequest(http.MethodGet, "/public-data/invalid-uuid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetProfile_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockSvc := new(MockUserService)
	handler := handlers.NewUserHandler(mockSvc, mockRepo)

	testUser := &models.User{
		BaseModel: models.BaseModel{ID: uuid.New()},
		Username:  "testuser",
	}
	r := common.SetupRouterWithAdmin(testUser, nil, nil, handler.RegisterRoutes)

	req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetProfile_Unauthenticated(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockSvc := new(MockUserService)
	handler := handlers.NewUserHandler(mockSvc, mockRepo)
	r := common.SetupRouterWithAdmin(nil, nil, nil, handler.RegisterRoutes)

	req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUpdateProfile_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockSvc := new(MockUserService)
	handler := handlers.NewUserHandler(mockSvc, mockRepo)

	testUser := &models.User{
		BaseModel: models.BaseModel{ID: uuid.New()},
		Username:  "olduser",
		Email:     "old@example.com",
	}

	r := common.SetupRouterWithAdmin(testUser, nil, nil, handler.RegisterRoutes)

	updateDataMap := map[string]interface{}{"username": "newuser"}
	mockRepo.On("GetByEmailorUsername", "newuser").Return((*models.User)(nil), nil)

	updatedUser := &models.User{
		BaseModel: models.BaseModel{ID: testUser.ID},
		Username:  "newuser",
		Email:     "old@example.com",
	}
	mockSvc.On("UpdateUserProfile", testUser.ID, updateDataMap).Return(updatedUser, nil)

	jsonBody, _ := json.Marshal(updateDataMap)
	req, _ := http.NewRequest(http.MethodPut, "/profile", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateProfile_DuplicateUsername(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockSvc := new(MockUserService)
	handler := handlers.NewUserHandler(mockSvc, mockRepo)

	testUser := &models.User{
		BaseModel: models.BaseModel{ID: uuid.New()},
		Username:  "olduser",
	}

	r := common.SetupRouterWithAdmin(testUser, nil, nil, handler.RegisterRoutes)

	updateDataMap := map[string]interface{}{"username": "existinguser"}
	existingUser := &models.User{
		BaseModel: models.BaseModel{ID: uuid.New()},
		Username:  "existinguser",
	}
	mockRepo.On("GetByEmailorUsername", "existinguser").Return(existingUser, nil)

	jsonBody, _ := json.Marshal(updateDataMap)
	req, _ := http.NewRequest(http.MethodPut, "/profile", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUser_AdminSuccess(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockSvc := new(MockUserService)
	handler := handlers.NewUserHandler(mockSvc, mockRepo)

	os.Setenv("ADMIN_API_KEY", "secret")
	defer os.Unsetenv("ADMIN_API_KEY")

	r := common.SetupRouterWithAdmin(nil, nil, middleware.RequireAdminAuth(), handler.RegisterRoutes)

	mockRepo.On("GetByEmailorUsername", "test@test.com").Return((*models.User)(nil), nil)
	mockRepo.On("GetByEmailorUsername", "adminuser").Return((*models.User)(nil), nil)
	mockRepo.On("Create", mock.Anything).Return(&models.User{
		BaseModel: models.BaseModel{ID: uuid.New()},
		Username:  "adminuser",
		Email:     "test@test.com",
	}, nil)

	reqBody := map[string]string{
		"email":    "test@test.com",
		"username": "adminuser",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/admin/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Admin-Api-Key", "secret")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestListUsers_AdminSuccess(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockSvc := new(MockUserService)
	handler := handlers.NewUserHandler(mockSvc, mockRepo)

	os.Setenv("ADMIN_API_KEY", "secret")
	defer os.Unsetenv("ADMIN_API_KEY")

	r := common.SetupRouterWithAdmin(nil, nil, middleware.RequireAdminAuth(), handler.RegisterRoutes)

	mockRepo.On("List").Return([]models.User{}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/admin/users", nil)
	req.Header.Set("X-Admin-Api-Key", "secret")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUser_AdminSuccess(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockSvc := new(MockUserService)
	handler := handlers.NewUserHandler(mockSvc, mockRepo)

	os.Setenv("ADMIN_API_KEY", "secret")
	defer os.Unsetenv("ADMIN_API_KEY")

	r := common.SetupRouterWithAdmin(nil, nil, middleware.RequireAdminAuth(), handler.RegisterRoutes)

	testID := uuid.New()
	mockRepo.On("GetByID", testID).Return(&models.User{
		BaseModel: models.BaseModel{ID: testID},
		Username:  "testuser",
	}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/admin/users/"+testID.String(), nil)
	req.Header.Set("X-Admin-Api-Key", "secret")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
