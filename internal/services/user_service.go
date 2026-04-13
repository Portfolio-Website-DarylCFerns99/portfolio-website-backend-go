package services

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"portfolio-website-backend/internal/dto"
	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/repository"
)

type UserService interface {
	ValidateFeaturedSkills(userID uuid.UUID, skillIDs []string) []string
	UpdateUserProfile(id uuid.UUID, updateData map[string]interface{}) (*models.User, error)
	GetPublicPortfolioData(userID uuid.UUID) (*dto.PublicDataResponse, error)
}

type userService struct {
	db       *gorm.DB
	userRepo repository.UserRepository
}

// NewUserService takes gorm.DB directly to allow raw queries if needed across multiple repos
func NewUserService(db *gorm.DB, userRepo repository.UserRepository) UserService {
	return &userService{
		db:       db,
		userRepo: userRepo,
	}
}

func (s *userService) ValidateFeaturedSkills(userID uuid.UUID, skillIDs []string) []string {
	if len(skillIDs) == 0 {
		return []string{}
	}

	log.Printf("Validating %d featured skill IDs for user %s", len(skillIDs), userID)

	// Convert string IDs to UUIDs for querying
	var validSkillIDs []string

	// Query the database to find which skill IDs exist
	// In Python this was: `self.db.query(Skill.id).filter(Skill.id.in_(skill_ids)).all()`
	var existingSkills []models.Skill
	if err := s.db.Where("id IN ?", skillIDs).Select("id").Find(&existingSkills).Error; err != nil {
		log.Printf("Error validating skills: %v", err)
		return validSkillIDs
	}

	for _, skill := range existingSkills {
		validSkillIDs = append(validSkillIDs, skill.ID.String())
	}

	discardedCount := len(skillIDs) - len(validSkillIDs)
	if discardedCount > 0 {
		log.Printf("Discarded %d invalid skill IDs", discardedCount)
	}

	log.Printf("Validated %d skill IDs for user %s", len(validSkillIDs), userID)
	return validSkillIDs
}

func (s *userService) UpdateUserProfile(id uuid.UUID, updateData map[string]interface{}) (*models.User, error) {
	// Validate featured skill IDs if they're being updated
	if featured, ok := updateData["featured_skill_ids"]; ok {
		// Type assertion for the slice
		if skillIDs, isSlice := featured.([]string); isSlice {
			validIDs := s.ValidateFeaturedSkills(id, skillIDs)
			updateData["featured_skill_ids"] = models.JSONStringArray(validIDs)
		} else if skillInterfaceSlice, isInterfaceSlice := featured.([]interface{}); isInterfaceSlice {
			var skillIDs []string
			for _, v := range skillInterfaceSlice {
				if str, isStr := v.(string); isStr {
					skillIDs = append(skillIDs, str)
				}
			}
			validIDs := s.ValidateFeaturedSkills(id, skillIDs)
			updateData["featured_skill_ids"] = models.JSONStringArray(validIDs)
		}
	}

	// Convert social_links to JSONMapArray so GORM serializes it correctly instead of storing it as a generic object or raw bytes
	if sl, ok := updateData["social_links"]; ok {
		// Attempt to read it as a slice of maps
		if slSlice, isSlice := sl.([]interface{}); isSlice {
			var jsonMapArray models.JSONMapArray
			for _, item := range slSlice {
				if m, isMap := item.(map[string]interface{}); isMap {
					jsonMapArray = append(jsonMapArray, m)
				}
			}
			updateData["social_links"] = jsonMapArray
		} else if slMap, isMap := sl.(map[string]interface{}); isMap {
			// If frontend accidentally sends a single object, wrap it
			updateData["social_links"] = models.JSONMapArray{slMap}
		} else if slBytes, isBytes := sl.([]byte); isBytes {
			var parsed models.JSONMapArray
			if err := parsed.Scan(slBytes); err == nil {
				updateData["social_links"] = parsed
			}
		}
	}

	// Delegate to repository
	return s.userRepo.Update(id, updateData)
}

func calculateTotalExperience(experiences []models.Experience) interface{} {
	var totalYears float64
	now := time.Now()

	for _, exp := range experiences {
		if exp.Type != "experience" {
			continue
		}
		end := now
		if exp.EndDate != nil {
			end = exp.EndDate.Time
		}

		days := end.Sub(exp.StartDate.Time).Hours() / 24.0
		years := days / 365.25
		totalYears += years
	}

	floorYears := int(math.Floor(totalYears))
	if floorYears < 2 {
		return floorYears
	}
	if totalYears > float64(floorYears) {
		return fmt.Sprintf("%d+", floorYears)
	}
	return fmt.Sprintf("%d", floorYears)
}

func (s *userService) GetPublicPortfolioData(userID uuid.UUID) (*dto.PublicDataResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found")
	}

	var skillGroups []models.SkillGroup
	s.db.Preload("Skills", "is_visible = ?", true).Where("is_visible = ?", true).Find(&skillGroups)

	var experiences []models.Experience
	s.db.Where("is_visible = ?", true).Find(&experiences)

	var categories []models.ProjectCategory
	s.db.Where("is_visible = ?", true).Find(&categories)

	var projects []models.Project
	s.db.Where("is_visible = ?", true).Find(&projects)

	var reviews []models.Review
	s.db.Where("is_visible = ?", true).Find(&reviews)

	var timelineData []dto.TimelineItem
	for _, exp := range experiences {
		item := dto.TimelineItem{
			ID:          exp.ID.String(),
			Type:        exp.Type,
			Period:      exp.StartDate.Time.Format("Jan 2006") + " - ",
			Year:        exp.StartDate.Time.Year(),
			Description: exp.Description,
		}
		if exp.EndDate != nil {
			item.Period += exp.EndDate.Time.Format("Jan 2006")
		} else {
			item.Period += "Present"
		}

		if exp.Type == "education" {
			item.Degree = exp.Title
			item.Institution = exp.Organization
		} else {
			item.Title = exp.Title
			item.Company = exp.Organization
		}
		timelineData = append(timelineData, item)
	}

	var formattedSkillGroups []dto.SkillGroupData
	for _, group := range skillGroups {
		var mappedSkills []dto.SkillData
		for _, skill := range group.Skills {
			mappedSkills = append(mappedSkills, dto.SkillData{
				ID:          skill.ID.String(),
				Name:        skill.Name,
				Proficiency: skill.Proficiency,
				Color:       skill.Color,
				Icon:        skill.Icon,
			})
		}
		if mappedSkills == nil {
			mappedSkills = []dto.SkillData{}
		}
		formattedSkillGroups = append(formattedSkillGroups, dto.SkillGroupData{
			Name:   group.Name,
			Skills: mappedSkills,
		})
	}

	var formattedProjects []dto.ProjectData
	for _, p := range projects {
		catID := ""
		if p.ProjectCategoryID != nil {
			catID = p.ProjectCategoryID.String()
		}
		formattedProjects = append(formattedProjects, dto.ProjectData{
			ID:                p.ID.String(),
			Type:              p.Type,
			Title:             p.Title,
			Description:       p.Description,
			Image:             p.Image,
			Tags:              p.Tags,
			URL:               p.URL,
			AdditionalData:    p.AdditionalData,
			CreatedAt:         p.CreatedAt,
			ProjectCategoryID: catID,
		})
	}

	var featuredSkills []dto.SkillData
	if len(user.FeaturedSkillIDs) > 0 {
		var fSkills []models.Skill
		s.db.Where("id IN ?", []string(user.FeaturedSkillIDs)).Find(&fSkills)
		for _, fs := range fSkills {
			featuredSkills = append(featuredSkills, dto.SkillData{
				ID:          fs.ID.String(),
				Name:        fs.Name,
				Proficiency: fs.Proficiency,
				Color:       fs.Color,
				Icon:        fs.Icon,
			})
		}
	}
	if featuredSkills == nil {
		featuredSkills = []dto.SkillData{}
	}
	if timelineData == nil {
		timelineData = []dto.TimelineItem{}
	}
	if formattedProjects == nil {
		formattedProjects = []dto.ProjectData{}
	}
	if categories == nil {
		categories = []models.ProjectCategory{}
	}
	if reviews == nil {
		reviews = []models.Review{}
	}

	totalExp := calculateTotalExperience(experiences)

	loc := ""
	if user.Location != nil && *user.Location != "" {
		loc = "Based in " + *user.Location
	}

	aboutDesc := ""
	aboutShort := ""
	var aboutImage *string
	if user.About != nil {
		if d, ok := user.About["description"].(string); ok {
			aboutDesc = d
		}
		if sd, ok := user.About["shortdescription"].(string); ok {
			aboutShort = sd
		}
		if img, ok := user.About["image"].(string); ok {
			aboutImage = &img
		}
	}

	return &dto.PublicDataResponse{
		Name:         user.Name,
		Surname:      user.Surname,
		Title:        user.Title,
		Email:        user.Email,
		Phone:        user.Phone,
		Location:     loc,
		Availability: user.Availability,
		Avatar:       user.Avatar,
		HeroStats: map[string]interface{}{
			"experience": totalExp,
		},
		SocialLinks:    []map[string]interface{}(user.SocialLinks),
		FeaturedSkills: featuredSkills,
		About: map[string]interface{}{
			"title":            "More about",
			"highlight":        "Myself",
			"subtitle":         "About",
			"description":      aboutDesc,
			"shortdescription": aboutShort,
			"image":            aboutImage,
		},
		ProjectsSection: map[string]interface{}{
			"subtitle":  "Projects",
			"title":     "My",
			"highlight": "Projects",
		},
		SkillsSection: map[string]interface{}{
			"subtitle":  "Skills",
			"title":     "My",
			"highlight": "Skills",
		},
		TimelineSection: map[string]interface{}{
			"subtitle":  "Experience & Education",
			"title":     "My",
			"highlight": "Experience & Education",
		},
		SkillGroups:       formattedSkillGroups,
		TimelineData:      timelineData,
		ProjectCategories: categories,
		Projects:          formattedProjects,
		Reviews:           reviews,
	}, nil
}
