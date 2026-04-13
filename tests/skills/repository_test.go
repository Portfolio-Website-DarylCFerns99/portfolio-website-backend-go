package skills

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/repository"
	"portfolio-website-backend/tests/common"
)

func TestSkillRepo_SkillGroup_CreateAndGet(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewSkillRepository(db)

	userID := uuid.New()
	common.CreateTestUser(db, userID)

	group := &models.SkillGroup{
		Name:      "Languages",
		IsVisible: true,
		UserID:    userID,
	}

	created, err := repo.CreateSkillGroup(group)
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "Languages", created.Name)

	groups, err := repo.GetSkillGroups(userID, 0, 10, false)
	assert.NoError(t, err)
	assert.Len(t, groups, 1)

	count, err := repo.CountSkillGroups(userID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestSkillRepo_Skill_CreateAndGet(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewSkillRepository(db)

	userID := uuid.New()
	common.CreateTestUser(db, userID)

	group := &models.SkillGroup{
		Name:      "Languages",
		IsVisible: true,
		UserID:    userID,
	}
	createdGroup, _ := repo.CreateSkillGroup(group)

	skill := &models.Skill{
		Name:         "Go",
		Proficiency:  5,
		IsVisible:    true,
		UserID:       userID,
		SkillGroupID: createdGroup.ID,
	}

	created, err := repo.CreateSkill(skill)
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "Go", created.Name)

	skills, err := repo.GetSkills(userID, 0, 10)
	assert.NoError(t, err)
	assert.Len(t, skills, 1)
}

func TestSkillRepo_UpdateAndDeleteSkillGroup(t *testing.T) {
	db := common.SetupTestDB()
	repo := repository.NewSkillRepository(db)

	userID := uuid.New()
	common.CreateTestUser(db, userID)
	created, _ := repo.CreateSkillGroup(&models.SkillGroup{Name: "Frameworks", UserID: userID})

	updated, err := repo.UpdateSkillGroup(userID, created.ID, map[string]interface{}{"name": "Web Frameworks"})
	assert.NoError(t, err)
	assert.Equal(t, "Web Frameworks", updated.Name)

	err = repo.DeleteSkillGroup(userID, created.ID)
	assert.NoError(t, err)

	deleted, err := repo.GetSkillGroupByID(userID, created.ID, false)
	assert.NoError(t, err)
	assert.Nil(t, deleted)
}

