package repository

import (
	"portfolio-website-backend/internal/dto"
	"portfolio-website-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatRepository interface {
	GetSession(sessionID, userID uuid.UUID) (*models.ChatSession, error)
	CreateSession(session *models.ChatSession) error
	GetRecentMessages(sessionID uuid.UUID, limit int) ([]models.ChatMessage, error)
	AddMessage(message *models.ChatMessage) error
	GetAllSessions(userID uuid.UUID, limit, offset int) ([]dto.ChatSessionResponse, error)
	GetSessionMessages(sessionID, userID uuid.UUID) ([]dto.ChatMessageResponse, error)
}

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db: db}
}

// GetSession returns a session by ID and UserID
func (r *chatRepository) GetSession(sessionID, userID uuid.UUID) (*models.ChatSession, error) {
	var session models.ChatSession
	err := r.db.Where("id = ? AND user_id = ?", sessionID, userID).First(&session).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

// CreateSession creates a new chat session
func (r *chatRepository) CreateSession(session *models.ChatSession) error {
	return r.db.Create(session).Error
}

// GetRecentMessages gets the last N messages for a session
func (r *chatRepository) GetRecentMessages(sessionID uuid.UUID, limit int) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := r.db.Where("session_id = ?", sessionID).Order("created_at asc").Limit(limit).Find(&messages).Error
	return messages, err
}

// AddMessage adds a new message to a session
func (r *chatRepository) AddMessage(message *models.ChatMessage) error {
	return r.db.Create(message).Error
}

// GetAllSessions gets all chat sessions for a user, with message counts and last active times
func (r *chatRepository) GetAllSessions(userID uuid.UUID, limit, offset int) ([]dto.ChatSessionResponse, error) {
	var results []dto.ChatSessionResponse

	// We use raw SQL or Query builder to get count and max(created_at) of messages
	query := `
		SELECT 
			s.id, 
			s.user_id,
			s.created_at, 
			COALESCE(COUNT(m.id), 0) as message_count, 
			COALESCE(MAX(m.created_at), s.created_at) as last_active
		FROM chat_sessions s
		LEFT JOIN chat_messages m ON s.id = m.session_id
		WHERE s.user_id = ?
		GROUP BY s.id
		ORDER BY last_active DESC
		LIMIT ? OFFSET ?
	`
	err := r.db.Raw(query, userID, limit, offset).Scan(&results).Error
	return results, err
}

// GetSessionMessages returns formatted messages for a specific session
func (r *chatRepository) GetSessionMessages(sessionID, userID uuid.UUID) ([]dto.ChatMessageResponse, error) {
	// First verify ownership
	session, err := r.GetSession(sessionID, userID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, gorm.ErrRecordNotFound
	}

	var messages []models.ChatMessage
	err = r.db.Where("session_id = ?", sessionID).Order("created_at asc").Find(&messages).Error
	if err != nil {
		return nil, err
	}

	var dtos []dto.ChatMessageResponse
	for _, m := range messages {
		// Map 'sender' for backward frontend compatibility if needed
		sender := m.Sender
		if sender == "assistant" {
			sender = "bot"
		}
		dtos = append(dtos, dto.ChatMessageResponse{
			ID:        m.ID,
			Sender:    sender,
			Content:   m.Content,
			CreatedAt: m.CreatedAt,
		})
	}
	return dtos, nil
}
