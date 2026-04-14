package chatbot

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"portfolio-website-backend/internal/dto"
	"portfolio-website-backend/internal/handlers"
	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/tests/common"
)

func TestChatHandler_SyncContext(t *testing.T) {
	mockChatSvc := new(MockChatService)
	mockVectorSvc := new(MockVectorService)
	mockChatRepo := new(MockChatRepository)

	handler := handlers.NewChatHandler(mockChatSvc, mockVectorSvc, mockChatRepo)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	mockVectorSvc.On("SyncUserData", mock.Anything, testUser.ID).Return(map[string]interface{}{
		"status":         "success",
		"vectors_synced": float64(5),
	}, nil)

	req, _ := http.NewRequest(http.MethodPost, "/chatbot/sync", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	// Unmarshal to verify response payload
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp["status"])
	assert.Equal(t, float64(5), resp["vectors_synced"])
}

func TestChatHandler_GetChatSessions(t *testing.T) {
	mockChatSvc := new(MockChatService)
	mockVectorSvc := new(MockVectorService)
	mockChatRepo := new(MockChatRepository)

	handler := handlers.NewChatHandler(mockChatSvc, mockVectorSvc, mockChatRepo)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	mockChatRepo.On("GetAllSessions", testUser.ID, 50, 0).Return([]dto.ChatSessionResponse{
		{
			ID:           uuid.New(),
			UserID:       testUser.ID,
			CreatedAt:    time.Now(),
			MessageCount: 10,
			LastActive:   time.Now(),
		},
	}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/chatbot/sessions", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var sessions []dto.ChatSessionResponse
	err := json.Unmarshal(w.Body.Bytes(), &sessions)
	assert.NoError(t, err)
	assert.Len(t, sessions, 1)
}

func TestChatHandler_GetSessionMessages_Success(t *testing.T) {
	mockChatSvc := new(MockChatService)
	mockVectorSvc := new(MockVectorService)
	mockChatRepo := new(MockChatRepository)

	handler := handlers.NewChatHandler(mockChatSvc, mockVectorSvc, mockChatRepo)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	testSessionID := uuid.New()
	mockChatRepo.On("GetSessionMessages", testSessionID, testUser.ID).Return([]dto.ChatMessageResponse{
		{
			ID:        uuid.New(),
			Sender:    "bot",
			Content:   "Hello!",
			CreatedAt: time.Now(),
		},
	}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/chatbot/sessions/"+testSessionID.String()+"/messages", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var messages []dto.ChatMessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &messages)
	assert.NoError(t, err)
	assert.Len(t, messages, 1)
	assert.Equal(t, "bot", messages[0].Sender)
}

func TestChatHandler_GetSessionMessages_InvalidUUID(t *testing.T) {
	mockChatSvc := new(MockChatService)
	mockVectorSvc := new(MockVectorService)
	mockChatRepo := new(MockChatRepository)

	handler := handlers.NewChatHandler(mockChatSvc, mockVectorSvc, mockChatRepo)

	testUser := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	r := common.SetupRouter(testUser, handler.RegisterRoutes)

	req, _ := http.NewRequest(http.MethodGet, "/chatbot/sessions/invalid-uuid/messages", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChatHandler_WebsocketChat_MissingQueries(t *testing.T) {
	mockChatSvc := new(MockChatService)
	mockVectorSvc := new(MockVectorService)
	mockChatRepo := new(MockChatRepository)

	handler := handlers.NewChatHandler(mockChatSvc, mockVectorSvc, mockChatRepo)

	r := common.SetupRouter(nil, handler.RegisterRoutes)

	req, _ := http.NewRequest(http.MethodGet, "/chatbot/ws/chat", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChatHandler_WebsocketChat_InvalidUUIDs(t *testing.T) {
	mockChatSvc := new(MockChatService)
	mockVectorSvc := new(MockVectorService)
	mockChatRepo := new(MockChatRepository)

	handler := handlers.NewChatHandler(mockChatSvc, mockVectorSvc, mockChatRepo)

	r := common.SetupRouter(nil, handler.RegisterRoutes)

	req, _ := http.NewRequest(http.MethodGet, "/chatbot/ws/chat?session_id=bad&user_id=bad", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
