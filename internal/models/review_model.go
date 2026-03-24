package models

type Review struct {
	BaseModel

	Name           string `gorm:"type:varchar(100);not null" json:"name"`
	Content        string `gorm:"type:text;not null" json:"content"`
	Rating         int    `gorm:"type:int;not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	WhereKnownFrom *string `gorm:"type:varchar(200)" json:"where_known_from,omitempty"`
	IsVisible      bool   `gorm:"default:false" json:"is_visible"`
}
