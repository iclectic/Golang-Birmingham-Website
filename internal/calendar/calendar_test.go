package calendar

import (
	"strings"
	"testing"
	"time"

	"github.com/iclectic/golang-birmingham-website/internal/domain"
)

func TestGenerateICS(t *testing.T) {
	start := time.Date(2026, 11, 26, 18, 30, 0, 0, time.UTC)
	event := domain.Event{
		ID: "workshop", Slug: "workshop", Title: "Go, Microservices; and more",
		Summary:   "A long summary that is definitely going to need folding because it is longer than seventy-five octets.",
		StartTime: start, EndTime: start.Add(150 * time.Minute),
		VenueName: "Tech Hub", VenueAddress: "Colmore Row, Birmingham",
		Speakers: []domain.Speaker{{Name: "Sarah Chen", Role: "Engineer", Organisation: "TechCorp"}},
	}
	ics := string(GenerateICS(event, "https://golangbirmingham.com", start))

	for _, want := range []string{
		"BEGIN:VCALENDAR\r\n",
		"UID:workshop@golangbirmingham.com\r\n",
		"DTSTART:20261126T183000Z\r\n",
		"DTEND:20261126T210000Z\r\n",
		`SUMMARY:Go\, Microservices\; and more` + "\r\n",
		`LOCATION:Tech Hub\, Colmore Row\, Birmingham` + "\r\n",
		"URL:https://golangbirmingham.com/events/workshop\r\n",
		"END:VCALENDAR\r\n",
	} {
		if !strings.Contains(ics, want) {
			t.Errorf("ICS missing %q", want)
		}
	}
	unfolded := strings.ReplaceAll(ics, "\r\n ", "")
	if !strings.Contains(unfolded, `- Sarah Chen (Engineer)\, TechCorp`) {
		t.Errorf("speaker missing from description:\n%s", unfolded)
	}
	for _, line := range strings.Split(ics, "\r\n") {
		if len(line) > 75 {
			t.Errorf("line longer than 75 octets: %q", line)
		}
	}
}

func TestEscapeICSString(t *testing.T) {
	tests := map[string]string{
		"Hello World":    "Hello World",
		"Hello, World":   `Hello\, World`,
		"Hello; World":   `Hello\; World`,
		`Hello\World`:    `Hello\\World`,
		"Hello\nWorld":   `Hello\nWorld`,
		"Hello\r\nWorld": `Hello\nWorld`,
	}
	for in, want := range tests {
		if got := EscapeICSString(in); got != want {
			t.Errorf("EscapeICSString(%q) = %q, want %q", in, got, want)
		}
	}
}
