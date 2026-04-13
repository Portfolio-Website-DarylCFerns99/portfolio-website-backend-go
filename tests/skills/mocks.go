package skills

import (
	"portfolio-website-backend/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockSkillService struct {
	mock.Mock
}

// SkillGroups
func (m *MockSkillService) CreateSkillGroup(group *models.SkillGroup) (*models.SkillGroup, error) {
	args := m.Called(group)
	if args.Get(0) != nil {
		return args.Get(0).(*models.SkillGroup), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSkillService) GetSkillGroups(userID uuid.UUID, skip, limit int, publicOnly bool) ([]models.SkillGroup, int64, error) {
	args := m.Called(userID, skip, limit, publicOnly)
	var cat []models.SkillGroup
	if args.Get(0) != nil {
		cat = args.Get(0).([]models.SkillGroup)
	}
	var total int64
	if args.Get(1) != nil {
		if t, ok := args.Get(1).(int64); ok {
			total = t
		} else if t, ok := args.Get(1).(int); ok {
			total = int64(t)
		}
	}
	return cat, total, args.Error(2)
}

func (m *MockSkillService) GetSkillGroupByID(userID uuid.UUID, id uuid.UUID, publicOnly bool) (*models.SkillGroup, error) {
	args := m.Called(userID, id, publicOnly)
	if args.Get(0) != nil {
		return args.Get(0).(*models.SkillGroup), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSkillService) UpdateSkillGroup(userID uuid.UUID, id uuid.UUID, data map[string]interface{}) (*models.SkillGroup, error) {
	args := m.Called(userID, id, data)
	if args.Get(0) != nil {
		return args.Get(0).(*models.SkillGroup), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSkillService) UpdateSkillGroupVisibility(userID uuid.UUID, id uuid.UUID, isVisible bool) (*models.SkillGroup, error) {
	args := m.Called(userID, id, isVisible)
	if args.Get(0) != nil {
		return args.Get(0).(*models.SkillGroup), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSkillService) DeleteSkillGroup(userID uuid.UUID, id uuid.UUID) error {
	args := m.Called(userID, id)
	return args.Error(0)
}

// Skills
func (m *MockSkillService) CreateSkill(skill *models.Skill) (*models.Skill, error) {
	args := m.Called(skill)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Skill), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSkillService) GetSkills(userID uuid.UUID, skip, limit int, publicOnly bool) ([]models.Skill, int64, error) {
	args := m.Called(userID, skip, limit, publicOnly)
	var cat []models.Skill
	if args.Get(0) != nil {
		cat = args.Get(0).([]models.Skill)
	}
	var total int64
	if args.Get(1) != nil {
		if t, ok := args.Get(1).(int64); ok {
			total = t
		} else if t, ok := args.Get(1).(int); ok {
			total = int64(t)
		}
	}
	return cat, total, args.Error(2)
}

func (m *MockSkillService) GetSkillByID(userID uuid.UUID, id uuid.UUID) (*models.Skill, error) {
	args := m.Called(userID, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Skill), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSkillService) UpdateSkill(userID uuid.UUID, id uuid.UUID, data map[string]interface{}) (*models.Skill, error) {
	args := m.Called(userID, id, data)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Skill), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSkillService) UpdateSkillVisibility(userID uuid.UUID, id uuid.UUID, isVisible bool) (*models.Skill, error) {
	args := m.Called(userID, id, isVisible)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Skill), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSkillService) DeleteSkill(userID uuid.UUID, id uuid.UUID) error {
	args := m.Called(userID, id)
	return args.Error(0)
}

type MockSkillRepository struct {
	mock.Mock
}

func (m *MockSkillRepository) CreateSkillGroup(group *models.SkillGroup) (*models.SkillGroup, error) {
	args := m.Called(group)
	if args.Get(0) != nil {
		return args.Get(0).(*models.SkillGroup), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSkillRepository) GetSkillGroups(userID uuid.UUID, skip, limit int, loadSkills bool) ([]models.SkillGroup, error) {
	args := m.Called(userID, skip, limit, loadSkills)
	if args.Get(0) != nil {
		return args.Get(0).([]models.SkillGroup), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSkillRepository) GetVisibleSkillGroups(userID uuid.UUID, skip, limit int, loadSkills bool) ([]models.SkillGroup, error) {
	args := m.Called(userID, skip, limit, loadSkills)
	if args.Get(0) != nil {
		return args.Get(0).([]models.SkillGroup), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSkillRepository) CountSkillGroups(userID uuid.UUID) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSkillRepository) CountVisibleSkillGroups(userID uuid.UUID) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSkillRepository) GetSkillGroupByID(userID uuid.UUID, id uuid.UUID, loadSkills bool) (*models.SkillGroup, error) {
	args := m.Called(userID, id, loadSkills)
	if args.Get(0) != nil {
		return args.Get(0).(*models.SkillGroup), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSkillRepository) GetVisibleSkillGroupByID(userID uuid.UUID, id uuid.UUID, loadSkills bool) (*models.SkillGroup, error) {
	args := m.Called(userID, id, loadSkills)
	if args.Get(0) != nil {
		return args.Get(0).(*models.SkillGroup), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSkillRepository) UpdateSkillGroup(userID uuid.UUID, id uuid.UUID, data map[string]interface{}) (*models.SkillGroup, error) {
	args := m.Called(userID, id, data)
	if args.Get(0) != nil {
		return args.Get(0).(*models.SkillGroup), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSkillRepository) DeleteSkillGroup(userID uuid.UUID, id uuid.UUID) error {
	args := m.Called(userID, id)
	return args.Error(0)
}

func (m *MockSkillRepository) CreateSkill(skill *models.Skill) (*models.Skill, error) {
	args := m.Called(skill)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Skill), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSkillRepository) GetSkills(userID uuid.UUID, skip, limit int) ([]models.Skill, error) {
	args := m.Called(userID, skip, limit)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Skill), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSkillRepository) GetVisibleSkills(userID uuid.UUID, skip, limit int) ([]models.Skill, error) {
	args := m.Called(userID, skip, limit)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Skill), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSkillRepository) CountSkills(userID uuid.UUID) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSkillRepository) CountVisibleSkills(userID uuid.UUID) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSkillRepository) GetSkillByID(userID uuid.UUID, id uuid.UUID) (*models.Skill, error) {
	args := m.Called(userID, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Skill), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSkillRepository) UpdateSkill(userID uuid.UUID, id uuid.UUID, data map[string]interface{}) (*models.Skill, error) {
	args := m.Called(userID, id, data)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Skill), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSkillRepository) DeleteSkill(userID uuid.UUID, id uuid.UUID) error {
	args := m.Called(userID, id)
	return args.Error(0)
}
