package skills

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/services"
)

func TestSkillService_CreateSkillGroup(t *testing.T) {
	mockRepo := new(MockSkillRepository)
	svc := services.NewSkillService(mockRepo)

	group := &models.SkillGroup{Name: "Languages"}
	mockRepo.On("CreateSkillGroup", group).Return(group, nil)

	created, err := svc.CreateSkillGroup(group)

	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "Languages", created.Name)
	mockRepo.AssertExpectations(t)
}

func TestSkillService_GetSkillGroups(t *testing.T) {
	mockRepo := new(MockSkillRepository)
	svc := services.NewSkillService(mockRepo)

	userID := uuid.New()

	// Test public only
	mockRepo.On("GetVisibleSkillGroups", userID, 0, 10, true).Return([]models.SkillGroup{{Name: "Visible"}}, nil)
	mockRepo.On("CountVisibleSkillGroups", userID).Return(int64(1), nil)

	groups, total, err := svc.GetSkillGroups(userID, 0, 10, true)
	assert.NoError(t, err)
	assert.Len(t, groups, 1)
	assert.Equal(t, int64(1), total)

	// Test all
	mockRepo.On("GetSkillGroups", userID, 0, 10, true).Return([]models.SkillGroup{{Name: "All"}}, nil)
	mockRepo.On("CountSkillGroups", userID).Return(int64(2), nil)

	groupsAll, totalAll, err := svc.GetSkillGroups(userID, 0, 10, false)
	assert.NoError(t, err)
	assert.Len(t, groupsAll, 1)
	assert.Equal(t, int64(2), totalAll)
}

func TestSkillService_CreateSkill_Success(t *testing.T) {
	mockRepo := new(MockSkillRepository)
	svc := services.NewSkillService(mockRepo)

	userID := uuid.New()
	groupID := uuid.New()

	// Mock getting group succeeds
	mockRepo.On("GetSkillGroupByID", userID, groupID, false).Return(&models.SkillGroup{BaseModel: models.BaseModel{ID: groupID}}, nil)
	
	skill := &models.Skill{Name: "Go", UserID: userID, SkillGroupID: groupID}
	mockRepo.On("CreateSkill", skill).Return(skill, nil)

	created, err := svc.CreateSkill(skill)
	assert.NoError(t, err)
	assert.NotNil(t, created)
}

func TestSkillService_CreateSkill_GroupNotFound(t *testing.T) {
	mockRepo := new(MockSkillRepository)
	svc := services.NewSkillService(mockRepo)

	userID := uuid.New()
	groupID := uuid.New()

	// Mock getting group fails/returns nil
	mockRepo.On("GetSkillGroupByID", userID, groupID, false).Return((*models.SkillGroup)(nil), nil)
	
	skill := &models.Skill{Name: "Go", UserID: userID, SkillGroupID: groupID}
	
	created, err := svc.CreateSkill(skill)
	assert.Error(t, err)
	assert.Nil(t, created)
	assert.Equal(t, "skill group not found or does not belong to user", err.Error())
}

func TestSkillService_UpdateSkillGroup_NotFound(t *testing.T) {
	mockRepo := new(MockSkillRepository)
	svc := services.NewSkillService(mockRepo)

	userID := uuid.New()
	id := uuid.New()

	mockRepo.On("UpdateSkillGroup", userID, id, mock.Anything).Return((*models.SkillGroup)(nil), errors.New("not found"))

	updated, err := svc.UpdateSkillGroup(userID, id, map[string]interface{}{"name": "new"})
	assert.Error(t, err)
	assert.Nil(t, updated)
}
