package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"portfolio-website-backend/internal/models"
)

type ReviewRepository interface {
	Create(review *models.Review) (*models.Review, error)
	GetAll(userID uuid.UUID, skip, limit int) ([]models.Review, error)
	GetVisible(userID uuid.UUID, skip, limit int) ([]models.Review, error)
	GetByID(userID uuid.UUID, id uuid.UUID) (*models.Review, error)
	GetVisibleByID(userID uuid.UUID, id uuid.UUID) (*models.Review, error)
	Update(userID uuid.UUID, id uuid.UUID, data map[string]interface{}) (*models.Review, error)
	Delete(userID uuid.UUID, id uuid.UUID) error
	Count(userID uuid.UUID) (int64, error)
	CountVisible(userID uuid.UUID) (int64, error)
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(review *models.Review) (*models.Review, error) {
	if err := r.db.Create(review).Error; err != nil {
		return nil, err
	}
	return review, nil
}

func (r *reviewRepository) GetAll(userID uuid.UUID, skip, limit int) ([]models.Review, error) {
	var reviews []models.Review
	if err := r.db.Where("user_id = ?", userID).Offset(skip).Limit(limit).Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *reviewRepository) GetVisible(userID uuid.UUID, skip, limit int) ([]models.Review, error) {
	var reviews []models.Review
	if err := r.db.Where("user_id = ? AND is_visible = ?", userID, true).Offset(skip).Limit(limit).Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *reviewRepository) GetByID(userID uuid.UUID, id uuid.UUID) (*models.Review, error) {
	var review models.Review
	if err := r.db.Where("user_id = ?", userID).First(&review, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) GetVisibleByID(userID uuid.UUID, id uuid.UUID) (*models.Review, error) {
	var review models.Review
	if err := r.db.Where("user_id = ? AND is_visible = ?", userID, true).First(&review, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) Update(userID uuid.UUID, id uuid.UUID, data map[string]interface{}) (*models.Review, error) {
	var review models.Review
	if err := r.db.Where("user_id = ?", userID).First(&review, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if err := r.db.Model(&review).Updates(data).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) Delete(userID uuid.UUID, id uuid.UUID) error {
	result := r.db.Where("user_id = ?", userID).Delete(&models.Review{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}
	return nil
}

func (r *reviewRepository) Count(userID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Review{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *reviewRepository) CountVisible(userID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Review{}).Where("user_id = ? AND is_visible = ?", userID, true).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
