package experiences

import (
	"portfolio-website-backend/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockExperienceService struct {
	mock.Mock
}

func (m *MockExperienceService) CreateExperience(experience *models.Experience) (*models.Experience, error) {
	args := m.Called(experience)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Experience), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockExperienceService) GetExperiences(userID uuid.UUID, skip, limit int, onlyVisible bool) ([]models.Experience, int64, error) {
	args := m.Called(userID, skip, limit, onlyVisible)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Experience), args.Get(1).(int64), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func (m *MockExperienceService) GetExperiencesByType(userID uuid.UUID, expType string, skip, limit int, onlyVisible bool) ([]models.Experience, int64, error) {
	args := m.Called(userID, expType, skip, limit, onlyVisible)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Experience), args.Get(1).(int64), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func (m *MockExperienceService) GetExperienceByID(userID uuid.UUID, id uuid.UUID, onlyVisible bool) (*models.Experience, error) {
	args := m.Called(userID, id, onlyVisible)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Experience), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockExperienceService) UpdateExperience(userID uuid.UUID, id uuid.UUID, updateData map[string]interface{}) (*models.Experience, error) {
	args := m.Called(userID, id, updateData)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Experience), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockExperienceService) UpdateExperienceVisibility(userID uuid.UUID, id uuid.UUID, isVisible bool) (*models.Experience, error) {
	args := m.Called(userID, id, isVisible)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Experience), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockExperienceService) DeleteExperience(userID uuid.UUID, id uuid.UUID) error {
	args := m.Called(userID, id)
	return args.Error(0)
}
