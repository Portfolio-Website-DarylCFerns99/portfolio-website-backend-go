package models

type ProjectCategory struct {
	BaseModel

	Name      string `gorm:"type:varchar(100);not null;index" json:"name"`
	IsVisible bool   `gorm:"default:true;not null" json:"is_visible"`
}
