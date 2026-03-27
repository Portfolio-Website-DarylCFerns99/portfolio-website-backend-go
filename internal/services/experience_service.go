package services

import (
	"log"

	"github.com/google/uuid"

	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/repository"
)

type ExperienceService interface {
	CreateExperience(experience *models.Experience) (*models.Experience, error)
	GetExperiences(userID uuid.UUID, skip, limit int, onlyVisible bool) ([]models.Experience, int64, error)
	GetExperiencesByType(userID uuid.UUID, expType string, skip, limit int, onlyVisible bool) ([]models.Experience, int64, error)
	GetExperienceByID(userID uuid.UUID, id uuid.UUID, onlyVisible bool) (*models.Experience, error)
	UpdateExperience(userID uuid.UUID, id uuid.UUID, updateData map[string]interface{}) (*models.Experience, error)
	UpdateExperienceVisibility(userID uuid.UUID, id uuid.UUID, isVisible bool) (*models.Experience, error)
	DeleteExperience(userID uuid.UUID, id uuid.UUID) error
}

type experienceService struct {
	repo repository.ExperienceRepository
}

func NewExperienceService(repo repository.ExperienceRepository) ExperienceService {
	return &experienceService{repo: repo}
}

func (s *experienceService) CreateExperience(experience *models.Experience) (*models.Experience, error) {
	log.Printf("Creating %s entry: %s", experience.Type, experience.Title)
	return s.repo.Create(experience)
}

func (s *experienceService) GetExperiences(userID uuid.UUID, skip, limit int, onlyVisible bool) ([]models.Experience, int64, error) {
	log.Printf("Retrieving all experiences for user %s (skip=%d, limit=%d, only_visible=%v)", userID, skip, limit, onlyVisible)

	var experiences []models.Experience
	var total int64
	var err error

	if onlyVisible {
		experiences, err = s.repo.GetVisible(userID, skip, limit)
		total, _ = s.repo.CountVisible(userID)
	} else {
		experiences, err = s.repo.GetAll(userID, skip, limit)
		total, _ = s.repo.Count(userID)
	}

	return experiences, total, err
}

func (s *experienceService) GetExperiencesByType(userID uuid.UUID, expType string, skip, limit int, onlyVisible bool) ([]models.Experience, int64, error) {
	log.Printf("Retrieving %s entries for user %s (skip=%d, limit=%d, only_visible=%v)", expType, userID, skip, limit, onlyVisible)

	experiences, err := s.repo.GetByType(userID, expType, skip, limit, onlyVisible)
	total, _ := s.repo.CountByType(userID, expType, onlyVisible)

	return experiences, total, err
}

func (s *experienceService) GetExperienceByID(userID uuid.UUID, id uuid.UUID, onlyVisible bool) (*models.Experience, error) {
	log.Printf("Retrieving experience with ID: %s for user %s, only_visible=%v", id, userID, onlyVisible)

	if onlyVisible {
		return s.repo.GetVisibleByID(userID, id)
	}
	return s.repo.GetByID(userID, id)
}

func (s *experienceService) UpdateExperience(userID uuid.UUID, id uuid.UUID, updateData map[string]interface{}) (*models.Experience, error) {
	log.Printf("Updating experience with ID: %s for user %s", id, userID)

	experience, err := s.repo.GetByID(userID, id)
	if err != nil || experience == nil {
		log.Printf("Experience with ID %s not found for user %s", id, userID)
		return nil, err
	}

	return s.repo.Update(userID, id, updateData)
}

func (s *experienceService) UpdateExperienceVisibility(userID uuid.UUID, id uuid.UUID, isVisible bool) (*models.Experience, error) {
	log.Printf("Updating visibility for experience ID: %s to %v for user %s", id, isVisible, userID)

	experience, err := s.repo.GetByID(userID, id)
	if err != nil || experience == nil {
		log.Printf("Experience with ID %s not found for user %s", id, userID)
		return nil, err
	}

	updateData := map[string]interface{}{"is_visible": isVisible}
	return s.repo.Update(userID, id, updateData)
}

func (s *experienceService) DeleteExperience(userID uuid.UUID, id uuid.UUID) error {
	log.Printf("Deleting experience with ID: %s for user %s", id, userID)
	return s.repo.Delete(userID, id)
}
