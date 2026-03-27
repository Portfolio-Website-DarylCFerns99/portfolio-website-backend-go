package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"portfolio-website-backend/internal/models"
)

type ExperienceRepository interface {
	Create(experience *models.Experience) (*models.Experience, error)
	GetAll(userID uuid.UUID, skip, limit int) ([]models.Experience, error)
	GetVisible(userID uuid.UUID, skip, limit int) ([]models.Experience, error)
	GetByType(userID uuid.UUID, expType string, skip, limit int, onlyVisible bool) ([]models.Experience, error)
	GetByID(userID uuid.UUID, id uuid.UUID) (*models.Experience, error)
	GetVisibleByID(userID uuid.UUID, id uuid.UUID) (*models.Experience, error)
	Update(userID uuid.UUID, id uuid.UUID, data map[string]interface{}) (*models.Experience, error)
	Delete(userID uuid.UUID, id uuid.UUID) error
	Count(userID uuid.UUID) (int64, error)
	CountVisible(userID uuid.UUID) (int64, error)
	CountByType(userID uuid.UUID, expType string, onlyVisible bool) (int64, error)
}

type experienceRepository struct {
	db *gorm.DB
}

func NewExperienceRepository(db *gorm.DB) ExperienceRepository {
	return &experienceRepository{db: db}
}

func (r *experienceRepository) Create(experience *models.Experience) (*models.Experience, error) {
	if err := r.db.Create(experience).Error; err != nil {
		return nil, err
	}
	return experience, nil
}

func (r *experienceRepository) GetAll(userID uuid.UUID, skip, limit int) ([]models.Experience, error) {
	var exp []models.Experience
	if err := r.db.Where("user_id = ?", userID).Offset(skip).Limit(limit).Find(&exp).Error; err != nil {
		return nil, err
	}
	return exp, nil
}

func (r *experienceRepository) GetVisible(userID uuid.UUID, skip, limit int) ([]models.Experience, error) {
	var exp []models.Experience
	if err := r.db.Where("user_id = ? AND is_visible = ?", userID, true).Offset(skip).Limit(limit).Find(&exp).Error; err != nil {
		return nil, err
	}
	return exp, nil
}

func (r *experienceRepository) GetByType(userID uuid.UUID, expType string, skip, limit int, onlyVisible bool) ([]models.Experience, error) {
	var exp []models.Experience
	q := r.db.Where("user_id = ? AND type = ?", userID, expType)
	if onlyVisible {
		q = q.Where("is_visible = ?", true)
	}
	if err := q.Offset(skip).Limit(limit).Find(&exp).Error; err != nil {
		return nil, err
	}
	return exp, nil
}

func (r *experienceRepository) GetByID(userID uuid.UUID, id uuid.UUID) (*models.Experience, error) {
	var exp models.Experience
	if err := r.db.Where("user_id = ?", userID).First(&exp, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found, mirroring Python
		}
		return nil, err
	}
	return &exp, nil
}

func (r *experienceRepository) GetVisibleByID(userID uuid.UUID, id uuid.UUID) (*models.Experience, error) {
	var exp models.Experience
	if err := r.db.Where("user_id = ? AND is_visible = ?", userID, true).First(&exp, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &exp, nil
}

func (r *experienceRepository) Update(userID uuid.UUID, id uuid.UUID, data map[string]interface{}) (*models.Experience, error) {
	var exp models.Experience
	if err := r.db.Where("user_id = ?", userID).First(&exp, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if err := r.db.Model(&exp).Updates(data).Error; err != nil {
		return nil, err
	}
	return &exp, nil
}

func (r *experienceRepository) Delete(userID uuid.UUID, id uuid.UUID) error {
	result := r.db.Where("user_id = ?", userID).Delete(&models.Experience{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}
	return nil
}

func (r *experienceRepository) Count(userID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Experience{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *experienceRepository) CountVisible(userID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Experience{}).Where("user_id = ? AND is_visible = ?", userID, true).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *experienceRepository) CountByType(userID uuid.UUID, expType string, onlyVisible bool) (int64, error) {
	var count int64
	q := r.db.Model(&models.Experience{}).Where("user_id = ? AND type = ?", userID, expType)
	if onlyVisible {
		q = q.Where("is_visible = ?", true)
	}
	if err := q.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
