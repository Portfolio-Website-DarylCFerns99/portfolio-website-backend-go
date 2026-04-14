package repository

import (
	"portfolio-website-backend/internal/models"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VectorRepository struct {
	db *gorm.DB
}

func NewVectorRepository(db *gorm.DB) *VectorRepository {
	return &VectorRepository{db: db}
}

// ClearAllVectors deletes all existing vectors for a given user
func (r *VectorRepository) ClearAllVectors(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.VectorEmbedding{}).Error
}

// AddEmbedding adds a new vector embedding record
func (r *VectorRepository) AddEmbedding(embedding *models.VectorEmbedding) error {
	return r.db.Create(embedding).Error
}

// Search searches for similar vectors filtered by user and optionally source types
func (r *VectorRepository) Search(userID uuid.UUID, queryVector []float32, limit int, filters []string) ([]models.VectorEmbedding, error) {
	var results []models.VectorEmbedding
	
	pgvec := pgvector.NewVector(queryVector)
	
	query := r.db.Where("user_id = ?", userID)
	
	if len(filters) > 0 {
		query = query.Where("source_type IN ?", filters)
	}
	
	// pgvector L2 distance operator is <->
	err := query.Clauses(clause.OrderBy{Expression: clause.Expr{SQL: "embedding <-> ?", Vars: []interface{}{pgvec}}}).Limit(limit).Find(&results).Error
	
	return results, err
}
