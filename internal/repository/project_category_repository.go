package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"portfolio-website-backend/internal/models"
)

type ProjectCategoryRepository interface {
	Create(category *models.ProjectCategory) (*models.ProjectCategory, error)
	GetAll(userID uuid.UUID, skip, limit int) ([]models.ProjectCategory, int64, error)
	GetVisible(userID uuid.UUID, skip, limit int) ([]models.ProjectCategory, int64, error)
	GetByID(userID uuid.UUID, id uuid.UUID) (*models.ProjectCategory, error)
	GetVisibleByID(userID uuid.UUID, id uuid.UUID) (*models.ProjectCategory, error)
	Update(userID uuid.UUID, id uuid.UUID, data map[string]interface{}) (*models.ProjectCategory, error)
	Delete(userID uuid.UUID, id uuid.UUID) error
}

type projectCategoryRepository struct {
	db *gorm.DB
}

func NewProjectCategoryRepository(db *gorm.DB) ProjectCategoryRepository {
	return &projectCategoryRepository{db: db}
}

func (r *projectCategoryRepository) Create(category *models.ProjectCategory) (*models.ProjectCategory, error) {
	if err := r.db.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (r *projectCategoryRepository) GetAll(userID uuid.UUID, skip, limit int) ([]models.ProjectCategory, int64, error) {
	var categories []models.ProjectCategory
	var total int64

	query := r.db.Model(&models.ProjectCategory{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(skip).Limit(limit).Find(&categories).Error; err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

func (r *projectCategoryRepository) GetVisible(userID uuid.UUID, skip, limit int) ([]models.ProjectCategory, int64, error) {
	var categories []models.ProjectCategory
	var total int64

	query := r.db.Model(&models.ProjectCategory{}).Where("user_id = ? AND is_visible = ?", userID, true)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(skip).Limit(limit).Find(&categories).Error; err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

func (r *projectCategoryRepository) GetByID(userID uuid.UUID, id uuid.UUID) (*models.ProjectCategory, error) {
	var category models.ProjectCategory
	if err := r.db.Where("user_id = ?", userID).First(&category, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *projectCategoryRepository) GetVisibleByID(userID uuid.UUID, id uuid.UUID) (*models.ProjectCategory, error) {
	var category models.ProjectCategory
	if err := r.db.Where("user_id = ? AND is_visible = ?", userID, true).First(&category, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *projectCategoryRepository) Update(userID uuid.UUID, id uuid.UUID, data map[string]interface{}) (*models.ProjectCategory, error) {
	var category models.ProjectCategory
	if err := r.db.Where("user_id = ?", userID).First(&category, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if err := r.db.Model(&category).Updates(data).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *projectCategoryRepository) Delete(userID uuid.UUID, id uuid.UUID) error {
	result := r.db.Where("user_id = ?", userID).Delete(&models.ProjectCategory{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}
	return nil
}
