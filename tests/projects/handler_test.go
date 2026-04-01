package projects

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

func TestProjectHandler_Create(t *testing.T) {
	mockSvc := new(MockProjectService)
	handler := handlers.NewProjectHandler(mockSvc)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	mockSvc.On("CreateProject", mock.Anything).Return(&models.Project{
		BaseModel: models.BaseModel{ID: uuid.New()},
		Title:     "New Project",
	}, nil)

	body, _ := json.Marshal(map[string]interface{}{"title": "New Project", "type": "custom", "is_visible": true})
	req, _ := http.NewRequest(http.MethodPost, "/projects", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestProjectHandler_GetByID(t *testing.T) {
	mockSvc := new(MockProjectService)
	handler := handlers.NewProjectHandler(mockSvc)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	testID := uuid.New()
	mockSvc.On("GetProjectByID", testUser.ID, testID, false).Return(&models.Project{
		BaseModel: models.BaseModel{ID: testID},
		Title:     "Test Project",
	}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/projects/"+testID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProjectHandler_GetProjects(t *testing.T) {
	mockSvc := new(MockProjectService)
	handler := handlers.NewProjectHandler(mockSvc)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	mockSvc.On("GetProjects", testUser.ID, 0, 100, false).Return([]models.Project{}, int64(0), nil)

	req, _ := http.NewRequest(http.MethodGet, "/projects", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProjectHandler_UpdateVisibility(t *testing.T) {
	mockSvc := new(MockProjectService)
	handler := handlers.NewProjectHandler(mockSvc)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	testID := uuid.New()
	mockSvc.On("UpdateProjectVisibility", testUser.ID, testID, true).Return(&models.Project{
		BaseModel: models.BaseModel{ID: testID},
		IsVisible: true,
	}, nil)

	body, _ := json.Marshal(map[string]bool{"is_visible": true})
	req, _ := http.NewRequest(http.MethodPatch, "/projects/"+testID.String()+"/visibility", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProjectHandler_GetPublic(t *testing.T) {
	mockSvc := new(MockProjectService)
	handler := handlers.NewProjectHandler(mockSvc)
	r := common.SetupRouter(nil, handler.RegisterRoutes)

	targetUserID := uuid.New()
	mockSvc.On("GetProjects", targetUserID, 0, 100, true).Return([]models.Project{}, int64(0), nil)

	req, _ := http.NewRequest(http.MethodGet, "/projects/public/"+targetUserID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProjectHandler_Delete(t *testing.T) {
	mockSvc := new(MockProjectService)
	handler := handlers.NewProjectHandler(mockSvc)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	testID := uuid.New()
	mockSvc.On("DeleteProject", testUser.ID, testID).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/projects/"+testID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
