package skills

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

func TestSkillHandler_CreateSkillGroup(t *testing.T) {
	mockSvc := new(MockSkillService)
	handler := handlers.NewSkillHandler(mockSvc)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	mockSvc.On("CreateSkillGroup", mock.Anything).Return(&models.SkillGroup{
		BaseModel: models.BaseModel{ID: uuid.New()},
		Name:      "Languages",
	}, nil)

	body, _ := json.Marshal(map[string]interface{}{"name": "Languages", "is_visible": true})
	req, _ := http.NewRequest(http.MethodPost, "/skills/groups", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestSkillHandler_GetSkillGroups(t *testing.T) {
	mockSvc := new(MockSkillService)
	handler := handlers.NewSkillHandler(mockSvc)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	mockSvc.On("GetSkillGroups", testUser.ID, 0, 100, false).Return([]models.SkillGroup{}, int64(0), nil)

	req, _ := http.NewRequest(http.MethodGet, "/skills/groups", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSkillHandler_CreateSkill(t *testing.T) {
	mockSvc := new(MockSkillService)
	handler := handlers.NewSkillHandler(mockSvc)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	groupID := uuid.New()
	mockSvc.On("CreateSkill", mock.Anything).Return(&models.Skill{
		BaseModel:    models.BaseModel{ID: uuid.New()},
		SkillGroupID: groupID,
		Name:         "Go",
	}, nil)

	body, _ := json.Marshal(map[string]interface{}{"skill_group_id": groupID, "name": "Go"})
	req, _ := http.NewRequest(http.MethodPost, "/skills", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestSkillHandler_GetSkills(t *testing.T) {
	mockSvc := new(MockSkillService)
	handler := handlers.NewSkillHandler(mockSvc)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	mockSvc.On("GetSkills", testUser.ID, 0, 100, false).Return([]models.Skill{}, int64(0), nil)

	req, _ := http.NewRequest(http.MethodGet, "/skills", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSkillHandler_UpdateVisibility(t *testing.T) {
	mockSvc := new(MockSkillService)
	handler := handlers.NewSkillHandler(mockSvc)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	testID := uuid.New()
	mockSvc.On("UpdateSkillVisibility", testUser.ID, testID, true).Return(&models.Skill{
		BaseModel: models.BaseModel{ID: testID},
		IsVisible: true,
	}, nil)

	body, _ := json.Marshal(map[string]bool{"is_visible": true})
	req, _ := http.NewRequest(http.MethodPatch, "/skills/"+testID.String()+"/visibility", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
