package models

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	BaseModel

	Title            string     `gorm:"type:varchar(100);not null;index" json:"title"`
	Description      *string    `gorm:"type:text" json:"description,omitempty"`
	Type             string     `gorm:"type:varchar(50);not null" json:"type"` // "github" or "custom"
	Image            *string    `gorm:"type:text" json:"image,omitempty"`
	Tags             JSONStringArray `gorm:"type:jsonb;default:'[]'" json:"tags"`
	URL              *string    `gorm:"type:varchar(255)" json:"url,omitempty"`
	AdditionalData   JSONMap    `gorm:"type:jsonb" json:"additional_data,omitempty"`
	ExpiryDate       *time.Time `json:"expiry_date,omitempty"`
	IsVisible        bool       `gorm:"default:true" json:"is_visible"`

	ProjectCategoryID *uuid.UUID `gorm:"type:uuid" json:"project_category_id,omitempty"`
	// Relationship can be added if needed: Category ProjectCategory `gorm:"foreignKey:ProjectCategoryID"`
}
