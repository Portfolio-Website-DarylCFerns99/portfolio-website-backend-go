package services

import (
	"log"
	"time"

	"github.com/google/uuid"

	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/repository"
	"portfolio-website-backend/internal/utils"
)

type ProjectService interface {
	CreateProject(project *models.Project) (*models.Project, error)
	GetProjects(userID uuid.UUID, skip, limit int, onlyVisible bool) ([]models.Project, int64, error)
	GetProjectByID(userID uuid.UUID, id uuid.UUID, onlyVisible bool) (*models.Project, error)
	UpdateProject(userID uuid.UUID, id uuid.UUID, updateData map[string]interface{}) (*models.Project, error)
	UpdateProjectVisibility(userID uuid.UUID, id uuid.UUID, isVisible bool) (*models.Project, error)
	DeleteProject(userID uuid.UUID, id uuid.UUID) error
}

type projectService struct {
	repo repository.ProjectRepository
}

func NewProjectService(repo repository.ProjectRepository) ProjectService {
	return &projectService{repo: repo}
}

func (s *projectService) refreshGithubData(project *models.Project) {
	if project.Type == "github" && project.URL != nil && *project.URL != "" {
		_, githubFullData, err := utils.FetchGithubData(*project.URL)
		if err == nil {
			project.AdditionalData = githubFullData
			newExpiry := time.Now().UTC().Add(24 * time.Hour)
			project.ExpiryDate = &newExpiry
		} else {
			log.Printf("Error refreshing GitHub data for project %s: %v", project.ID, err)
		}
	}
}

func (s *projectService) hydrateInitialGithubData(project *models.Project) {
	if project.Type == "github" && project.URL != nil && *project.URL != "" {
		basicData, githubFullData, err := utils.FetchGithubData(*project.URL)
		if err == nil {
			if project.Title == "" {
				if title, ok := basicData["title"].(string); ok {
					project.Title = title
				}
			}
			if project.Description == nil || *project.Description == "" {
				if desc, ok := basicData["description"].(string); ok {
					project.Description = &desc
				}
			}
			newExpiry := time.Now().UTC().Add(24 * time.Hour)
			project.ExpiryDate = &newExpiry
			project.AdditionalData = githubFullData
		} else {
			log.Printf("Error fetching initial GitHub data for project: %v", err)
		}
	}
}

func (s *projectService) CreateProject(project *models.Project) (*models.Project, error) {
	log.Printf("Creating %s project: %s", project.Type, project.Title)

	s.hydrateInitialGithubData(project)

	return s.repo.Create(project)
}

func (s *projectService) checkAndRefreshProject(userID uuid.UUID, project *models.Project) {
	if project.Type == "github" && project.ExpiryDate != nil && project.ExpiryDate.Before(time.Now().UTC()) {
		log.Printf("Refreshing GitHub data for project %s", project.ID)
		s.refreshGithubData(project)
		// Update DB with refreshed data
		updateData := map[string]interface{}{
			"additional_data": project.AdditionalData,
			"expiry_date":     project.ExpiryDate,
		}
		_, _ = s.repo.Update(userID, project.ID, updateData)
	}
}

func (s *projectService) GetProjects(userID uuid.UUID, skip, limit int, onlyVisible bool) ([]models.Project, int64, error) {
	log.Printf("Retrieving all projects for user %s (skip=%d, limit=%d, only_visible=%v)", userID, skip, limit, onlyVisible)

	var projects []models.Project
	var total int64
	var err error

	if onlyVisible {
		projects, err = s.repo.GetVisible(userID, skip, limit)
		total, _ = s.repo.CountVisible(userID)
	} else {
		projects, err = s.repo.GetAll(userID, skip, limit)
		total, _ = s.repo.Count(userID)
	}

	if err == nil {
		for i := range projects {
			s.checkAndRefreshProject(userID, &projects[i])
		}
	}

	return projects, total, err
}

func (s *projectService) GetProjectByID(userID uuid.UUID, id uuid.UUID, onlyVisible bool) (*models.Project, error) {
	log.Printf("Retrieving project with ID: %s for user %s, only_visible=%v", id, userID, onlyVisible)

	var project *models.Project
	var err error

	if onlyVisible {
		project, err = s.repo.GetVisibleByID(userID, id)
	} else {
		project, err = s.repo.GetByID(userID, id)
	}

	if err == nil && project != nil {
		s.checkAndRefreshProject(userID, project)
	}

	return project, err
}

func (s *projectService) UpdateProject(userID uuid.UUID, id uuid.UUID, updateData map[string]interface{}) (*models.Project, error) {
	log.Printf("Updating project with ID: %s for user %s", id, userID)

	project, err := s.repo.GetByID(userID, id)
	if err != nil || project == nil {
		log.Printf("Project with ID %s not found for user %s", id, userID)
		return nil, err
	}

	if pType, ok := updateData["type"].(string); ok && pType == "github" {
		url := project.URL
		if updateUrlStr, uOk := updateData["url"].(string); uOk {
			url = &updateUrlStr
		}
		if url != nil && *url != "" {
			_, githubFullData, fetchErr := utils.FetchGithubData(*url)
			if fetchErr == nil {
				updateData["additional_data"] = githubFullData
				newExpiry := time.Now().UTC().Add(24 * time.Hour)
				updateData["expiry_date"] = &newExpiry
			} else {
				log.Printf("Error fetching GitHub data during update for project %s: %v", id, fetchErr)
			}
		}
	} else if pType == "custom" {
		updateData["expiry_date"] = nil
	}

	return s.repo.Update(userID, id, updateData)
}

func (s *projectService) UpdateProjectVisibility(userID uuid.UUID, id uuid.UUID, isVisible bool) (*models.Project, error) {
	log.Printf("Updating visibility for project ID: %s to %v for user %s", id, isVisible, userID)

	project, err := s.repo.GetByID(userID, id)
	if err != nil || project == nil {
		log.Printf("Project with ID %s not found for user %s", id, userID)
		return nil, err
	}

	updateData := map[string]interface{}{"is_visible": isVisible}
	return s.repo.Update(userID, id, updateData)
}

func (s *projectService) DeleteProject(userID uuid.UUID, id uuid.UUID) error {
	log.Printf("Deleting project with ID: %s for user %s", id, userID)
	return s.repo.Delete(userID, id)
}
