package dto

import (
	"time"

	"github.com/google/uuid"
)

// ChatSessionResponse represents a chat session list item
type ChatSessionResponse struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	MessageCount int       `json:"message_count"`
	LastActive   time.Time `json:"last_active"`
}

// ChatMessageResponse represents a single message in a session
type ChatMessageResponse struct {
	ID        uuid.UUID `json:"id"`
	Sender    string    `json:"sender"` // mapped from 'role'
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// WsMessage represents messages exchanged over WebSocket connection
type WsMessage struct {
	Type    string      `json:"type"`              // 'history', 'content', 'end', 'error'
	Payload interface{} `json:"payload,omitempty"` // For history: []ChatMessageResponse, for content: string
}
