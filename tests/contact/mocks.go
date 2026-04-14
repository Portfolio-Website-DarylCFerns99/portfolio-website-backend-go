package contact

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"portfolio-website-backend/internal/dto"
	"portfolio-website-backend/internal/models"
)

// ---------------------------------------------------------------------------
// MockContactService — used by handler tests
// ---------------------------------------------------------------------------

type MockContactService struct {
	mock.Mock
}

func (m *MockContactService) SendContactEmail(userID uuid.UUID, req dto.ContactRequest) error {
	args := m.Called(userID, req)
	return args.Error(0)
}

// ---------------------------------------------------------------------------
// MockUserRepository — used by service tests
// ---------------------------------------------------------------------------

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
