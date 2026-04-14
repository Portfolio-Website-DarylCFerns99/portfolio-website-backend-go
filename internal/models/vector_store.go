package models

import (
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

type VectorEmbedding struct {
	BaseModel

	Content      string          `gorm:"type:text;not null" json:"content"`
	Embedding    pgvector.Vector `gorm:"type:vector(768);not null" json:"embedding"`
	MetadataJSON JSONMap         `gorm:"type:jsonb" json:"metadata_json,omitempty"`
	SourceType   string          `gorm:"type:varchar(50);not null;index" json:"source_type"`
	SourceID     *uuid.UUID      `gorm:"type:uuid" json:"source_id,omitempty"`
	UserID       uuid.UUID       `gorm:"type:uuid;not null;index" json:"user_id"`
	User         *User           `gorm:"foreignKey:UserID;references:ID" json:"-"`
}
