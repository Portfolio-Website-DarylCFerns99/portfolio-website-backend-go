package services

import (
	"errors"

	"github.com/google/uuid"

	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/repository"
)

type ReviewService interface {
	CreateReview(review *models.Review) (*models.Review, error)
	GetReviews(userID uuid.UUID, skip, limit int, publicOnly bool) ([]models.Review, int64, error)
	GetReviewByID(userID uuid.UUID, id uuid.UUID, publicOnly bool) (*models.Review, error)
	UpdateReview(userID uuid.UUID, id uuid.UUID, data map[string]interface{}) (*models.Review, error)
	UpdateReviewVisibility(userID uuid.UUID, id uuid.UUID, isVisible bool) (*models.Review, error)
	DeleteReview(userID uuid.UUID, id uuid.UUID) error
}

type reviewService struct {
	repo repository.ReviewRepository
}

func NewReviewService(repo repository.ReviewRepository) ReviewService {
	return &reviewService{repo: repo}
}

func (s *reviewService) CreateReview(review *models.Review) (*models.Review, error) {
	return s.repo.Create(review)
}

func (s *reviewService) GetReviews(userID uuid.UUID, skip, limit int, publicOnly bool) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64
	var err error

	if publicOnly {
		reviews, err = s.repo.GetVisible(userID, skip, limit)
		if err == nil {
			total, err = s.repo.CountVisible(userID)
		}
	} else {
		reviews, err = s.repo.GetAll(userID, skip, limit)
		if err == nil {
			total, err = s.repo.Count(userID)
		}
	}

	if err != nil {
		return nil, 0, err
	}
	return reviews, total, nil
}

func (s *reviewService) GetReviewByID(userID uuid.UUID, id uuid.UUID, publicOnly bool) (*models.Review, error) {
	if publicOnly {
		return s.repo.GetVisibleByID(userID, id)
	}
	return s.repo.GetByID(userID, id)
}

func (s *reviewService) UpdateReview(userID uuid.UUID, id uuid.UUID, data map[string]interface{}) (*models.Review, error) {
	updated, err := s.repo.Update(userID, id, data)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, errors.New("review not found")
	}
	return updated, nil
}

func (s *reviewService) UpdateReviewVisibility(userID uuid.UUID, id uuid.UUID, isVisible bool) (*models.Review, error) {
	data := map[string]interface{}{"is_visible": isVisible}
	return s.UpdateReview(userID, id, data)
}

func (s *reviewService) DeleteReview(userID uuid.UUID, id uuid.UUID) error {
	return s.repo.Delete(userID, id)
}
