package models

import (
	"github.com/google/uuid"
)

type ProjectCategory struct {
	BaseModel

	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	User        *User     `gorm:"foreignKey:UserID;references:ID" json:"-"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	IsVisible   bool      `gorm:"default:false" json:"is_visible"`
}
