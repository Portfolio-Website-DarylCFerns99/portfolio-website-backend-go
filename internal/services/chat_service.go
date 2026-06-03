package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"portfolio-website-backend/internal/config"
	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/repository"
	"portfolio-website-backend/internal/utils"
)

type ChatService interface {
	GetOrCreateSession(sessionID, userID uuid.UUID) error
	SaveMessage(sessionID uuid.UUID, role, content string) error
	BuildChatHistory(sessionID uuid.UUID, limit int) ([]*genai.Content, []map[string]interface{}, error)
	GenerateStream(ctx context.Context, sessionID, userID uuid.UUID, query string, history []*genai.Content) (*genai.GenerateContentResponseIterator, *genai.Client, error)
}

type chatService struct {
	db         *gorm.DB
	chatRepo   repository.ChatRepository
	vectorSvc  VectorService
	llmFactory *utils.LLMFactory
}

func NewChatService(db *gorm.DB) ChatService {
	return &chatService{
		db:         db,
		chatRepo:   repository.NewChatRepository(db),
		vectorSvc:  NewVectorService(db),
		llmFactory: utils.NewLLMFactory(),
	}
}

func (s *chatService) GetOrCreateSession(sessionID, userID uuid.UUID) error {
	session, err := s.chatRepo.GetSession(sessionID, userID)
	if err != nil {
		return err
	}
	if session == nil {
		newSession := &models.ChatSession{
			BaseModel: models.BaseModel{ID: sessionID},
			UserID:    userID,
			Title:     nil, // Generate title later if needed
		}
		return s.chatRepo.CreateSession(newSession)
	}
	return nil
}

func (s *chatService) SaveMessage(sessionID uuid.UUID, role, content string) error {
	msg := &models.ChatMessage{
		SessionID: sessionID,
		Sender:    role,
		Role:      role,
		Content:   content,
	}
	return s.chatRepo.AddMessage(msg)
}

func (s *chatService) BuildChatHistory(sessionID uuid.UUID, limit int) ([]*genai.Content, []map[string]interface{}, error) {
	msgs, err := s.chatRepo.GetRecentMessages(sessionID, limit)
	if err != nil {
		return nil, nil, err
	}

	var history []*genai.Content
	var payload []map[string]interface{}

	for _, m := range msgs {
		role := "user"
		sender := "user"
		if m.Sender == "assistant" || m.Sender == "bot" {
			role = "model"
			sender = "bot"
		}

		history = append(history, &genai.Content{
			Parts: []genai.Part{genai.Text(m.Content)},
			Role:  role,
		})

		payload = append(payload, map[string]interface{}{
			"sender": sender,
			"text":   m.Content,
		})
	}
	return history, payload, nil
}

func (s *chatService) condenseQuery(ctx context.Context, client *genai.Client, query string, history []*genai.Content) string {
	if len(history) == 0 {
		return query
	}

	modelName := config.Envs.GeminiModel
	if modelName == "" {
		modelName = "gemini-1.5-flash"
	}
	model := client.GenerativeModel(modelName)

	var builder strings.Builder
	for _, h := range history {
		role := "User"
		if h.Role == "model" {
			role = "Assistant"
		}
		for _, part := range h.Parts {
			if txt, ok := part.(genai.Text); ok {
				builder.WriteString(fmt.Sprintf("%s: %s\n", role, string(txt)))
			}
		}
	}
	builder.WriteString(fmt.Sprintf("User: %s\n", query))

	prompt := fmt.Sprintf(`Given the following conversation history and the user's latest follow-up question, rewrite it into a single standalone search query containing all necessary context. Do NOT add any preamble, explanation, or formatting. Output ONLY the rewritten standalone query.

Conversation History:
%s
Standalone query:`, builder.String())

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil || resp == nil || len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return query
	}

	if txt, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
		condensed := strings.TrimSpace(string(txt))
		if condensed != "" {
			return condensed
		}
	}
	return query
}

func (s *chatService) GenerateStream(ctx context.Context, sessionID, userID uuid.UUID, query string, history []*genai.Content) (*genai.GenerateContentResponseIterator, *genai.Client, error) {
	client, err := s.llmFactory.CreateGeminiClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	searchQuery := query
	if len(history) > 0 {
		searchQuery = s.condenseQuery(ctx, client, query, history)
	}

	docs, err := s.vectorSvc.Search(ctx, searchQuery, userID, 7, nil)
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	var contexts []string
	for _, d := range docs {
		contexts = append(contexts, fmt.Sprintf("- %s", d.Content))
	}
	contextText := strings.Join(contexts, "\n")

	// Get User info for prompt
	var user models.User
	s.db.First(&user, "id = ?", userID)
	portfolioOwner := "the portfolio owner"
	if user.Name != nil && *user.Name != "" {
		portfolioOwner = *user.Name
	}

	// 3. Construct System Prompt
	systemPrompt := fmt.Sprintf(`You are a helpful portfolio assistant for %s.
Use the following context to answer the user's question.
If the answer is not in the context, just say you don't know, but be friendly.

CONTEXT:
%s`, portfolioOwner, contextText)

	modelName := config.Envs.GeminiModel
	if modelName == "" {
		modelName = "gemini-1.5-flash"
	}

	model := client.GenerativeModel(modelName)
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemPrompt)},
	}

	// Start chat with history
	chat := model.StartChat()
	chat.History = history

	iter := chat.SendMessageStream(ctx, genai.Text(query))
	return iter, client, nil
}
