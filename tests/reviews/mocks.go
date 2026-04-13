package reviews

import (
	"portfolio-website-backend/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockReviewService struct {
	mock.Mock
}

func (m *MockReviewService) CreateReview(review *models.Review) (*models.Review, error) {
	args := m.Called(review)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Review), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockReviewService) GetReviews(userID uuid.UUID, skip, limit int, publicOnly bool) ([]models.Review, int64, error) {
	args := m.Called(userID, skip, limit, publicOnly)
	var reviews []models.Review
	if args.Get(0) != nil {
		reviews = args.Get(0).([]models.Review)
	}
	var total int64
	if args.Get(1) != nil {
		if t, ok := args.Get(1).(int64); ok {
			total = t
		} else if t, ok := args.Get(1).(int); ok {
			total = int64(t)
		}
	}
	return reviews, total, args.Error(2)
}

func (m *MockReviewService) GetReviewByID(userID uuid.UUID, id uuid.UUID, publicOnly bool) (*models.Review, error) {
	args := m.Called(userID, id, publicOnly)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Review), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockReviewService) UpdateReview(userID uuid.UUID, id uuid.UUID, data map[string]interface{}) (*models.Review, error) {
	args := m.Called(userID, id, data)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Review), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockReviewService) UpdateReviewVisibility(userID uuid.UUID, id uuid.UUID, isVisible bool) (*models.Review, error) {
	args := m.Called(userID, id, isVisible)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Review), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockReviewService) DeleteReview(userID uuid.UUID, id uuid.UUID) error {
	args := m.Called(userID, id)
	return args.Error(0)
}

type MockReviewRepository struct {
	mock.Mock
}

func (m *MockReviewRepository) Create(review *models.Review) (*models.Review, error) {
	args := m.Called(review)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Review), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockReviewRepository) GetAll(userID uuid.UUID, skip, limit int) ([]models.Review, error) {
	args := m.Called(userID, skip, limit)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Review), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockReviewRepository) GetVisible(userID uuid.UUID, skip, limit int) ([]models.Review, error) {
	args := m.Called(userID, skip, limit)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Review), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockReviewRepository) Count(userID uuid.UUID) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockReviewRepository) CountVisible(userID uuid.UUID) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockReviewRepository) GetByID(userID uuid.UUID, id uuid.UUID) (*models.Review, error) {
	args := m.Called(userID, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Review), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockReviewRepository) GetVisibleByID(userID uuid.UUID, id uuid.UUID) (*models.Review, error) {
	args := m.Called(userID, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Review), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockReviewRepository) Update(userID uuid.UUID, id uuid.UUID, data map[string]interface{}) (*models.Review, error) {
	args := m.Called(userID, id, data)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Review), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockReviewRepository) Delete(userID uuid.UUID, id uuid.UUID) error {
	args := m.Called(userID, id)
	return args.Error(0)
}
