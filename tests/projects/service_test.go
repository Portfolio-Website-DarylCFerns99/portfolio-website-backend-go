package projects

import (
	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/repository"
	"portfolio-website-backend/internal/services"
	"portfolio-website-backend/tests/common"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProjectService_CreateAndRetrieve(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewProjectRepository(db)
	svc := services.NewProjectService(repo)

	userID := uuid.New()
	common.CreateTestUser(db, userID)
	created, err := svc.CreateProject(&models.Project{
		Title:     "Custom Service Project",
		Type:      "custom",
		IsVisible: false,
		UserID:    userID,
	})
	assert.NoError(t, err)

	// Test Visibility Filter from Service
	visibleFetch, _ := svc.GetProjectByID(userID, created.ID, true)
	assert.Nil(t, visibleFetch)

	allFetch, _ := svc.GetProjectByID(userID, created.ID, false)
	assert.NotNil(t, allFetch)
	assert.Equal(t, "Custom Service Project", allFetch.Title)
}

func TestProjectService_UpdateVisibility(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewProjectRepository(db)
	svc := services.NewProjectService(repo)

	userID := uuid.New()
	common.CreateTestUser(db, userID)
	created, _ := svc.CreateProject(&models.Project{Title: "Hidden Project", Type: "custom", IsVisible: false, UserID: userID})

	updated, err := svc.UpdateProjectVisibility(userID, created.ID, true)
	assert.NoError(t, err)
	assert.True(t, updated.IsVisible)

	visibleFetch, _ := svc.GetProjectByID(userID, created.ID, true)
	assert.NotNil(t, visibleFetch)
}
