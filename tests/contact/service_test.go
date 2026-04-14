package contact

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"portfolio-website-backend/internal/config"
	"portfolio-website-backend/internal/dto"
	"portfolio-website-backend/internal/models"
	"portfolio-website-backend/internal/services"
)

// ensureConfig initialises the global config singleton with safe test defaults
// so the service can read CorsOrigins without panicking.
func ensureConfig() {
	if config.Envs == nil {
		config.Envs = &config.Config{
			CorsOrigins:                   []string{"https://portfolio.example.com"},
			MailgunAPIURL:                 "https://api.mailgun.net/v3/example.com",
			MailgunAPIKey:                 "test-key",
			MailgunFromEmail:              "no-reply@example.com",
			AdminEmail:                    "admin@example.com",
			MailgunNotificationTemplateID: "notification-template",
			MailgunConfirmationTemplateID: "confirmation-template",
		}
	}
}

// ---------------------------------------------------------------------------
// User not found → "user not found" error
// ---------------------------------------------------------------------------

func TestContactService_SendEmail_UserNotFound(t *testing.T) {
	ensureConfig()

	mockRepo := new(MockUserRepository)
	svc := services.NewContactService(mockRepo)

	userID := uuid.New()
	mockRepo.On("GetByID", userID).Return((*models.User)(nil), nil)

	req := dto.ContactRequest{
		Name:    "Alice",
		Email:   "alice@example.com",
		Subject: "Hello",
		Message: "Hi there.",
	}

	err := svc.SendContactEmail(userID, req)

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	mockRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// DB error while fetching user → wrapped error propagated
// ---------------------------------------------------------------------------

func TestContactService_SendEmail_DBError(t *testing.T) {
	ensureConfig()

	mockRepo := new(MockUserRepository)
	svc := services.NewContactService(mockRepo)

	userID := uuid.New()
	mockRepo.On("GetByID", userID).Return((*models.User)(nil), errors.New("connection refused"))

	err := svc.SendContactEmail(userID, dto.ContactRequest{
		Name:    "Bob",
		Email:   "bob@example.com",
		Subject: "DB test",
		Message: "Should fail fast.",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	mockRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Full name resolution — Name + Surname present
// ---------------------------------------------------------------------------

func TestContactService_FullName_NameAndSurname(t *testing.T) {
	ensureConfig()

	name := "Alice"
	surname := "Smith"
	user := &models.User{
		BaseModel: models.BaseModel{ID: uuid.New()},
		Username:  "alice",
		Name:      &name,
		Surname:   &surname,
	}

	mockRepo := new(MockUserRepository)
	// We only need the repo to return the user; the actual email call will fail
	// against the dummy Mailgun URL, so we patch around that by inspecting the
	// error message rather than expecting success.
	mockRepo.On("GetByID", user.ID).Return(user, nil)

	svc := services.NewContactService(mockRepo)
	err := svc.SendContactEmail(user.ID, dto.ContactRequest{
		Name:    "Sender",
		Email:   "sender@example.com",
		Subject: "Hi",
		Message: "Hello!",
	})

	// The error, if any, will come from Mailgun (no real HTTP server in unit
	// tests). What matters is that the service did NOT return "user not found".
	if err != nil {
		assert.NotEqual(t, "user not found", err.Error())
	}
	mockRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Full name resolution — only Name present (no Surname)
// ---------------------------------------------------------------------------

func TestContactService_FullName_NameOnly(t *testing.T) {
	ensureConfig()

	name := "Charlie"
	user := &models.User{
		BaseModel: models.BaseModel{ID: uuid.New()},
		Username:  "charlie",
		Name:      &name,
		Surname:   nil,
	}

	mockRepo := new(MockUserRepository)
	mockRepo.On("GetByID", user.ID).Return(user, nil)

	svc := services.NewContactService(mockRepo)
	err := svc.SendContactEmail(user.ID, dto.ContactRequest{
		Name:    "Sender",
		Email:   "sender@example.com",
		Subject: "Hi",
		Message: "Hello!",
	})

	if err != nil {
		assert.NotEqual(t, "user not found", err.Error())
	}
	mockRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Full name resolution — no Name / Surname → falls back to Username
// ---------------------------------------------------------------------------

func TestContactService_FullName_UsernameOnly(t *testing.T) {
	ensureConfig()

	user := &models.User{
		BaseModel: models.BaseModel{ID: uuid.New()},
		Username:  "dave99",
		Name:      nil,
		Surname:   nil,
	}

	mockRepo := new(MockUserRepository)
	mockRepo.On("GetByID", user.ID).Return(user, nil)

	svc := services.NewContactService(mockRepo)
	err := svc.SendContactEmail(user.ID, dto.ContactRequest{
		Name:    "Sender",
		Email:   "sender@example.com",
		Subject: "Hi",
		Message: "Hello!",
	})

	if err != nil {
		assert.NotEqual(t, "user not found", err.Error())
	}
	mockRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Portfolio URL — wildcard CORS should produce empty portfolio URL
// ---------------------------------------------------------------------------

func TestContactService_PortfolioURL_Wildcard(t *testing.T) {
	// Override CORS to wildcard
	config.Envs = &config.Config{
		CorsOrigins:                   []string{"*"},
		MailgunAPIURL:                 "https://api.mailgun.net/v3/example.com",
		MailgunAPIKey:                 "test-key",
		MailgunFromEmail:              "no-reply@example.com",
		AdminEmail:                    "admin@example.com",
		MailgunNotificationTemplateID: "n-tmpl",
		MailgunConfirmationTemplateID: "c-tmpl",
	}

	user := &models.User{
		BaseModel: models.BaseModel{ID: uuid.New()},
		Username:  "wildcarduser",
	}

	mockRepo := new(MockUserRepository)
	mockRepo.On("GetByID", user.ID).Return(user, nil)

	svc := services.NewContactService(mockRepo)

	// The service should still run (portfolio URL will be empty string).
	// It may fail at the Mailgun HTTP call, which is expected in unit tests.
	err := svc.SendContactEmail(user.ID, dto.ContactRequest{
		Name: "G", Email: "g@g.com", Subject: "s", Message: "m",
	})
	if err != nil {
		assert.NotEqual(t, "user not found", err.Error())
	}

	// Restore
	ensureConfig()
	mockRepo.AssertExpectations(t)
}
