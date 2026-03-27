package users_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/repository"
	"portfolio-website-backend/internal/services"
	"portfolio-website-backend/tests/common"
)

// func TestService_ValidateFeaturedSkills(t *testing.T) {
// 	db := common.SetupTestDB()
// 	repo := repository.NewUserRepository(db)
// 	svc := services.NewUserService(db, repo)

// 	// Create a dummy skill in DB so we can validate true existence
// 	skillID := uuid.New()
// 	db.Create(&models.Skill{
// 		BaseModel: models.BaseModel{ID: skillID},
// 		Name:      "Go Validation Test",
// 	})

// 	// Provide one real ID and one fake ID
// 	valid := svc.ValidateFeaturedSkills(uuid.New(), []string{skillID.String(), uuid.New().String()})

// 	// Should successfully filter out the fake one
// 	assert.Len(t, valid, 1)
// 	assert.Equal(t, skillID.String(), valid[0])
// }

func TestService_GetPublicPortfolioData(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewUserRepository(db)
	svc := services.NewUserService(db, repo)

	email := "portfoliotest" + uuid.New().String()[:8] + "@test.com"
	user, err := repo.Create(&models.User{
		Email:    email,
		Username: "portfoliotest_" + uuid.New().String()[:8],
	})
	assert.NoError(t, err)

	data, err := svc.GetPublicPortfolioData(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, email, data.Email)
}

func TestService_UpdateUserProfile(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewUserRepository(db)
	svc := services.NewUserService(db, repo)

	email := "svc_update" + uuid.New().String()[:8] + "@test.com"
	user, err := repo.Create(&models.User{
		Email:    email,
		Username: "svc_update_" + uuid.New().String()[:8],
	})
	assert.NoError(t, err)

	// Validate generic updates proxy through correctly
	updates := map[string]interface{}{"title": "Principal Engineer"}
	updated, err := svc.UpdateUserProfile(user.ID, updates)

	assert.NoError(t, err)
	assert.Equal(t, "Principal Engineer", *updated.Title)
}
