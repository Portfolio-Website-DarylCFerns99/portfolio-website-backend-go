package models

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	BaseModel

	ProjectCategoryID *uuid.UUID       `gorm:"type:uuid;index" json:"project_category_id,omitempty"`
	ProjectCategory   *ProjectCategory `gorm:"foreignKey:ProjectCategoryID" json:"project_category,omitempty"`
	UserID            uuid.UUID        `gorm:"type:uuid;not null;index" json:"user_id"`
	Type              string           `gorm:"type:varchar(50);not null" json:"type"` // e.g. 'github', 'custom'
	Title             string           `gorm:"type:varchar(100);not null;index" json:"title"`
	Description       *string          `gorm:"type:text" json:"description,omitempty"`
	Image             *string          `gorm:"type:text" json:"image,omitempty"` // Base64 image data
	Tags              JSONStringArray  `gorm:"type:jsonb;default:'[]'" json:"tags"`
	URL               *string          `gorm:"type:varchar(255)" json:"url,omitempty"`
	AdditionalData    JSONMap          `gorm:"type:jsonb" json:"additional_data,omitempty"`   // Store complete GitHub API response
	ExpiryDate        *time.Time       `gorm:"type:timestamptz" json:"expiry_date,omitempty"` // Expiry date for non-custom projects
	IsVisible         bool             `gorm:"default:false" json:"is_visible"`
	PublishedAt       *time.Time       `gorm:"type:timestamptz" json:"published_at,omitempty"`
}
