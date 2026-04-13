package project_categories

import (
	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/repository"
	"portfolio-website-backend/internal/services"
	"portfolio-website-backend/tests/common"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProjectCategoryService_CreateAndRetrieve(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewProjectCategoryRepository(db)
	svc := services.NewProjectCategoryService(repo)

	userID := uuid.New()
	common.CreateTestUser(db, userID)
	created, err := svc.CreateCategory(&models.ProjectCategory{
		Name:      "Service Category",
		IsVisible: false,
		UserID:    userID,
	})
	assert.NoError(t, err)

	// Test Visibility Filter from Service
	visibleFetch, _ := svc.GetCategoryByID(userID, created.ID, true)
	assert.Nil(t, visibleFetch)

	allFetch, _ := svc.GetCategoryByID(userID, created.ID, false)
	assert.NotNil(t, allFetch)
	assert.Equal(t, "Service Category", allFetch.Name)
}

func TestProjectCategoryService_UpdateVisibility(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewProjectCategoryRepository(db)
	svc := services.NewProjectCategoryService(repo)

	userID := uuid.New()
	common.CreateTestUser(db, userID)
	created, _ := svc.CreateCategory(&models.ProjectCategory{Name: "Hidden Category", IsVisible: false, UserID: userID})

	updated, err := svc.UpdateCategory(userID, created.ID, map[string]interface{}{"is_visible": true})
	assert.NoError(t, err)
	assert.True(t, updated.IsVisible)

	visibleFetch, _ := svc.GetCategoryByID(userID, created.ID, true)
	assert.NotNil(t, visibleFetch)
}
