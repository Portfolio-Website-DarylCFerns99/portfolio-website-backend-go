package models

type ProjectCategory struct {
	BaseModel

	Name        string  `gorm:"type:varchar(100);not null" json:"name"`
	Description *string `gorm:"type:text" json:"description,omitempty"`
	IsVisible   bool    `gorm:"default:false" json:"is_visible"`
}
