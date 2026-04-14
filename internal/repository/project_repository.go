package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"portfolio-website-backend/internal/models"
)

type ProjectRepository interface {
	Create(project *models.Project) (*models.Project, error)
	GetAll(userID uuid.UUID, skip, limit int) ([]models.Project, error)
	GetVisible(userID uuid.UUID, skip, limit int) ([]models.Project, error)
	GetByID(userID uuid.UUID, id uuid.UUID) (*models.Project, error)
	GetVisibleByID(userID uuid.UUID, id uuid.UUID) (*models.Project, error)
	Update(userID uuid.UUID, id uuid.UUID, data map[string]interface{}) (*models.Project, error)
	Delete(userID uuid.UUID, id uuid.UUID) error
	Count(userID uuid.UUID) (int64, error)
	CountVisible(userID uuid.UUID) (int64, error)
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(project *models.Project) (*models.Project, error) {
	if err := r.db.Create(project).Error; err != nil {
		return nil, err
	}
	return project, nil
}

func (r *projectRepository) GetAll(userID uuid.UUID, skip, limit int) ([]models.Project, error) {
	var projects []models.Project
	if err := r.db.Where("user_id = ?", userID).Offset(skip).Limit(limit).Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepository) GetVisible(userID uuid.UUID, skip, limit int) ([]models.Project, error) {
	var projects []models.Project
	if err := r.db.Where("user_id = ? AND is_visible = ?", userID, true).Offset(skip).Limit(limit).Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepository) GetByID(userID uuid.UUID, id uuid.UUID) (*models.Project, error) {
	var project models.Project
	if err := r.db.Where("user_id = ?", userID).First(&project, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) GetVisibleByID(userID uuid.UUID, id uuid.UUID) (*models.Project, error) {
	var project models.Project
	if err := r.db.Where("user_id = ? AND is_visible = ?", userID, true).First(&project, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) Update(userID uuid.UUID, id uuid.UUID, data map[string]interface{}) (*models.Project, error) {
	var project models.Project
	if err := r.db.Where("user_id = ?", userID).First(&project, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Normalize JSONB fields: JSON deserialization yields []interface{} for arrays,
	// but Postgres needs the custom driver.Valuer (JSONStringArray) to encode them
	// as proper JSON bytes. Without this, GORM sends a Go record instead of jsonb.
	if raw, ok := data["tags"]; ok {
		switch v := raw.(type) {
		case []interface{}:
			tags := make(models.JSONStringArray, 0, len(v))
			for _, item := range v {
				if s, ok := item.(string); ok {
					tags = append(tags, s)
				}
			}
			data["tags"] = tags
		case []string:
			data["tags"] = models.JSONStringArray(v)
		}
	}

	if err := r.db.Model(&project).Updates(data).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) Delete(userID uuid.UUID, id uuid.UUID) error {
	result := r.db.Where("user_id = ?", userID).Delete(&models.Project{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}
	return nil
}

func (r *projectRepository) Count(userID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Project{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *projectRepository) CountVisible(userID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Project{}).Where("user_id = ? AND is_visible = ?", userID, true).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
