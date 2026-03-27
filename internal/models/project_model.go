package models

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	BaseModel

	ProjectCategoryID *uuid.UUID             `gorm:"type:uuid;index" json:"project_category_id,omitempty"`
	Type              string                 `gorm:"type:varchar(50);not null" json:"type"` // e.g. 'personal', 'professional'
	Title             string                 `gorm:"type:varchar(200);not null" json:"title"`
	Description       *string                `gorm:"type:text" json:"description,omitempty"`
	Image             *string                `gorm:"type:text" json:"image,omitempty"`
	Tags              JSONStringArray        `gorm:"type:jsonb;default:'[]'" json:"tags"`
	URL               *string                `gorm:"type:text" json:"url,omitempty"`
	AdditionalData    JSONMap                `gorm:"type:jsonb" json:"additional_data,omitempty"`
	IsVisible         bool                   `gorm:"default:false" json:"is_visible"`
	PublishedAt       *time.Time             `gorm:"type:timestamptz" json:"published_at,omitempty"`
}
