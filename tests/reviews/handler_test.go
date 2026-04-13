package reviews

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

func TestReviewHandler_Create(t *testing.T) {
	mockSvc := new(MockReviewService)
	handler := handlers.NewReviewHandler(mockSvc)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	mockSvc.On("CreateReview", mock.Anything).Return(&models.Review{
		BaseModel: models.BaseModel{ID: uuid.New()},
		Name:      "Jane Doe",
		Content:   "Amazing work!",
	}, nil)

	body, _ := json.Marshal(map[string]interface{}{"name": "Jane Doe", "content": "Amazing work!", "rating": 5})
	req, _ := http.NewRequest(http.MethodPost, "/reviews", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestReviewHandler_GetByID(t *testing.T) {
	mockSvc := new(MockReviewService)
	handler := handlers.NewReviewHandler(mockSvc)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	testID := uuid.New()
	mockSvc.On("GetReviewByID", testUser.ID, testID, false).Return(&models.Review{
		BaseModel: models.BaseModel{ID: testID},
		Content:   "Great service!",
	}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/reviews/"+testID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestReviewHandler_GetReviews(t *testing.T) {
	mockSvc := new(MockReviewService)
	handler := handlers.NewReviewHandler(mockSvc)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	mockSvc.On("GetReviews", testUser.ID, 0, 100, false).Return([]models.Review{}, int64(0), nil)

	req, _ := http.NewRequest(http.MethodGet, "/reviews", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestReviewHandler_UpdateVisibility(t *testing.T) {
	mockSvc := new(MockReviewService)
	handler := handlers.NewReviewHandler(mockSvc)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	testID := uuid.New()
	mockSvc.On("UpdateReviewVisibility", testUser.ID, testID, true).Return(&models.Review{
		BaseModel: models.BaseModel{ID: testID},
		IsVisible: true,
	}, nil)

	body, _ := json.Marshal(map[string]bool{"is_visible": true})
	req, _ := http.NewRequest(http.MethodPatch, "/reviews/"+testID.String()+"/visibility", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestReviewHandler_GetPublic(t *testing.T) {
	mockSvc := new(MockReviewService)
	handler := handlers.NewReviewHandler(mockSvc)
	r := common.SetupRouter(nil, handler.RegisterRoutes)

	targetUserID := uuid.New()
	mockSvc.On("GetReviews", targetUserID, 0, 100, true).Return([]models.Review{}, int64(0), nil)

	req, _ := http.NewRequest(http.MethodGet, "/reviews/public/"+targetUserID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestReviewHandler_Delete(t *testing.T) {
	mockSvc := new(MockReviewService)
	handler := handlers.NewReviewHandler(mockSvc)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	testID := uuid.New()
	mockSvc.On("DeleteReview", testUser.ID, testID).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/reviews/"+testID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
