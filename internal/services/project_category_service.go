package services

import (
	"log"

	"github.com/google/uuid"

	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/repository"
)

type ProjectCategoryService interface {
	CreateCategory(category *models.ProjectCategory) (*models.ProjectCategory, error)
	GetCategories(userID uuid.UUID, skip, limit int, onlyVisible bool) ([]models.ProjectCategory, int64, error)
	GetCategoryByID(userID uuid.UUID, id uuid.UUID, onlyVisible bool) (*models.ProjectCategory, error)
	UpdateCategory(userID uuid.UUID, id uuid.UUID, updateData map[string]interface{}) (*models.ProjectCategory, error)
	DeleteCategory(userID uuid.UUID, id uuid.UUID) error
}

type projectCategoryService struct {
	repo repository.ProjectCategoryRepository
}

func NewProjectCategoryService(repo repository.ProjectCategoryRepository) ProjectCategoryService {
	return &projectCategoryService{repo: repo}
}

func (s *projectCategoryService) CreateCategory(category *models.ProjectCategory) (*models.ProjectCategory, error) {
	log.Printf("Creating project category: %s", category.Name)
	return s.repo.Create(category)
}

func (s *projectCategoryService) GetCategories(userID uuid.UUID, skip, limit int, onlyVisible bool) ([]models.ProjectCategory, int64, error) {
	log.Printf("Retrieving project categories for user %s (skip=%d, limit=%d, only_visible=%v)", userID, skip, limit, onlyVisible)

	if onlyVisible {
		return s.repo.GetVisible(userID, skip, limit)
	}
	return s.repo.GetAll(userID, skip, limit)
}

func (s *projectCategoryService) GetCategoryByID(userID uuid.UUID, id uuid.UUID, onlyVisible bool) (*models.ProjectCategory, error) {
	log.Printf("Retrieving project category with ID: %s for user %s, only_visible=%v", id, userID, onlyVisible)

	if onlyVisible {
		return s.repo.GetVisibleByID(userID, id)
	}
	return s.repo.GetByID(userID, id)
}

func (s *projectCategoryService) UpdateCategory(userID uuid.UUID, id uuid.UUID, updateData map[string]interface{}) (*models.ProjectCategory, error) {
	log.Printf("Updating project category with ID: %s for user %s", id, userID)

	category, err := s.repo.GetByID(userID, id)
	if err != nil || category == nil {
		log.Printf("Project category with ID %s not found for user %s", id, userID)
		return nil, err
	}

	return s.repo.Update(userID, id, updateData)
}

func (s *projectCategoryService) DeleteCategory(userID uuid.UUID, id uuid.UUID) error {
	log.Printf("Deleting project category with ID: %s for user %s", id, userID)
	return s.repo.Delete(userID, id)
}
