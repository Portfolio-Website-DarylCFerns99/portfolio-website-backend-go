package dto

import (
	"time"

	"portfolio-website-backend/internal/models"
)

type TimelineItem struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	Title       string  `json:"title,omitempty"`
	Degree      string  `json:"degree,omitempty"`
	Company     string  `json:"company,omitempty"`
	Institution string  `json:"institution,omitempty"`
	Period      string  `json:"period"`
	Year        int     `json:"year"`
	Description *string `json:"description"`
}

type SkillData struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Proficiency int     `json:"proficiency"`
	Color       *string `json:"color"`
	Icon        *string `json:"icon"`
}

type SkillGroupData struct {
	Name   string      `json:"name"`
	Skills []SkillData `json:"skills"`
}

type ProjectData struct {
	ID                string                 `json:"id"`
	Type              string                 `json:"type"`
	Title             string                 `json:"title"`
	Description       *string                `json:"description"`
	Image             *string                `json:"image"`
	Tags              []string               `json:"tags"`
	URL               *string                `json:"url"`
	AdditionalData    map[string]interface{} `json:"additional_data"`
	CreatedAt         time.Time              `json:"created_at"`
	ProjectCategoryID string                 `json:"project_category_id"`
}

type PublicDataResponse struct {
	Name              *string                  `json:"name"`
	Surname           *string                  `json:"surname"`
	Title             *string                  `json:"title"`
	Email             string                   `json:"email"`
	Phone             *string                  `json:"phone"`
	Location          string                   `json:"location"`
	Availability      *string                  `json:"availability"`
	Avatar            *string                  `json:"avatar"`
	HeroStats         map[string]interface{}   `json:"heroStats"`
	SocialLinks       []map[string]interface{} `json:"socialLinks"`
	FeaturedSkills    []SkillData              `json:"featuredSkills"`
	About             map[string]interface{}   `json:"about"`
	ProjectsSection   map[string]interface{}   `json:"projectsSection"`
	SkillsSection     map[string]interface{}   `json:"skillsSection"`
	TimelineSection   map[string]interface{}   `json:"timelineSection"`
	SkillGroups       []SkillGroupData         `json:"skillGroups"`
	TimelineData      []TimelineItem           `json:"timelineData"`
	ProjectCategories []models.ProjectCategory `json:"projectCategories"`
	Projects          []ProjectData            `json:"projects"`
	Reviews           []models.Review          `json:"reviews"`
}
