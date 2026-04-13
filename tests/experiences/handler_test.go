package experiences

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"portfolio-website-backend/internal/handlers"
	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/tests/common"
)

func TestExperienceHandler_Create(t *testing.T) {
	mockSvc := new(MockExperienceService)
	handler := handlers.NewExperienceHandler(mockSvc)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	mockSvc.On("CreateExperience", mock.Anything).Return(&models.Experience{
		BaseModel: models.BaseModel{ID: uuid.New()},
		Title:     "New Job",
	}, nil)

	body, _ := json.Marshal(map[string]string{"title": "New Job", "type": "experience", "organization": "ACME", "start_date": "2024-01-01"})
	req, _ := http.NewRequest(http.MethodPost, "/experiences", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestExperienceHandler_GetByID(t *testing.T) {
	mockSvc := new(MockExperienceService)
	handler := handlers.NewExperienceHandler(mockSvc)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	testID := uuid.New()
	mockSvc.On("GetExperienceByID", testUser.ID, testID, false).Return(&models.Experience{
		BaseModel: models.BaseModel{ID: testID},
		Title:     "Test Role",
	}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/experiences/"+testID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestExperienceHandler_GetExperiences(t *testing.T) {
	mockSvc := new(MockExperienceService)
	handler := handlers.NewExperienceHandler(mockSvc)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	mockSvc.On("GetExperiences", testUser.ID, 0, 100, false).Return([]models.Experience{}, int64(0), nil)

	req, _ := http.NewRequest(http.MethodGet, "/experiences", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestExperienceHandler_UpdateVisibility(t *testing.T) {
	mockSvc := new(MockExperienceService)
	handler := handlers.NewExperienceHandler(mockSvc)

	// Inject the mock test user to bypass standard authentication Middleware
	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	testID := uuid.New()
	// Body sends is_visible=true, so mock must match true
	mockSvc.On("UpdateExperienceVisibility", testUser.ID, testID, true).Return(&models.Experience{
		BaseModel: models.BaseModel{ID: testID},
		IsVisible: true,
	}, nil)

	body, _ := json.Marshal(map[string]bool{"is_visible": true})
	req, _ := http.NewRequest(http.MethodPatch, "/experiences/"+testID.String()+"/visibility", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestExperienceHandler_GetPublic(t *testing.T) {
	mockSvc := new(MockExperienceService)
	handler := handlers.NewExperienceHandler(mockSvc)
	r := common.SetupRouter(nil, handler.RegisterRoutes)

	targetUserID := uuid.New()
	mockSvc.On("GetExperiences", targetUserID, 0, 100, true).Return([]models.Experience{}, int64(0), nil)

	req, _ := http.NewRequest(http.MethodGet, "/experiences/public/"+targetUserID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestExperienceHandler_Delete(t *testing.T) {
	mockSvc := new(MockExperienceService)
	handler := handlers.NewExperienceHandler(mockSvc)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	testID := uuid.New()
	mockSvc.On("DeleteExperience", testUser.ID, testID).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/experiences/"+testID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
