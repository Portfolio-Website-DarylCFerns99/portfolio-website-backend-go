package reviews

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/services"
)

func TestReviewService_CreateReview(t *testing.T) {
	mockRepo := new(MockReviewRepository)
	svc := services.NewReviewService(mockRepo)

	review := &models.Review{Content: "Great!"}
	mockRepo.On("Create", review).Return(review, nil)

	created, err := svc.CreateReview(review)

	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "Great!", created.Content)
	mockRepo.AssertExpectations(t)
}

func TestReviewService_GetReviews(t *testing.T) {
	mockRepo := new(MockReviewRepository)
	svc := services.NewReviewService(mockRepo)

	userID := uuid.New()

	// Test public only
	mockRepo.On("GetVisible", userID, 0, 10).Return([]models.Review{{Content: "Public"}}, nil)
	mockRepo.On("CountVisible", userID).Return(int64(1), nil)

	reviews, total, err := svc.GetReviews(userID, 0, 10, true)
	assert.NoError(t, err)
	assert.Len(t, reviews, 1)
	assert.Equal(t, int64(1), total)

	// Test all
	mockRepo.On("GetAll", userID, 0, 10).Return([]models.Review{{Content: "All"}}, nil)
	mockRepo.On("Count", userID).Return(int64(2), nil)

	reviewsAll, totalAll, err := svc.GetReviews(userID, 0, 10, false)
	assert.NoError(t, err)
	assert.Len(t, reviewsAll, 1)
	assert.Equal(t, int64(2), totalAll)
}

func TestReviewService_UpdateReview_NotFound(t *testing.T) {
	mockRepo := new(MockReviewRepository)
	svc := services.NewReviewService(mockRepo)

	userID := uuid.New()
	id := uuid.New()

	mockRepo.On("Update", userID, id, map[string]interface{}{"content": "updated"}).Return((*models.Review)(nil), errors.New("not found"))

	updated, err := svc.UpdateReview(userID, id, map[string]interface{}{"content": "updated"})
	assert.Error(t, err)
	assert.Nil(t, updated)
}
