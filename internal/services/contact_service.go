package services

import (
	"fmt"

	"github.com/google/uuid"

	"portfolio-website-backend/internal/config"
	"portfolio-website-backend/internal/dto"
	"portfolio-website-backend/internal/repository"
	"portfolio-website-backend/internal/utils"
)

// ContactService handles the business logic for the contact form.
type ContactService interface {
	SendContactEmail(userID uuid.UUID, req dto.ContactRequest) error
}

type contactService struct {
	userRepo repository.UserRepository
}

// NewContactService creates a new ContactService.
func NewContactService(userRepo repository.UserRepository) ContactService {
	return &contactService{userRepo: userRepo}
}

func (s *contactService) SendContactEmail(userID uuid.UUID, req dto.ContactRequest) error {
	// Fetch the portfolio owner's user record
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Build the full name for the email templates
	yourName := user.Username
	if user.Name != nil && user.Surname != nil {
		yourName = *user.Name + " " + *user.Surname
	} else if user.Name != nil {
		yourName = *user.Name
	}

	// Use the first CORS origin as the portfolio URL (same logic as the Python backend)
	portfolioURL := ""
	if len(config.Envs.CorsOrigins) > 0 && config.Envs.CorsOrigins[0] != "*" {
		portfolioURL = config.Envs.CorsOrigins[0]
	}

	mailer := utils.NewMailgunEmail()

	// Send confirmation email to the person who filled in the form
	if err := mailer.SendConfirmationEmail(
		req.Name,
		req.Email,
		req.Subject,
		req.Message,
		user.SocialLinks,
		portfolioURL,
		yourName,
	); err != nil {
		return fmt.Errorf("failed to send confirmation email: %w", err)
	}

	// Send notification email to the admin
	if err := mailer.SendNotificationEmail(
		req.Name,
		req.Email,
		req.Subject,
		req.Message,
		yourName,
	); err != nil {
		return fmt.Errorf("failed to send notification email: %w", err)
	}

	return nil
}
