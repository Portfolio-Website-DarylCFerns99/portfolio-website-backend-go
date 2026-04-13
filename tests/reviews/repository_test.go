package reviews

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/repository"
	"portfolio-website-backend/tests/common"
)

func TestReviewRepo_CreateAndGet(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewReviewRepository(db)

	userID := uuid.New()
	common.CreateTestUser(db, userID)
	review := &models.Review{
		Name:      "Jane Doe",
		Content:   "Excellent service!",
		Rating:    5,
		IsVisible: true,
		UserID:    userID,
	}

	created, err := repo.Create(review)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, created.ID)
	assert.Equal(t, "Jane Doe", created.Name)

	fetched, err := repo.GetByID(userID, created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
}

func TestReviewRepo_GetVisibility(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewReviewRepository(db)

	userID := uuid.New()
	common.CreateTestUser(db, userID)
	repo.Create(&models.Review{Name: "Alice", Content: "Visible review", IsVisible: true, UserID: userID})
	repo.Create(&models.Review{Name: "Bob", Content: "Hidden review", IsVisible: false, UserID: userID})

	// Only visible
	visibleReviews, err := repo.GetVisible(userID, 0, 10)
	assert.NoError(t, err)
	assert.Len(t, visibleReviews, 1)
	assert.Equal(t, "Alice", visibleReviews[0].Name)

	// Count visible
	countVisible, err := repo.CountVisible(userID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), countVisible)

	// All
	allReviews, err := repo.GetAll(userID, 0, 10)
	assert.NoError(t, err)
	assert.Len(t, allReviews, 2)

	// Count all
	countAll, err := repo.Count(userID)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), countAll)
}

func TestReviewRepo_UpdateAndDelete(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewReviewRepository(db)

	userID := uuid.New()
	common.CreateTestUser(db, userID)
	created, _ := repo.Create(&models.Review{Name: "Temporary", Content: "Old content", UserID: userID})

	// Update
	updated, err := repo.Update(userID, created.ID, map[string]interface{}{"content": "Updated content"})
	assert.NoError(t, err)
	assert.Equal(t, "Updated content", updated.Content)

	// Delete
	err = repo.Delete(userID, created.ID)
	assert.NoError(t, err)

	// Verify missing
	deleted, err := repo.GetByID(userID, created.ID)
	assert.NoError(t, err)
	assert.Nil(t, deleted)
}
