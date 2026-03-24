package models

import (
	"github.com/google/uuid"
)

type ChatSession struct {
	BaseModel

	Messages []ChatMessage `gorm:"foreignKey:SessionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"messages,omitempty"`
}

type ChatMessage struct {
	BaseModel

	SessionID uuid.UUID `gorm:"type:uuid;not null;index" json:"session_id"`
	Sender    string    `gorm:"type:varchar(50);not null" json:"sender"` // 'user' or 'bot'
	Content   string    `gorm:"type:text;not null" json:"content"`       

	Session   *ChatSession `gorm:"foreignKey:SessionID" json:"session,omitempty"`
}
