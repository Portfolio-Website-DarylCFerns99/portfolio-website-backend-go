package project_categories

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

func TestProjectCategoryHandler_Create(t *testing.T) {
	mockSvc := new(MockProjectCategoryService)
	handler := handlers.NewProjectCategoryHandler(mockSvc)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	mockSvc.On("CreateCategory", mock.Anything).Return(&models.ProjectCategory{
		BaseModel: models.BaseModel{ID: uuid.New()},
		Name:      "New Category",
	}, nil)

	body, _ := json.Marshal(map[string]interface{}{"name": "New Category", "is_visible": true})
	req, _ := http.NewRequest(http.MethodPost, "/project-categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestProjectCategoryHandler_GetByID(t *testing.T) {
	mockSvc := new(MockProjectCategoryService)
	handler := handlers.NewProjectCategoryHandler(mockSvc)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	testID := uuid.New()
	mockSvc.On("GetCategoryByID", testUser.ID, testID, false).Return(&models.ProjectCategory{
		BaseModel: models.BaseModel{ID: testID},
		Name:      "Test Category",
	}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/project-categories/"+testID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProjectCategoryHandler_GetCategories(t *testing.T) {
	mockSvc := new(MockProjectCategoryService)
	handler := handlers.NewProjectCategoryHandler(mockSvc)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	mockSvc.On("GetCategories", testUser.ID, 0, 100, false).Return([]models.ProjectCategory{}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/project-categories", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProjectCategoryHandler_UpdateVisibility(t *testing.T) {
	mockSvc := new(MockProjectCategoryService)
	handler := handlers.NewProjectCategoryHandler(mockSvc)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	testID := uuid.New()
	mockSvc.On("UpdateCategory", testUser.ID, testID, map[string]interface{}{"is_visible": true}).Return(&models.ProjectCategory{
		BaseModel: models.BaseModel{ID: testID},
		IsVisible: true,
	}, nil)

	body, _ := json.Marshal(map[string]bool{"is_visible": true})
	req, _ := http.NewRequest(http.MethodPatch, "/project-categories/"+testID.String()+"/visibility", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProjectCategoryHandler_GetPublic(t *testing.T) {
	mockSvc := new(MockProjectCategoryService)
	handler := handlers.NewProjectCategoryHandler(mockSvc)
	r := common.SetupRouter(nil, handler.RegisterRoutes)

	targetUserID := uuid.New()
	mockSvc.On("GetCategories", targetUserID, 0, 100, true).Return([]models.ProjectCategory{}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/project-categories/public/"+targetUserID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProjectCategoryHandler_Delete(t *testing.T) {
	mockSvc := new(MockProjectCategoryService)
	handler := handlers.NewProjectCategoryHandler(mockSvc)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	testID := uuid.New()
	mockSvc.On("DeleteCategory", testUser.ID, testID).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/project-categories/"+testID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
