package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"portfolio-website-backend/internal/config"
)

// MailgunEmail handles sending emails via the Mailgun API.
type MailgunEmail struct {
	baseURL              string
	apiKey               string
	fromEmail            string
	adminEmail           string
	notificationTemplate string
	confirmationTemplate string
}

// NewMailgunEmail creates a new MailgunEmail instance using the app config.
func NewMailgunEmail() *MailgunEmail {
	cfg := config.Envs
	return &MailgunEmail{
		baseURL:              cfg.MailgunAPIURL,
		apiKey:               cfg.MailgunAPIKey,
		fromEmail:            cfg.MailgunFromEmail,
		adminEmail:           cfg.AdminEmail,
		notificationTemplate: cfg.MailgunNotificationTemplateID,
		confirmationTemplate: cfg.MailgunConfirmationTemplateID,
	}
}

// formatSocialLinks filters and normalises social links for the email template.
// Only the platforms supported by the email template are kept.
func formatSocialLinks(socialLinks []map[string]interface{}) []map[string]string {
	validPlatforms := map[string]bool{
		"linkedin":  true,
		"github":    true,
		"twitter":   true,
		"instagram": true,
		"facebook":  true,
	}

	formatted := make([]map[string]string, 0, len(socialLinks))
	for _, link := range socialLinks {
		platform, _ := link["platform"].(string)
		platform = strings.ToLower(platform)
		if !validPlatforms[platform] {
			continue
		}
		linkURL, _ := link["url"].(string)
		if linkURL == "" {
			continue
		}
		tooltip, _ := link["tooltip"].(string)
		if tooltip == "" {
			tooltip = strings.ToTitle(platform)
		}
		formatted = append(formatted, map[string]string{
			"platform": platform,
			"url":      linkURL,
			"tooltip":  tooltip,
		})
	}
	return formatted
}

// sendMessage is the shared helper that posts a message to the Mailgun API.
func (m *MailgunEmail) sendMessage(to, subject, templateID string, templateVars map[string]interface{}) error {
	tvJSON, err := json.Marshal(templateVars)
	if err != nil {
		return fmt.Errorf("mailgun: failed to marshal template variables: %w", err)
	}

	formData := url.Values{}
	formData.Set("from", m.fromEmail)
	formData.Set("to", to)
	formData.Set("subject", subject)
	formData.Set("template", templateID)
	formData.Set("t:variables", string(tvJSON))

	req, err := http.NewRequest(http.MethodPost, m.baseURL+"/messages", strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("mailgun: failed to create request: %w", err)
	}
	req.SetBasicAuth("api", m.apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("mailgun: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("mailgun: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// SendConfirmationEmail sends a confirmation email to the person who submitted the contact form.
func (m *MailgunEmail) SendConfirmationEmail(
	name, email, subject, message string,
	socialLinks []map[string]interface{},
	portfolioURL, yourName string,
) error {
	templateVars := map[string]interface{}{
		"name":         name,
		"email":        email,
		"subject":      subject,
		"message":      message,
		"social_links": formatSocialLinks(socialLinks),
		"portfolio_url": portfolioURL,
		"your_name":    yourName,
	}
	return m.sendMessage(email, subject, m.confirmationTemplate, templateVars)
}

// SendNotificationEmail sends an admin notification email about a new contact form submission.
func (m *MailgunEmail) SendNotificationEmail(
	name, email, subject, message, yourName string,
) error {
	templateVars := map[string]interface{}{
		"name":      name,
		"email":     email,
		"subject":   subject,
		"message":   message,
		"your_name": yourName,
	}
	return m.sendMessage(m.adminEmail, subject, m.notificationTemplate, templateVars)
}
