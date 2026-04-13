package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"portfolio-website-backend/internal/models"
)

type SkillRepository interface {
	// SkillGroups
	CreateSkillGroup(group *models.SkillGroup) (*models.SkillGroup, error)
	GetSkillGroups(userID uuid.UUID, skip, limit int, eagerLoad bool) ([]models.SkillGroup, error)
	GetVisibleSkillGroups(userID uuid.UUID, skip, limit int, eagerLoad bool) ([]models.SkillGroup, error)
	GetSkillGroupByID(userID uuid.UUID, id uuid.UUID, eagerLoad bool) (*models.SkillGroup, error)
	GetVisibleSkillGroupByID(userID uuid.UUID, id uuid.UUID, eagerLoad bool) (*models.SkillGroup, error)
	UpdateSkillGroup(userID uuid.UUID, id uuid.UUID, data map[string]interface{}) (*models.SkillGroup, error)
	DeleteSkillGroup(userID uuid.UUID, id uuid.UUID) error
	CountSkillGroups(userID uuid.UUID) (int64, error)
	CountVisibleSkillGroups(userID uuid.UUID) (int64, error)

	// Skills
	CreateSkill(skill *models.Skill) (*models.Skill, error)
	GetSkills(userID uuid.UUID, skip, limit int) ([]models.Skill, error)
	GetVisibleSkills(userID uuid.UUID, skip, limit int) ([]models.Skill, error)
	GetSkillByID(userID uuid.UUID, id uuid.UUID) (*models.Skill, error)
	UpdateSkill(userID uuid.UUID, id uuid.UUID, data map[string]interface{}) (*models.Skill, error)
	DeleteSkill(userID uuid.UUID, id uuid.UUID) error
	CountSkills(userID uuid.UUID) (int64, error)
	CountVisibleSkills(userID uuid.UUID) (int64, error)
}

type skillRepository struct {
	db *gorm.DB
}

func NewSkillRepository(db *gorm.DB) SkillRepository {
	return &skillRepository{db: db}
}

// --- Skill Group Methods ---

func (r *skillRepository) CreateSkillGroup(group *models.SkillGroup) (*models.SkillGroup, error) {
	if err := r.db.Create(group).Error; err != nil {
		return nil, err
	}
	return group, nil
}

func (r *skillRepository) GetSkillGroups(userID uuid.UUID, skip, limit int, eagerLoad bool) ([]models.SkillGroup, error) {
	var groups []models.SkillGroup
	query := r.db.Where("user_id = ?", userID)
	if eagerLoad {
		query = query.Preload("Skills")
	}
	if err := query.Offset(skip).Limit(limit).Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *skillRepository) GetVisibleSkillGroups(userID uuid.UUID, skip, limit int, eagerLoad bool) ([]models.SkillGroup, error) {
	var groups []models.SkillGroup
	query := r.db.Where("user_id = ? AND is_visible = ?", userID, true)
	if eagerLoad {
		// Only preload visible skills if getting visible groups? Or usually we just preload visible skills.
		query = query.Preload("Skills", "is_visible = ?", true)
	}
	if err := query.Offset(skip).Limit(limit).Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *skillRepository) GetSkillGroupByID(userID uuid.UUID, id uuid.UUID, eagerLoad bool) (*models.SkillGroup, error) {
	var group models.SkillGroup
	query := r.db.Where("user_id = ?", userID)
	if eagerLoad {
		query = query.Preload("Skills")
	}
	if err := query.First(&group, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &group, nil
}

func (r *skillRepository) GetVisibleSkillGroupByID(userID uuid.UUID, id uuid.UUID, eagerLoad bool) (*models.SkillGroup, error) {
	var group models.SkillGroup
	query := r.db.Where("user_id = ? AND is_visible = ?", userID, true)
	if eagerLoad {
		query = query.Preload("Skills", "is_visible = ?", true)
	}
	if err := query.First(&group, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &group, nil
}

func (r *skillRepository) UpdateSkillGroup(userID uuid.UUID, id uuid.UUID, data map[string]interface{}) (*models.SkillGroup, error) {
	var group models.SkillGroup
	if err := r.db.Where("user_id = ?", userID).First(&group, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if err := r.db.Model(&group).Updates(data).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *skillRepository) DeleteSkillGroup(userID uuid.UUID, id uuid.UUID) error {
	result := r.db.Where("user_id = ?", userID).Delete(&models.SkillGroup{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}
	return nil
}

func (r *skillRepository) CountSkillGroups(userID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.Model(&models.SkillGroup{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *skillRepository) CountVisibleSkillGroups(userID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.Model(&models.SkillGroup{}).Where("user_id = ? AND is_visible = ?", userID, true).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// --- Skill Methods ---

func (r *skillRepository) CreateSkill(skill *models.Skill) (*models.Skill, error) {
	if err := r.db.Create(skill).Error; err != nil {
		return nil, err
	}
	return skill, nil
}

func (r *skillRepository) GetSkills(userID uuid.UUID, skip, limit int) ([]models.Skill, error) {
	var skills []models.Skill
	if err := r.db.Where("user_id = ?", userID).Offset(skip).Limit(limit).Find(&skills).Error; err != nil {
		return nil, err
	}
	return skills, nil
}

func (r *skillRepository) GetVisibleSkills(userID uuid.UUID, skip, limit int) ([]models.Skill, error) {
	var skills []models.Skill
	if err := r.db.Where("user_id = ? AND is_visible = ?", userID, true).Offset(skip).Limit(limit).Find(&skills).Error; err != nil {
		return nil, err
	}
	return skills, nil
}

func (r *skillRepository) CountSkills(userID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Skill{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *skillRepository) CountVisibleSkills(userID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Skill{}).Where("user_id = ? AND is_visible = ?", userID, true).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *skillRepository) GetSkillByID(userID uuid.UUID, id uuid.UUID) (*models.Skill, error) {
	var skill models.Skill
	if err := r.db.Where("user_id = ?", userID).First(&skill, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &skill, nil
}

func (r *skillRepository) UpdateSkill(userID uuid.UUID, id uuid.UUID, data map[string]interface{}) (*models.Skill, error) {
	var skill models.Skill
	if err := r.db.Where("user_id = ?", userID).First(&skill, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if err := r.db.Model(&skill).Updates(data).Error; err != nil {
		return nil, err
	}
	return &skill, nil
}

func (r *skillRepository) DeleteSkill(userID uuid.UUID, id uuid.UUID) error {
	result := r.db.Where("user_id = ?", userID).Delete(&models.Skill{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}
	return nil
}
