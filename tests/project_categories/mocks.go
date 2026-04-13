package project_categories

import (
	"portfolio-website-backend/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockProjectCategoryService struct {
	mock.Mock
}

func (m *MockProjectCategoryService) CreateCategory(category *models.ProjectCategory) (*models.ProjectCategory, error) {
	args := m.Called(category)
	if args.Get(0) != nil {
		return args.Get(0).(*models.ProjectCategory), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProjectCategoryService) GetCategories(userID uuid.UUID, skip, limit int, onlyVisible bool) ([]models.ProjectCategory, int64, error) {
	args := m.Called(userID, skip, limit, onlyVisible)
	var cat []models.ProjectCategory
	if args.Get(0) != nil {
		cat = args.Get(0).([]models.ProjectCategory)
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

func (m *MockProjectCategoryService) GetCategoryByID(userID uuid.UUID, id uuid.UUID, onlyVisible bool) (*models.ProjectCategory, error) {
	args := m.Called(userID, id, onlyVisible)
	if args.Get(0) != nil {
		return args.Get(0).(*models.ProjectCategory), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProjectCategoryService) UpdateCategory(userID uuid.UUID, id uuid.UUID, updateData map[string]interface{}) (*models.ProjectCategory, error) {
	args := m.Called(userID, id, updateData)
	if args.Get(0) != nil {
		return args.Get(0).(*models.ProjectCategory), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProjectCategoryService) DeleteCategory(userID uuid.UUID, id uuid.UUID) error {
	args := m.Called(userID, id)
	return args.Error(0)
}
