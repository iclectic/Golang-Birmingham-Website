// Package config reads runtime configuration from environment variables.
package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// Config is the server's runtime configuration.
type Config struct {
	Port        string
	BaseURL     string
	Environment string
	ContentDir  string
}

// Load reads configuration from the environment, applying defaults.
func Load() (*Config, error) {
	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		BaseURL:     strings.TrimRight(getEnv("BASE_URL", "http://localhost:8080"), "/"),
		Environment: getEnv("ENVIRONMENT", "development"),
		ContentDir:  getEnv("CONTENT_DIR", "content"),
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return cfg, nil
}

// Validate checks that the configuration is usable.
func (c *Config) Validate() error {
	if c.Port == "" {
		return fmt.Errorf("PORT cannot be empty")
	}
	u, err := url.Parse(c.BaseURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("BASE_URL must be an absolute URL, got %q", c.BaseURL)
	}
	if c.ContentDir == "" {
		return fmt.Errorf("CONTENT_DIR cannot be empty")
	}
	return nil
}

// IsProduction reports whether the server is running in production.
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
