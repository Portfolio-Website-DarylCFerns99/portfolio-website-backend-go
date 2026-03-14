package config

import (
	"log"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	// Database settings
	DatabaseURL string `env:"DATABASE_URL"`

	// API settings
	APIPrefix string `env:"API_PREFIX" envDefault:"/api/v1"`
	Debug     bool   `env:"DEBUG" envDefault:"false"`

	// Retry settings
	MaxDBRetries int     `env:"MAX_DB_RETRIES" envDefault:"3"`
	RetryBackoff float64 `env:"RETRY_BACKOFF" envDefault:"0.5"`

	// Authentication settings
	JWTSecretKey           string `env:"JWT_SECRET_KEY" envDefault:"YOUR_DEFAULT_SECRET_KEY_CHANGE_THIS"`
	AccessTokenExpireMins  int    `env:"ACCESS_TOKEN_EXPIRE_MINUTES" envDefault:"30"`

	// Mailgun settings
	MailgunAPIURL                  string `env:"MAILGUN_API_URL"`
	MailgunAPIKey                  string `env:"MAILGUN_API_KEY"`
	MailgunFromEmail               string `env:"MAILGUN_FROM_EMAIL" envDefault:"your-email@domain.com"`
	AdminEmail                     string `env:"ADMIN_EMAIL" envDefault:"your-email@domain.com"`
	MailgunNotificationTemplateID  string `env:"MAILGUN_NOTIFICATION_TEMPLATE_ID" envDefault:"test-id"`
	MailgunConfirmationTemplateID  string `env:"MAILGUN_CONFIRMATION_TEMPLATE_ID" envDefault:"test-id"`

	// LLM settings
	GeminiAPIKey          string `env:"GEMINI_API_KEY"`
	GeminiModel           string `env:"GEMINI_MODEL"`
	GeminiEmbeddingModel  string `env:"GEMINI_EMBEDDING_MODEL"`

	// CORS settings
	CorsOriginsRaw string   `env:"CORS_ORIGINS" envDefault:"*"`
	CorsOrigins    []string `env:"-"`
}

var Envs *Config

func LoadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Unable to parse environment variables: %e", err)
	}

	// Parse CORS origins
	if cfg.CorsOriginsRaw == "*" {
		cfg.CorsOrigins = []string{"*"}
	} else {
		splitOrigins := strings.Split(cfg.CorsOriginsRaw, ",")
		for _, origin := range splitOrigins {
			trimmed := strings.TrimSpace(origin)
			if trimmed != "" {
				cfg.CorsOrigins = append(cfg.CorsOrigins, trimmed)
			}
		}
	}

	Envs = &cfg
	return Envs
}
