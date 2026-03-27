package users_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/repository"
	"portfolio-website-backend/tests/common"
)

func TestRepo_CreateAndGet(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewUserRepository(db)

	email := "repo" + uuid.New().String()[:8] + "@test.com"
	user := &models.User{
		Email:    email,
		Username: "repo_user_" + uuid.New().String()[:8],
	}

	created, err := repo.Create(user)
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.NotEqual(t, uuid.Nil, created.ID)

	fetched, err := repo.GetByID(created.ID)
	assert.NoError(t, err)
	assert.Equal(t, email, fetched.Email)
}

func TestRepo_List(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewUserRepository(db)

	email := "listtest" + uuid.New().String()[:8] + "@test.com"
	repo.Create(&models.User{
		Email:    email,
		Username: "listtest_" + uuid.New().String()[:8],
	})

	users, err := repo.List()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(users), 1)
}

func TestRepo_Update(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewUserRepository(db)

	email := "updatetest" + uuid.New().String()[:8] + "@test.com"
	created, _ := repo.Create(&models.User{
		Email:    email,
		Username: "updatetest_" + uuid.New().String()[:8],
	})

	updates := map[string]interface{}{"title": "Senior Developer"}
	updated, err := repo.Update(created.ID, updates)
	assert.NoError(t, err)
	assert.Equal(t, "Senior Developer", *updated.Title)
}

func TestRepo_GetByEmailorUsername(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewUserRepository(db)

	email := "lookup" + uuid.New().String()[:8] + "@test.com"
	username := "lookup_" + uuid.New().String()[:8]
	repo.Create(&models.User{
		Email:    email,
		Username: username,
	})

	// Lookup by email
	fetchedEmail, err := repo.GetByEmailorUsername(email)
	assert.NoError(t, err)
	assert.Equal(t, username, fetchedEmail.Username)

	// Lookup by username
	fetchedUsername, err := repo.GetByEmailorUsername(username)
	assert.NoError(t, err)
	assert.Equal(t, email, fetchedUsername.Email)
}
