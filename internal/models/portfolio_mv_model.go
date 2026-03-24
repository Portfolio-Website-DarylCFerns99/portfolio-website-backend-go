package models

import (
	"github.com/google/uuid"
)

// PortfolioMV matches the PostgreSQL materialized view structure
type PortfolioMV struct {
	UserID uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`

	Name             *string         `gorm:"type:varchar(100)" json:"name,omitempty"`
	Surname          *string         `gorm:"type:varchar(100)" json:"surname,omitempty"`
	Title            *string         `gorm:"type:varchar(100)" json:"title,omitempty"`
	Email            *string         `gorm:"type:varchar(100)" json:"email,omitempty"`
	Phone            *string         `gorm:"type:varchar(20)" json:"phone,omitempty"`
	Location         *string         `gorm:"type:varchar(100)" json:"location,omitempty"`
	Availability     *string         `gorm:"type:varchar(50)" json:"availability,omitempty"`
	Avatar           *string         `gorm:"type:text" json:"avatar,omitempty"`
	SocialLinks      JSONMap         `gorm:"type:jsonb" json:"social_links,omitempty"`
	About            JSONMap         `gorm:"type:jsonb" json:"about,omitempty"`
	FeaturedSkillIDs JSONStringArray `gorm:"type:jsonb" json:"featured_skill_ids,omitempty"`

	// Aggregated JSON fields from the materialized view
	Experiences       JSONMap `gorm:"type:jsonb" json:"experiences,omitempty"`
	Projects          JSONMap `gorm:"type:jsonb" json:"projects,omitempty"`
	SkillGroups       JSONMap `gorm:"type:jsonb" json:"skill_groups,omitempty"`
	ProjectCategories JSONMap `gorm:"type:jsonb" json:"project_categories,omitempty"`
	Reviews           JSONMap `gorm:"type:jsonb" json:"reviews,omitempty"`
}

// TableName overrides the table name used by GORM 
func (PortfolioMV) TableName() string {
	return "portfolio_mv"
}
