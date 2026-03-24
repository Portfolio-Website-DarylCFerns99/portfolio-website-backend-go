package models

import (
	"github.com/google/uuid"
)

type SkillGroup struct {
	BaseModel

	Name      string  `gorm:"type:varchar(100);not null;index" json:"name"`
	IsVisible bool    `gorm:"default:true;not null" json:"is_visible"`
	Skills    []Skill `gorm:"foreignKey:SkillGroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"skills,omitempty"`
}

type Skill struct {
	BaseModel

	Name        string `gorm:"type:varchar(100);not null;index" json:"name"`
	Proficiency int    `gorm:"type:int;not null" json:"proficiency"` // 1-5 rating
	Color       *string `gorm:"type:varchar(20)" json:"color,omitempty"`
	Icon        *string `gorm:"type:varchar(50)" json:"icon,omitempty"`

	SkillGroupID uuid.UUID   `gorm:"type:uuid;not null" json:"skill_group_id"`
	SkillGroup   *SkillGroup `gorm:"foreignKey:SkillGroupID" json:"skill_group,omitempty"`
	IsVisible    bool        `gorm:"default:true" json:"is_visible"`
}
