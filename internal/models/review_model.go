package models

type Review struct {
	BaseModel

	Name      string  `gorm:"type:varchar(100);not null" json:"name"`
	Position  *string `gorm:"type:varchar(100)" json:"position,omitempty"`
	Company   *string `gorm:"type:varchar(100)" json:"company,omitempty"`
	Content   string  `gorm:"type:text;not null" json:"content"`
	Rating    int     `gorm:"default:5" json:"rating"`
	Avatar    *string `gorm:"type:text" json:"avatar,omitempty"`
	IsVisible bool    `gorm:"default:false" json:"is_visible"`
}
