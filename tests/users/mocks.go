package users

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"portfolio-website-backend/internal/dto"
	"portfolio-website-backend/internal/models"
)

// MockUserRepository implements repository.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) (*models.User, error) {
	args := m.Called(user)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetByEmailorUsername(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) Update(id uuid.UUID, data map[string]interface{}) (*models.User, error) {
	args := m.Called(id, data)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) List() ([]models.User, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

// MockUserService implements services.UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) ValidateFeaturedSkills(userID uuid.UUID, skillIDs []string) []string {
	args := m.Called(userID, skillIDs)
	return args.Get(0).([]string)
}

func (m *MockUserService) UpdateUserProfile(id uuid.UUID, updateData map[string]interface{}) (*models.User, error) {
	args := m.Called(id, updateData)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserService) GetPublicPortfolioData(userID uuid.UUID) (*dto.PublicDataResponse, error) {
	args := m.Called(userID)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.PublicDataResponse), args.Error(1)
	}
	return nil, args.Error(1)
}
