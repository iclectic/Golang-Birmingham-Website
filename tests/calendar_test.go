package tests

import (
	"strings"
	"testing"
	"time"

	"github.com/placeholder/golang-birmingham-website/internal/calendar"
	"github.com/placeholder/golang-birmingham-website/internal/domain"
)

func TestCalendar_GenerateICS(t *testing.T) {
	// Create test event
	event := &domain.Event{
		ID:          "test-event-001",
		Slug:        "test-event",
		Title:       "Test Go Event",
		Summary:     "A test Go programming event",
		Description: "This is a detailed description of our test Go event.\n\nWe will cover Go basics.",
		StartTime:   time.Date(2027, 7, 25, 18, 30, 0, 0, time.UTC),
		EndTime:     time.Date(2027, 7, 25, 21, 0, 0, 0, time.UTC),
		VenueName:   "Test Venue",
		VenueAddress: "123 Test Street, Birmingham, B1 1AA",
		Speakers: []domain.Speaker{
			{
				Name:         "Test Speaker",
				Role:         "Developer",
				Organisation: "Test Company",
			},
		},
	}

	baseURL := "http://localhost:8080"

	// Generate ICS
	icsData, err := calendar.GenerateICS(event, baseURL)
	if err != nil {
		t.Fatalf("Failed to generate ICS: %v", err)
	}

	// Convert to string for testing
	icsString := string(icsData)

	// Check ICS header
	expectedHeaders := []string{
		"BEGIN:VCALENDAR",
		"VERSION:2.0",
		"PRODID:-//Golang Birmingham//Community Events//EN",
		"BEGIN:VEVENT",
		"END:VEVENT",
		"END:VCALENDAR",
	}

	for _, header := range expectedHeaders {
		if !strings.Contains(icsString, header) {
			t.Errorf("ICS data missing header: %s", header)
		}
	}

	// Check event details
	expectedDetails := []string{
		"UID:test-event-001@golangbirmingham.com",
		"DTSTART:20270725T183000Z",
		"DTEND:20270725T210000Z",
		"SUMMARY:Test Go Event",
		"LOCATION:Test Venue, 123 Test Street, Birmingham, B1 1AA",
		"URL:http://localhost:8080/events/test-event",
		"STATUS:CONFIRMED",
	}

	for _, detail := range expectedDetails {
		if !strings.Contains(icsString, detail) {
			t.Errorf("ICS data missing detail: %s", detail)
		}
	}

	// Check description contains summary and speaker info
	if !strings.Contains(icsString, event.Summary) {
		t.Errorf("ICS description missing summary: %s", event.Summary)
	}

	if !strings.Contains(icsString, "Test Speaker") {
		t.Errorf("ICS description missing speaker: Test Speaker")
	}
}

func TestCalendar_EscapeICSString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple string",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "String with comma",
			input:    "Hello, World",
			expected: "Hello\\, World",
		},
		{
			name:     "String with semicolon",
			input:    "Hello; World",
			expected: "Hello\\; World",
		},
		{
			name:     "String with backslash",
			input:    "Hello\\World",
			expected: "Hello\\\\World",
		},
		{
			name:     "String with newline",
			input:    "Hello\nWorld",
			expected: "Hello\\nWorld",
		},
		{
			name:     "Complex string",
			input:    "Hello, World; Test\\Newline",
			expected: "Hello\\, World\\; Test\\\\Newline",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calendar.EscapeICSString(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestCalendar_GenerateICSWithMultipleSpeakers(t *testing.T) {
	event := &domain.Event{
		ID:          "multi-speaker-event",
		Slug:        "multi-speaker",
		Title:       "Multi-Speaker Event",
		Summary:     "Event with multiple speakers",
		Description: "This event has multiple speakers.",
		StartTime:   time.Date(2027, 8, 15, 19, 0, 0, 0, time.UTC),
		EndTime:     time.Date(2027, 8, 15, 21, 0, 0, 0, time.UTC),
		VenueName:   "Test Venue",
		VenueAddress: "456 Test Avenue, Birmingham, B2 2BB",
		Speakers: []domain.Speaker{
			{
				Name:         "Alice Smith",
				Role:         "Senior Developer",
				Organisation: "Tech Corp",
			},
			{
				Name:         "Bob Johnson",
				Role:         "DevOps Engineer",
				Organisation: "Cloud Solutions",
			},
		},
	}

	baseURL := "http://localhost:8080"

	icsData, err := calendar.GenerateICS(event, baseURL)
	if err != nil {
		t.Fatalf("Failed to generate ICS: %v", err)
	}

	icsString := string(icsData)

	// Check that both speakers are included
	if !strings.Contains(icsString, "Alice Smith") {
		t.Errorf("ICS description missing speaker: Alice Smith")
	}

	if !strings.Contains(icsString, "Bob Johnson") {
		t.Errorf("ICS description missing speaker: Bob Johnson")
	}

	if !strings.Contains(icsString, "Senior Developer") {
		t.Errorf("ICS description missing role: Senior Developer")
	}

	if !strings.Contains(icsString, "DevOps Engineer") {
		t.Errorf("ICS description missing role: DevOps Engineer")
	}
}

func TestCalendar_GenerateICSWithEmptySpeakers(t *testing.T) {
	event := &domain.Event{
		ID:          "no-speakers-event",
		Slug:        "no-speakers",
		Title:       "No Speakers Event",
		Summary:     "Event without speakers",
		Description: "This event has no speakers listed.",
		StartTime:   time.Date(2027, 9, 10, 18, 0, 0, 0, time.UTC),
		EndTime:     time.Date(2027, 9, 10, 20, 0, 0, 0, time.UTC),
		VenueName:   "Test Venue",
		VenueAddress: "789 Test Road, Birmingham, B3 3CC",
		Speakers:    []domain.Speaker{},
	}

	baseURL := "http://localhost:8080"

	icsData, err := calendar.GenerateICS(event, baseURL)
	if err != nil {
		t.Fatalf("Failed to generate ICS: %v", err)
	}

	icsString := string(icsData)

	// Check that basic event details are still present
	expectedDetails := []string{
		"UID:no-speakers-event@golangbirmingham.com",
		"SUMMARY:No Speakers Event",
		"LOCATION:Test Venue, 789 Test Road, Birmingham, B3 3CC",
	}

	for _, detail := range expectedDetails {
		if !strings.Contains(icsString, detail) {
			t.Errorf("ICS data missing detail: %s", detail)
		}
	}

	// Should not contain "Speakers:" section
	if strings.Contains(icsString, "Speakers:") {
		t.Errorf("ICS should not contain 'Speakers:' section when no speakers are present")
	}
}

func TestCalendar_GenerateICSWithSpecialCharacters(t *testing.T) {
	event := &domain.Event{
		ID:          "special-chars-event",
		Slug:        "special-chars",
		Title:       "Event with Special Characters: & < > \" '",
		Summary:     "Summary with special chars: semicolon; comma, backslash\\",
		Description: "Description with newlines\nand special characters.",
		StartTime:   time.Date(2027, 10, 5, 17, 30, 0, 0, time.UTC),
		EndTime:     time.Date(2027, 10, 5, 19, 30, 0, 0, time.UTC),
		VenueName:   "Venue & Location",
		VenueAddress: "123 Main St; Suite 456, Birmingham, B1 1AA",
		Speakers: []domain.Speaker{
			{
				Name:         "John O'Connor",
				Role:         "Lead Developer",
				Organisation: "Company & Co.",
			},
		},
	}

	baseURL := "http://localhost:8080"

	icsData, err := calendar.GenerateICS(event, baseURL)
	if err != nil {
		t.Fatalf("Failed to generate ICS: %v", err)
	}

	icsString := string(icsData)

	// Check that special characters are properly escaped
	if !strings.Contains(icsString, "SUMMARY:Event with Special Characters: & < > \" '") {
		t.Errorf("ICS title not properly handled for special characters")
	}

	if !strings.Contains(icsString, "Venue \\& Location") {
		t.Errorf("ICS venue name not properly escaped: expected 'Venue \\& Location'")
	}

	if !strings.Contains(icsString, "123 Main St\\; Suite 456, Birmingham, B1 1AA") {
		t.Errorf("ICS venue address not properly escaped")
	}
}
