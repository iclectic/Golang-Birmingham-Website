package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port           string
	BaseURL        string
	Environment    string
	ContactEmail   string
	SpeakerFormURL string
	SocialLinks    map[string]string
	Security       SecurityConfig
}

type SecurityConfig struct {
	CSRFSecret    string
	SessionSecret string
	EnableHTTPS   bool
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:           getEnv("PORT", "8080"),
		BaseURL:        getEnv("BASE_URL", "http://localhost:8080"),
		Environment:    getEnv("ENVIRONMENT", "development"),
		ContactEmail:   getEnv("CONTACT_EMAIL", "conduct@example.com"),
		SpeakerFormURL: getEnv("SPEAKER_FORM_URL", "https://forms.example.com/speaker"),
		SocialLinks: map[string]string{
			"meetup":    getEnv("MEETUP_URL", "https://meetup.com/golang-birmingham"),
			"linkedin":  getEnv("LINKEDIN_URL", "https://linkedin.com/company/golang-birmingham"),
			"youtube":   getEnv("YOUTUBE_URL", "https://youtube.com/@golangbirmingham"),
			"github":    getEnv("GITHUB_URL", "https://github.com/golang-birmingham"),
		},
		Security: SecurityConfig{
			CSRFSecret:    getEnv("CSRF_SECRET", "dev-secret-change-in-production"),
			SessionSecret: getEnv("SESSION_SECRET", "dev-session-secret-change-in-production"),
			EnableHTTPS:   getBoolEnv("ENABLE_HTTPS", false),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Port == "" {
		return fmt.Errorf("PORT cannot be empty")
	}
	if c.BaseURL == "" {
		return fmt.Errorf("BASE_URL cannot be empty")
	}
	if c.ContactEmail == "" {
		return fmt.Errorf("CONTACT_EMAIL cannot be empty")
	}
	return nil
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
