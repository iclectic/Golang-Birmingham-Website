package tests

import (
	"os"
	"testing"

	"github.com/placeholder/golang-birmingham-website/internal/config"
)

func TestConfig_Load(t *testing.T) {
	// Set test environment variables
	os.Setenv("PORT", "9090")
	os.Setenv("BASE_URL", "https://test.example.com")
	os.Setenv("ENVIRONMENT", "test")
	os.Setenv("CONTACT_EMAIL", "test@example.com")
	os.Setenv("CSRF_SECRET", "test-secret")
	os.Setenv("SESSION_SECRET", "test-session-secret")
	defer func() {
		// Clean up environment variables
		os.Unsetenv("PORT")
		os.Unsetenv("BASE_URL")
		os.Unsetenv("ENVIRONMENT")
		os.Unsetenv("CONTACT_EMAIL")
		os.Unsetenv("CSRF_SECRET")
		os.Unsetenv("SESSION_SECRET")
	}()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check that values are loaded correctly
	if cfg.Port != "9090" {
		t.Errorf("Expected port 9090, got %s", cfg.Port)
	}

	if cfg.BaseURL != "https://test.example.com" {
		t.Errorf("Expected base URL https://test.example.com, got %s", cfg.BaseURL)
	}

	if cfg.Environment != "test" {
		t.Errorf("Expected environment test, got %s", cfg.Environment)
	}

	if cfg.ContactEmail != "test@example.com" {
		t.Errorf("Expected contact email test@example.com, got %s", cfg.ContactEmail)
	}
}

func TestConfig_LoadDefaults(t *testing.T) {
	// Clear relevant environment variables
	os.Unsetenv("PORT")
	os.Unsetenv("BASE_URL")
	os.Unsetenv("ENVIRONMENT")
	os.Unsetenv("CONTACT_EMAIL")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("BASE_URL")
		os.Unsetenv("ENVIRONMENT")
		os.Unsetenv("CONTACT_EMAIL")
	}()

	// Load configuration with defaults
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check defaults
	if cfg.Port != "8080" {
		t.Errorf("Expected default port 8080, got %s", cfg.Port)
	}

	if cfg.BaseURL != "http://localhost:8080" {
		t.Errorf("Expected default base URL http://localhost:8080, got %s", cfg.BaseURL)
	}

	if cfg.Environment != "development" {
		t.Errorf("Expected default environment development, got %s", cfg.Environment)
	}
}

func TestConfig_Validation(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
	}{
		{
			name: "Valid configuration",
			envVars: map[string]string{
				"PORT":          "8080",
				"BASE_URL":      "http://localhost:8080",
				"CONTACT_EMAIL": "test@example.com",
				"CSRF_SECRET":   "valid-secret",
				"SESSION_SECRET": "valid-session-secret",
			},
			wantErr: false,
		},
		{
			name: "Missing CSRF secret",
			envVars: map[string]string{
				"PORT":          "8080",
				"BASE_URL":      "http://localhost:8080",
				"CONTACT_EMAIL": "test@example.com",
				"CSRF_SECRET":   "",
				"SESSION_SECRET": "valid-session-secret",
			},
			wantErr: true,
		},
		{
			name: "Missing session secret",
			envVars: map[string]string{
				"PORT":          "8080",
				"BASE_URL":      "http://localhost:8080",
				"CONTACT_EMAIL": "test@example.com",
				"CSRF_SECRET":   "valid-secret",
				"SESSION_SECRET": "",
			},
			wantErr: true,
		},
		{
			name: "Invalid port",
			envVars: map[string]string{
				"PORT":          "invalid",
				"BASE_URL":      "http://localhost:8080",
				"CONTACT_EMAIL": "test@example.com",
				"CSRF_SECRET":   "valid-secret",
				"SESSION_SECRET": "valid-session-secret",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}
			defer func() {
				// Clean up environment variables
				for key := range tt.envVars {
					os.Unsetenv(key)
				}
			}()

			// Try to load configuration
			_, err := config.Load()

			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
