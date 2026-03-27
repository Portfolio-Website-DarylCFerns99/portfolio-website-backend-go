package models

import "github.com/google/uuid"

type SkillGroup struct {
	BaseModel

	Name      string  `gorm:"type:varchar(100);not null" json:"name"`
	IsVisible bool    `gorm:"default:false" json:"is_visible"`
	Skills    []Skill `gorm:"foreignKey:SkillGroupID" json:"skills,omitempty"`
}

type Skill struct {
	BaseModel

	SkillGroupID uuid.UUID `gorm:"type:uuid;not null;index" json:"skill_group_id"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	Proficiency  int       `gorm:"default:0" json:"proficiency"`
	Color        *string   `gorm:"type:varchar(50)" json:"color,omitempty"`
	Icon         *string   `gorm:"type:varchar(100)" json:"icon,omitempty"`
	IsVisible    bool      `gorm:"default:false" json:"is_visible"`
}
