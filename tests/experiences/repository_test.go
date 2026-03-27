package experiences

import (
	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/repository"
	"portfolio-website-backend/tests/common"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestExperienceRepo_CreateAndGet(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewExperienceRepository(db)

	userID := uuid.New()
	exp := &models.Experience{
		Title:        "Software Engineer",
		Organization: "Tech Corp",
		Type:         "experience",
		IsVisible:    true,
		UserID:       userID,
	}

	created, err := repo.Create(exp)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, created.ID)
	assert.Equal(t, "Software Engineer", created.Title)

	fetched, err := repo.GetByID(userID, created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
}

func TestExperienceRepo_GetByTypeAndVisibility(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewExperienceRepository(db)

	userID := uuid.New()
	repo.Create(&models.Experience{Type: "education", IsVisible: true, Title: "BS Computer Science", UserID: userID})
	repo.Create(&models.Experience{Type: "education", IsVisible: false, Title: "High School", UserID: userID})
	repo.Create(&models.Experience{Type: "certification", IsVisible: true, Title: "AWS Cert", UserID: userID})

	// Fetch only visible education
	visibleEdu, err := repo.GetByType(userID, "education", 0, 10, true)
	// fmt.Println(visibleEdu)
	assert.NoError(t, err)
	assert.Len(t, visibleEdu, 1)
	assert.Equal(t, "BS Computer Science", visibleEdu[0].Title)

	// Fetch ALL education
	allEdu, err := repo.GetByType(userID, "education", 0, 10, false)
	assert.NoError(t, err)
	assert.Len(t, allEdu, 2)
}

func TestExperienceRepo_UpdateAndDelete(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewExperienceRepository(db)

	userID := uuid.New()
	created, _ := repo.Create(&models.Experience{Title: "Temporary Role", IsVisible: true, UserID: userID})

	// Update
	updated, err := repo.Update(userID, created.ID, map[string]interface{}{"title": "Permanent Role"})
	assert.NoError(t, err)
	assert.Equal(t, "Permanent Role", updated.Title)

	// Delete
	err = repo.Delete(userID, created.ID)
	assert.NoError(t, err)

	// Verify Missing
	deleted, err := repo.GetByID(userID, created.ID)
	assert.NoError(t, err)
	assert.Nil(t, deleted)
}
