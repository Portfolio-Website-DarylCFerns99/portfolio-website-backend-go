package experiences

import (
	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/repository"
	"portfolio-website-backend/internal/services"
	"portfolio-website-backend/tests/common"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestExperienceService_CreateAndRetrieve(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewExperienceRepository(db)
	svc := services.NewExperienceService(repo)

	userID := uuid.New()
	created, err := svc.CreateExperience(&models.Experience{
		Title:     "Service Role",
		Type:      "experience",
		IsVisible: false,
		UserID:    userID,
	})
	assert.NoError(t, err)

	// Test Visibility Filter from Service
	visibleFetch, _ := svc.GetExperienceByID(userID, created.ID, true)
	assert.Nil(t, visibleFetch) // Should be nil because it's not visible

	allFetch, _ := svc.GetExperienceByID(userID, created.ID, false)
	assert.NotNil(t, allFetch)
	assert.Equal(t, "Service Role", allFetch.Title)
}

func TestExperienceService_UpdateVisibility(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewExperienceRepository(db)
	svc := services.NewExperienceService(repo)

	userID := uuid.New()
	created, _ := svc.CreateExperience(&models.Experience{Title: "Hidden Role", IsVisible: false, UserID: userID})

	// Toggle Visibility
	updated, err := svc.UpdateExperienceVisibility(userID, created.ID, true)
	assert.NoError(t, err)
	assert.True(t, updated.IsVisible)

	// Verify Filter works now
	visibleFetch, _ := svc.GetExperienceByID(userID, created.ID, true)
	assert.NotNil(t, visibleFetch)
}
