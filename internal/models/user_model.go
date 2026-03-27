package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONStringArray handles string arrays stored as JSON in Postgres
type JSONStringArray []string

func (j *JSONStringArray) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONStringArray, 0)
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, j)
	case string:
		return json.Unmarshal([]byte(v), j)
	default:
		return errors.New("type assertion to []byte failed for JSONStringArray")
	}
}

func (j JSONStringArray) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "[]", nil
	}
	return json.Marshal(j)
}

// JSONMap handles generic JSON objects stored in Postgres
type JSONMap map[string]interface{}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONMap)
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, j)
	case string:
		return json.Unmarshal([]byte(v), j)
	default:
		return errors.New("type assertion to []byte failed for JSONMap")
	}
}

func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return "{}", nil
	}
	return json.Marshal(j)
}

type User struct {
	BaseModel

	// Required fields
	Username       string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email          string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	HashedPassword string `gorm:"type:varchar(255);not null" json:"-"` // Hidden from JSON

	// Optional fields
	Name             *string         `gorm:"type:varchar(100)" json:"name,omitempty"`
	Surname          *string         `gorm:"type:varchar(100)" json:"surname,omitempty"`
	Title            *string         `gorm:"type:varchar(100)" json:"title,omitempty"`
	Phone            *string         `gorm:"type:varchar(20)" json:"phone,omitempty"`
	Location         *string         `gorm:"type:varchar(100)" json:"location,omitempty"`
	Availability     *string         `gorm:"type:varchar(50)" json:"availability,omitempty"`
	Avatar           *string         `gorm:"type:text" json:"avatar,omitempty"`
	SocialLinks      JSONMap         `gorm:"type:jsonb" json:"social_links,omitempty"`
	About            JSONMap         `gorm:"type:jsonb" json:"about,omitempty"`
	FeaturedSkillIDs JSONStringArray `gorm:"type:jsonb;default:'[]'" json:"featured_skill_ids"`
}
