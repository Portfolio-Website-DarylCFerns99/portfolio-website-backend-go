package models

import "github.com/google/uuid"

type ChatSession struct {
	BaseModel

	UserID   uuid.UUID     `gorm:"type:uuid;not null;index" json:"user_id"`
	User     *User         `gorm:"foreignKey:UserID;references:ID" json:"-"`
	Title    *string       `gorm:"type:varchar(200)" json:"title,omitempty"`
	Messages []ChatMessage `gorm:"foreignKey:SessionID" json:"messages,omitempty"`
}

type ChatMessage struct {
	BaseModel

	SessionID uuid.UUID `gorm:"type:uuid;not null;index" json:"session_id"`
	Sender    string    `gorm:"type:varchar(20);not null;column:sender" json:"sender"` // legacy db column to satisfy constraint
	Role      string    `gorm:"type:varchar(20);not null;column:role" json:"role"`     // newer db column to satisfy constraint
	Content   string    `gorm:"type:text;not null" json:"content"`
}
