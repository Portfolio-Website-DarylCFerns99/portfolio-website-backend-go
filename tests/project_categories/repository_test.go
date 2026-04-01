package project_categories

import (
	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/repository"
	"portfolio-website-backend/tests/common"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProjectCategoryRepo_CreateAndGet(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewProjectCategoryRepository(db)

	userID := uuid.New()
	cat := &models.ProjectCategory{
		Name:      "Open Source",
		IsVisible: true,
		UserID:    userID,
	}

	created, err := repo.Create(cat)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, created.ID)
	assert.Equal(t, "Open Source", created.Name)

	fetched, err := repo.GetByID(userID, created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
}

func TestProjectCategoryRepo_GetVisible(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewProjectCategoryRepository(db)

	userID := uuid.New()
	repo.Create(&models.ProjectCategory{Name: "Visible Category", IsVisible: true, UserID: userID})
	repo.Create(&models.ProjectCategory{Name: "Hidden Category", IsVisible: false, UserID: userID})

	// Fetch only visible
	visibleCats, err := repo.GetVisible(userID, 0, 10)
	assert.NoError(t, err)
	assert.Len(t, visibleCats, 1)
	assert.Equal(t, "Visible Category", visibleCats[0].Name)

	// Fetch ALL
	allCats, err := repo.GetAll(userID, 0, 10)
	assert.NoError(t, err)
	assert.Len(t, allCats, 2)
}

func TestProjectCategoryRepo_UpdateAndDelete(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewProjectCategoryRepository(db)

	userID := uuid.New()
	created, _ := repo.Create(&models.ProjectCategory{Name: "Old Name", IsVisible: true, UserID: userID})

	// Update
	updated, err := repo.Update(userID, created.ID, map[string]interface{}{"name": "New Name"})
	assert.NoError(t, err)
	assert.Equal(t, "New Name", updated.Name)

	// Delete
	err = repo.Delete(userID, created.ID)
	assert.NoError(t, err)

	// Verify Missing
	deleted, err := repo.GetByID(userID, created.ID)
	assert.NoError(t, err)
	assert.Nil(t, deleted)
}
