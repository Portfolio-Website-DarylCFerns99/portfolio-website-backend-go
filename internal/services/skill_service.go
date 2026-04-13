package services

import (
	"errors"

	"github.com/google/uuid"

	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/repository"
)

type SkillService interface {
	// SkillGroups
	CreateSkillGroup(group *models.SkillGroup) (*models.SkillGroup, error)
	GetSkillGroups(userID uuid.UUID, skip, limit int, publicOnly bool) ([]models.SkillGroup, int64, error)
	GetSkillGroupByID(userID uuid.UUID, id uuid.UUID, publicOnly bool) (*models.SkillGroup, error)
	UpdateSkillGroup(userID uuid.UUID, id uuid.UUID, data map[string]interface{}) (*models.SkillGroup, error)
	UpdateSkillGroupVisibility(userID uuid.UUID, id uuid.UUID, isVisible bool) (*models.SkillGroup, error)
	DeleteSkillGroup(userID uuid.UUID, id uuid.UUID) error

	// Skills
	CreateSkill(skill *models.Skill) (*models.Skill, error)
	GetSkills(userID uuid.UUID, skip, limit int, publicOnly bool) ([]models.Skill, int64, error)
	GetSkillByID(userID uuid.UUID, id uuid.UUID) (*models.Skill, error)
	UpdateSkill(userID uuid.UUID, id uuid.UUID, data map[string]interface{}) (*models.Skill, error)
	UpdateSkillVisibility(userID uuid.UUID, id uuid.UUID, isVisible bool) (*models.Skill, error)
	DeleteSkill(userID uuid.UUID, id uuid.UUID) error
}

type skillService struct {
	repo repository.SkillRepository
}

func NewSkillService(repo repository.SkillRepository) SkillService {
	return &skillService{repo: repo}
}

// --- Skill Group Methods ---

func (s *skillService) CreateSkillGroup(group *models.SkillGroup) (*models.SkillGroup, error) {
	return s.repo.CreateSkillGroup(group)
}

func (s *skillService) GetSkillGroups(userID uuid.UUID, skip, limit int, publicOnly bool) ([]models.SkillGroup, int64, error) {
	var groups []models.SkillGroup
	var total int64
	var err error

	if publicOnly {
		groups, err = s.repo.GetVisibleSkillGroups(userID, skip, limit, true)
		if err == nil {
			total, err = s.repo.CountVisibleSkillGroups(userID)
		}
	} else {
		groups, err = s.repo.GetSkillGroups(userID, skip, limit, true)
		if err == nil {
			total, err = s.repo.CountSkillGroups(userID)
		}
	}

	if err != nil {
		return nil, 0, err
	}
	return groups, total, nil
}

func (s *skillService) GetSkillGroupByID(userID uuid.UUID, id uuid.UUID, publicOnly bool) (*models.SkillGroup, error) {
	if publicOnly {
		return s.repo.GetVisibleSkillGroupByID(userID, id, true)
	}
	return s.repo.GetSkillGroupByID(userID, id, true)
}

func (s *skillService) UpdateSkillGroup(userID uuid.UUID, id uuid.UUID, data map[string]interface{}) (*models.SkillGroup, error) {
	updated, err := s.repo.UpdateSkillGroup(userID, id, data)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, errors.New("skill group not found")
	}
	return updated, nil
}

func (s *skillService) UpdateSkillGroupVisibility(userID uuid.UUID, id uuid.UUID, isVisible bool) (*models.SkillGroup, error) {
	data := map[string]interface{}{"is_visible": isVisible}
	return s.UpdateSkillGroup(userID, id, data)
}

func (s *skillService) DeleteSkillGroup(userID uuid.UUID, id uuid.UUID) error {
	return s.repo.DeleteSkillGroup(userID, id)
}

// --- Skill Methods ---

func (s *skillService) CreateSkill(skill *models.Skill) (*models.Skill, error) {
	// First check if user owns the skill group
	group, err := s.repo.GetSkillGroupByID(skill.UserID, skill.SkillGroupID, false)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, errors.New("skill group not found or does not belong to user")
	}
	return s.repo.CreateSkill(skill)
}

func (s *skillService) GetSkills(userID uuid.UUID, skip, limit int, publicOnly bool) ([]models.Skill, int64, error) {
	var skills []models.Skill
	var total int64
	var err error

	if publicOnly {
		skills, err = s.repo.GetVisibleSkills(userID, skip, limit)
		if err == nil {
			total, err = s.repo.CountVisibleSkills(userID)
		}
	} else {
		skills, err = s.repo.GetSkills(userID, skip, limit)
		if err == nil {
			total, err = s.repo.CountSkills(userID)
		}
	}

	if err != nil {
		return nil, 0, err
	}
	return skills, total, nil
}

func (s *skillService) GetSkillByID(userID uuid.UUID, id uuid.UUID) (*models.Skill, error) {
	return s.repo.GetSkillByID(userID, id)
}

func (s *skillService) UpdateSkill(userID uuid.UUID, id uuid.UUID, data map[string]interface{}) (*models.Skill, error) {
	updated, err := s.repo.UpdateSkill(userID, id, data)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, errors.New("skill not found")
	}
	return updated, nil
}

func (s *skillService) UpdateSkillVisibility(userID uuid.UUID, id uuid.UUID, isVisible bool) (*models.Skill, error) {
	data := map[string]interface{}{"is_visible": isVisible}
	return s.UpdateSkill(userID, id, data)
}

func (s *skillService) DeleteSkill(userID uuid.UUID, id uuid.UUID) error {
	return s.repo.DeleteSkill(userID, id)
}
