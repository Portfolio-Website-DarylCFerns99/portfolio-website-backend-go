package models

import (
	"time"
)

type Experience struct {
	BaseModel

	Type         string     `gorm:"type:varchar(50);not null" json:"type"` // 'experience' or 'education'
	Title        string     `gorm:"type:varchar(200);not null" json:"title"`
	Organization string     `gorm:"type:varchar(200);not null" json:"organization"`
	StartDate    time.Time  `gorm:"type:date;not null" json:"start_date"`
	EndDate      *time.Time `gorm:"type:date" json:"end_date,omitempty"`
	Description  *string    `gorm:"type:text" json:"description,omitempty"`
	IsVisible    bool       `gorm:"default:true" json:"is_visible"`
}
