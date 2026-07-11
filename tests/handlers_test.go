package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/placeholder/golang-birmingham-website/internal/config"
	"github.com/placeholder/golang-birmingham-website/internal/content"
	"github.com/placeholder/golang-birmingham-website/internal/domain"
	"github.com/placeholder/golang-birmingham-website/internal/handlers"
)

func TestHandlers_Home(t *testing.T) {
	// Setup test data
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
		Port:    "8080",
	}

	testContent := &content.Content{
		Site: domain.SiteConfig{
			Title:       "Golang Birmingham",
			Description: "Test description",
			Stats: domain.SiteStats{
				Members:  100,
				Events:   10,
				Speakers: 20,
				Talks:    30,
			},
		},
		Events: []domain.Event{
			{
				ID:        "test-event",
				Slug:      "test-event",
				Title:     "Test Event",
				Summary:   "Test summary",
				StartTime: time.Now().Add(24 * time.Hour),
				EndTime:   time.Now().Add(26 * time.Hour),
				VenueName: "Test Venue",
				Image:     "/static/images/events/test.jpg",
			},
		},
	}

	// Create handlers
	h := handlers.New(cfg, testContent)

	// Create request
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Call handler
	h.Home(w, req)

	// Check response
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check content type
	expectedContentType := "text/html; charset=utf-8"
	if resp.Header.Get("Content-Type") != expectedContentType {
		t.Errorf("Expected content type %s, got %s", expectedContentType, resp.Header.Get("Content-Type"))
	}
}

func TestHandlers_Events(t *testing.T) {
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
		Port:    "8080",
	}

	testContent := &content.Content{
		Site: domain.SiteConfig{
			Title:       "Golang Birmingham",
			Description: "Test description",
		},
		Events: []domain.Event{
			{
				ID:        "test-event",
				Slug:      "test-event",
				Title:     "Test Event",
				Summary:   "Test summary",
				StartTime: time.Now().Add(24 * time.Hour),
				EndTime:   time.Now().Add(26 * time.Hour),
				VenueName: "Test Venue",
				Image:     "/static/images/events/test.jpg",
			},
		},
	}

	h := handlers.New(cfg, testContent)

	req := httptest.NewRequest("GET", "/events", nil)
	w := httptest.NewRecorder()

	h.Events(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestHandlers_EventDetail(t *testing.T) {
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
		Port:    "8080",
	}

	testContent := &content.Content{
		Site: domain.SiteConfig{
			Title:       "Golang Birmingham",
			Description: "Test description",
		},
		Events: []domain.Event{
			{
				ID:        "test-event",
				Slug:      "test-event",
				Title:     "Test Event",
				Summary:   "Test summary",
				StartTime: time.Now().Add(24 * time.Hour),
				EndTime:   time.Now().Add(26 * time.Hour),
				VenueName: "Test Venue",
				Image:     "/static/images/events/test.jpg",
			},
		},
	}

	h := handlers.New(cfg, testContent)

	tests := []struct {
		name           string
		url            string
		expectedStatus int
	}{
		{
			name:           "Valid event",
			url:            "/events/test-event",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid event",
			url:            "/events/nonexistent",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Empty path",
			url:            "/events/",
			expectedStatus: http.StatusOK, // Should redirect to events list
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()

			h.EventDetail(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func TestHandlers_Health(t *testing.T) {
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
		Port:    "8080",
	}

	testContent := &content.Content{
		Site: domain.SiteConfig{
			Title: "Golang Birmingham",
		},
	}

	h := handlers.New(cfg, testContent)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	h.Health(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	expectedContentType := "application/json"
	if resp.Header.Get("Content-Type") != expectedContentType {
		t.Errorf("Expected content type %s, got %s", expectedContentType, resp.Header.Get("Content-Type"))
	}
}

func TestHandlers_Robots(t *testing.T) {
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
		Port:    "8080",
	}

	testContent := &content.Content{
		Site: domain.SiteConfig{
			Title: "Golang Birmingham",
		},
	}

	h := handlers.New(cfg, testContent)

	req := httptest.NewRequest("GET", "/robots.txt", nil)
	w := httptest.NewRecorder()

	h.Robots(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	expectedContentType := "text/plain"
	if resp.Header.Get("Content-Type") != expectedContentType {
		t.Errorf("Expected content type %s, got %s", expectedContentType, resp.Header.Get("Content-Type"))
	}
}

func TestHandlers_Sitemap(t *testing.T) {
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
		Port:    "8080",
	}

	testContent := &content.Content{
		Site: domain.SiteConfig{
			Title: "Golang Birmingham",
		},
	}

	h := handlers.New(cfg, testContent)

	req := httptest.NewRequest("GET", "/sitemap.xml", nil)
	w := httptest.NewRecorder()

	h.Sitemap(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	expectedContentType := "application/xml"
	if resp.Header.Get("Content-Type") != expectedContentType {
		t.Errorf("Expected content type %s, got %s", expectedContentType, resp.Header.Get("Content-Type"))
	}
}

func TestHandlers_Contact(t *testing.T) {
	cfg := &config.Config{
		BaseURL:      "http://localhost:8080",
		Port:         "8080",
		ContactEmail: "test@example.com",
	}

	testContent := &content.Content{
		Site: domain.SiteConfig{
			Title: "Golang Birmingham",
		},
	}

	h := handlers.New(cfg, testContent)

	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "GET request",
			method:         "GET",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST request",
			method:         "POST",
			expectedStatus: http.StatusOK, // Should handle form submission
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/contact", nil)
			w := httptest.NewRecorder()

			h.Contact(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func TestHandlers_NotFound(t *testing.T) {
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
		Port:    "8080",
	}

	testContent := &content.Content{
		Site: domain.SiteConfig{
			Title: "Golang Birmingham",
		},
	}

	h := handlers.New(cfg, testContent)

	req := httptest.NewRequest("GET", "/nonexistent-page", nil)
	w := httptest.NewRecorder()

	h.NotFound(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

// Table-driven test for static pages
func TestHandlers_StaticPages(t *testing.T) {
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
		Port:    "8080",
	}

	testContent := &content.Content{
		Site: domain.SiteConfig{
			Title: "Golang Birmingham",
		},
	}

	h := handlers.New(cfg, testContent)

	tests := []struct {
		name           string
		handler        func(http.ResponseWriter, *http.Request)
		url            string
		expectedStatus int
	}{
		{
			name:           "About page",
			handler:        h.About,
			url:            "/about",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Speakers page",
			handler:        h.Speakers,
			url:            "/speakers",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Team page",
			handler:        h.Team,
			url:            "/team",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Gallery page",
			handler:        h.Gallery,
			url:            "/gallery",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Code of Conduct page",
			handler:        h.CodeOfConduct,
			url:            "/code-of-conduct",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Privacy page",
			handler:        h.Privacy,
			url:            "/privacy",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()

			tt.handler(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}
