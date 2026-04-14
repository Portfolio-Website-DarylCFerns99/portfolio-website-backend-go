package chatbot

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"portfolio-website-backend/internal/dto"
	"portfolio-website-backend/internal/models"
)

type MockChatService struct {
	mock.Mock
}

func (m *MockChatService) GetOrCreateSession(sessionID, userID uuid.UUID) error {
	args := m.Called(sessionID, userID)
	return args.Error(0)
}

func (m *MockChatService) SaveMessage(sessionID uuid.UUID, role, content string) error {
	args := m.Called(sessionID, role, content)
	return args.Error(0)
}

func (m *MockChatService) BuildChatHistory(sessionID uuid.UUID, limit int) ([]*genai.Content, []map[string]interface{}, error) {
	args := m.Called(sessionID, limit)
	var contents []*genai.Content
	if c := args.Get(0); c != nil {
		contents = c.([]*genai.Content)
	}
	var payloads []map[string]interface{}
	if p := args.Get(1); p != nil {
		payloads = p.([]map[string]interface{})
	}
	return contents, payloads, args.Error(2)
}

func (m *MockChatService) GenerateStream(ctx context.Context, sessionID, userID uuid.UUID, query string, history []*genai.Content) (*genai.GenerateContentResponseIterator, *genai.Client, error) {
	args := m.Called(ctx, sessionID, userID, query, history)
	var iter *genai.GenerateContentResponseIterator
	if i := args.Get(0); i != nil {
		iter = i.(*genai.GenerateContentResponseIterator)
	}
	var client *genai.Client
	if c := args.Get(1); c != nil {
		client = c.(*genai.Client)
	}
	return iter, client, args.Error(2)
}

type MockVectorService struct {
	mock.Mock
}

func (m *MockVectorService) SyncUserData(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	args := m.Called(ctx, userID)
	var res map[string]interface{}
	if r := args.Get(0); r != nil {
		res = r.(map[string]interface{})
	}
	return res, args.Error(1)
}

func (m *MockVectorService) Search(ctx context.Context, text string, userID uuid.UUID, limit int, filters []string) ([]models.VectorEmbedding, error) {
	args := m.Called(ctx, text, userID, limit, filters)
	var res []models.VectorEmbedding
	if r := args.Get(0); r != nil {
		res = r.([]models.VectorEmbedding)
	}
	return res, args.Error(1)
}

type MockChatRepository struct {
	mock.Mock
}

func (m *MockChatRepository) GetSession(sessionID, userID uuid.UUID) (*models.ChatSession, error) {
	args := m.Called(sessionID, userID)
	var session *models.ChatSession
	if s := args.Get(0); s != nil {
		session = s.(*models.ChatSession)
	}
	return session, args.Error(1)
}

func (m *MockChatRepository) CreateSession(session *models.ChatSession) error {
	args := m.Called(session)
	return args.Error(0)
}

func (m *MockChatRepository) GetRecentMessages(sessionID uuid.UUID, limit int) ([]models.ChatMessage, error) {
	args := m.Called(sessionID, limit)
	var msgs []models.ChatMessage
	if msg := args.Get(0); msg != nil {
		msgs = msg.([]models.ChatMessage)
	}
	return msgs, args.Error(1)
}

func (m *MockChatRepository) AddMessage(message *models.ChatMessage) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *MockChatRepository) GetAllSessions(userID uuid.UUID, limit, offset int) ([]dto.ChatSessionResponse, error) {
	args := m.Called(userID, limit, offset)
	var res []dto.ChatSessionResponse
	if r := args.Get(0); r != nil {
		res = r.([]dto.ChatSessionResponse)
	}
	return res, args.Error(1)
}

func (m *MockChatRepository) GetSessionMessages(sessionID, userID uuid.UUID) ([]dto.ChatMessageResponse, error) {
	args := m.Called(sessionID, userID)
	var res []dto.ChatMessageResponse
	if r := args.Get(0); r != nil {
		res = r.([]dto.ChatMessageResponse)
	}
	return res, args.Error(1)
}
