package services

import (
	"context"
	"fmt"
	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/repository"
	"portfolio-website-backend/internal/utils"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type VectorService interface {
	SyncUserData(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
	Search(ctx context.Context, text string, userID uuid.UUID, limit int, filters []string) ([]models.VectorEmbedding, error)
}

type vectorService struct {
	db         *gorm.DB
	vectorRepo *repository.VectorRepository
	llmFactory *utils.LLMFactory
}

func NewVectorService(db *gorm.DB) VectorService {
	return &vectorService{
		db:         db,
		vectorRepo: repository.NewVectorRepository(db),
		llmFactory: utils.NewLLMFactory(),
	}
}

func (s *vectorService) SyncUserData(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	// First, clear existing vectors for user
	err := s.vectorRepo.ClearAllVectors(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to clear vectors: %v", err)
	}

	totalCount := 0

	// 1. Projects
	var projects []models.Project
	s.db.Where("user_id = ? AND is_visible = ?", userID, true).Find(&projects)
	for _, p := range projects {
		desc := ""
		if p.Description != nil {
			desc = *p.Description
		}
		content := fmt.Sprintf("Project: %s. Type: %s. Description: %s.", p.Title, p.Type, desc)
		err := s.saveVector(ctx, content, "project", &p.ID, userID)
		if err == nil {
			totalCount++
		}
	}

	// 2. Skills
	var skills []models.Skill
	s.db.Where("user_id = ? AND is_visible = ?", userID, true).Find(&skills)
	for _, sk := range skills {
		groupName := "Unknown"
		var sg models.SkillGroup
		if err := s.db.First(&sg, "id = ?", sk.SkillGroupID).Error; err == nil {
			groupName = sg.Name
		}
		content := fmt.Sprintf("Skill: %s. Proficiency: %d/5. Group: %s.", sk.Name, sk.Proficiency, groupName)
		err := s.saveVector(ctx, content, "skill", &sk.ID, userID)
		if err == nil {
			totalCount++
		}
	}

	// 3. Experience & Education
	var experiences []models.Experience
	s.db.Where("user_id = ? AND is_visible = ?", userID, true).Find(&experiences)
	for _, e := range experiences {
		typeStr := "Work Experience"
		if e.Type == "education" {
			typeStr = "Education"
		}
		
		endDate := "Present"
		if e.EndDate != nil && !e.EndDate.Time.IsZero() {
			endDate = e.EndDate.Time.Format("2006-01-02")
		}

		desc := ""
		if e.Description != nil {
			desc = *e.Description
		}

		content := fmt.Sprintf("%s: %s at %s. %s to %s. %s", typeStr, e.Title, e.Organization, e.StartDate.Time.Format("2006-01-02"), endDate, desc)
		err := s.saveVector(ctx, content, "experience", &e.ID, userID)
		if err == nil {
			totalCount++
		}
	}

	// 4. Reviews
	var reviews []models.Review
	s.db.Where("user_id = ? AND is_visible = ?", userID, true).Find(&reviews)
	for _, r := range reviews {
		whereKnown := "Client"
		if r.Company != nil {
			whereKnown = *r.Company
		}
		content := fmt.Sprintf("Review from %s (%s): '%s'. Rating: %d/5.", r.Name, whereKnown, r.Content, r.Rating)
		err := s.saveVector(ctx, content, "review", &r.ID, userID)
		if err == nil {
			totalCount++
		}
	}

	// 5. User Profile
	var user models.User
	s.db.First(&user, "id = ?", userID)
	if user.ID != uuid.Nil {
		title := "Developer"
		if user.Title != nil {
			title = *user.Title
		}
		loc := ""
		if user.Location != nil {
			loc = *user.Location
		}
		about := ""
		if user.About != nil {
			if text, ok := user.About["text"].(string); ok {
				about = text
			} else if text, ok := user.About["content"].(string); ok {
				about = text
			}
		}
		
		name := "Portfolio Owner"
		if user.Name != nil {
			name = *user.Name
		}

		content := fmt.Sprintf("About Me (%s, Portfolio Owner): %s. Location: %s. Bio: %s. Contact: %s.", name, title, loc, about, user.Email)
		err := s.saveVector(ctx, content, "user", &user.ID, userID)
		if err == nil {
			totalCount++
		}
	}

	return map[string]interface{}{
		"status":         "success",
		"vectors_synced": totalCount,
	}, nil
}

func (s *vectorService) saveVector(ctx context.Context, content string, sourceType string, sourceID *uuid.UUID, userID uuid.UUID) error {
	vec, err := s.llmFactory.EmbedQuery(ctx, content, 768)
	if err != nil {
		return err
	}

	embedding := models.VectorEmbedding{
		Content:    content,
		Embedding:  pgvector.NewVector(vec),
		SourceType: sourceType,
		SourceID:   sourceID,
		UserID:     userID,
	}

	return s.vectorRepo.AddEmbedding(&embedding)
}

func (s *vectorService) Search(ctx context.Context, text string, userID uuid.UUID, limit int, filters []string) ([]models.VectorEmbedding, error) {
	vec, err := s.llmFactory.EmbedQuery(ctx, text, 768)
	if err != nil {
		return nil, err
	}
	return s.vectorRepo.Search(userID, vec, limit, filters)
}
