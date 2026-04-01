package projects

import (
	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/repository"
	"portfolio-website-backend/tests/common"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProjectRepo_CreateAndGet(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewProjectRepository(db)

	userID := uuid.New()
	proj := &models.Project{
		Title:     "Awesome Go Tool",
		Type:      "custom",
		IsVisible: true,
		UserID:    userID,
	}

	created, err := repo.Create(proj)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, created.ID)
	assert.Equal(t, "Awesome Go Tool", created.Title)

	fetched, err := repo.GetByID(userID, created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
}

func TestProjectRepo_GetVisible(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewProjectRepository(db)

	userID := uuid.New()
	repo.Create(&models.Project{Type: "custom", IsVisible: true, Title: "Visible Project", UserID: userID})
	repo.Create(&models.Project{Type: "custom", IsVisible: false, Title: "Hidden Project", UserID: userID})

	// Fetch only visible
	visibleProjs, err := repo.GetVisible(userID, 0, 10)
	assert.NoError(t, err)
	assert.Len(t, visibleProjs, 1)
	assert.Equal(t, "Visible Project", visibleProjs[0].Title)

	// Fetch ALL
	allProjs, err := repo.GetAll(userID, 0, 10)
	assert.NoError(t, err)
	assert.Len(t, allProjs, 2)
}

func TestProjectRepo_UpdateAndDelete(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewProjectRepository(db)

	userID := uuid.New()
	created, _ := repo.Create(&models.Project{Title: "Old Name", IsVisible: true, Type: "custom", UserID: userID})

	// Update
	updated, err := repo.Update(userID, created.ID, map[string]interface{}{"title": "New Name"})
	assert.NoError(t, err)
	assert.Equal(t, "New Name", updated.Title)

	// Delete
	err = repo.Delete(userID, created.ID)
	assert.NoError(t, err)

	// Verify Missing
	deleted, err := repo.GetByID(userID, created.ID)
	assert.NoError(t, err)
	assert.Nil(t, deleted)
}
