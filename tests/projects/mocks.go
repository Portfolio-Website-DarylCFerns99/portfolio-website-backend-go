package projects

import (
	"portfolio-website-backend/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockProjectService struct {
	mock.Mock
}

func (m *MockProjectService) CreateProject(project *models.Project) (*models.Project, error) {
	args := m.Called(project)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Project), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProjectService) GetProjects(userID uuid.UUID, skip, limit int, onlyVisible bool) ([]models.Project, int64, error) {
	args := m.Called(userID, skip, limit, onlyVisible)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Project), args.Get(1).(int64), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func (m *MockProjectService) GetProjectByID(userID uuid.UUID, id uuid.UUID, onlyVisible bool) (*models.Project, error) {
	args := m.Called(userID, id, onlyVisible)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Project), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProjectService) UpdateProject(userID uuid.UUID, id uuid.UUID, updateData map[string]interface{}) (*models.Project, error) {
	args := m.Called(userID, id, updateData)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Project), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProjectService) UpdateProjectVisibility(userID uuid.UUID, id uuid.UUID, isVisible bool) (*models.Project, error) {
	args := m.Called(userID, id, isVisible)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Project), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProjectService) DeleteProject(userID uuid.UUID, id uuid.UUID) error {
	args := m.Called(userID, id)
	return args.Error(0)
}
