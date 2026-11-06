package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	for _, k := range []string{"PORT", "BASE_URL", "ENVIRONMENT", "CONTENT_DIR"} {
		t.Setenv(k, "")
	}
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != "8080" || cfg.BaseURL != "http://localhost:8080" || cfg.ContentDir != "content" || cfg.IsProduction() {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("PORT", "9000")
	t.Setenv("BASE_URL", "https://golangbirmingham.com/")
	t.Setenv("ENVIRONMENT", "production")
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != "9000" || cfg.BaseURL != "https://golangbirmingham.com" || !cfg.IsProduction() {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestValidateRejectsRelativeBaseURL(t *testing.T) {
	cfg := &Config{Port: "8080", BaseURL: "golangbirmingham.com", ContentDir: "content"}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for relative BASE_URL")
	}
}
