package models

import "github.com/google/uuid"

type ChatSession struct {
	BaseModel

	UserID   *uuid.UUID    `gorm:"type:uuid;index" json:"user_id,omitempty"`
	Title    *string       `gorm:"type:varchar(200)" json:"title,omitempty"`
	Messages []ChatMessage `gorm:"foreignKey:SessionID" json:"messages,omitempty"`
}

type ChatMessage struct {
	BaseModel

	SessionID uuid.UUID `gorm:"type:uuid;not null;index" json:"session_id"`
	Role      string    `gorm:"type:varchar(20);not null" json:"role"` // 'user' or 'assistant'
	Content   string    `gorm:"type:text;not null" json:"content"`
}
