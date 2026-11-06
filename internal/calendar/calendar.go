// Package calendar generates iCalendar (RFC 5545) files for events.
package calendar

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/iclectic/golang-birmingham-website/internal/domain"
)

const icsTime = "20060102T150405Z"

// GenerateICS renders a single-event VCALENDAR document.
func GenerateICS(event domain.Event, baseURL string, now time.Time) []byte {
	host := "golang-birmingham"
	if u, err := url.Parse(baseURL); err == nil && u.Hostname() != "" {
		host = u.Hostname()
	}

	description := event.Summary
	if len(event.Speakers) > 0 {
		description += "\n\nSpeakers:"
		for _, s := range event.Speakers {
			description += "\n- " + s.Name
			if s.Role != "" {
				description += " (" + s.Role + ")"
			}
			if s.Organisation != "" {
				description += ", " + s.Organisation
			}
		}
	}
	eventURL := baseURL + "/events/" + event.Slug
	description += "\n\n" + eventURL

	location := event.VenueName
	if event.VenueAddress != "" {
		location += ", " + event.VenueAddress
	}

	var b strings.Builder
	line := func(name, value string) { b.WriteString(fold(name + ":" + value)) }

	line("BEGIN", "VCALENDAR")
	line("VERSION", "2.0")
	line("PRODID", "-//Golang Birmingham//Events//EN")
	line("CALSCALE", "GREGORIAN")
	line("METHOD", "PUBLISH")
	line("BEGIN", "VEVENT")
	line("UID", fmt.Sprintf("%s@%s", event.ID, host))
	line("DTSTAMP", now.UTC().Format(icsTime))
	line("DTSTART", event.StartTime.UTC().Format(icsTime))
	line("DTEND", event.EndTime.UTC().Format(icsTime))
	line("SUMMARY", EscapeICSString(event.Title))
	line("DESCRIPTION", EscapeICSString(description))
	line("LOCATION", EscapeICSString(location))
	line("URL", eventURL)
	line("STATUS", "CONFIRMED")
	line("END", "VEVENT")
	line("END", "VCALENDAR")
	return []byte(b.String())
}

// EscapeICSString escapes text values per RFC 5545 section 3.3.11.
func EscapeICSString(s string) string {
	return strings.NewReplacer(
		`\`, `\\`,
		`,`, `\,`,
		`;`, `\;`,
		"\r\n", `\n`,
		"\n", `\n`,
		"\r", "",
	).Replace(s)
}

// fold splits a content line into 75-octet chunks, as RFC 5545 requires.
func fold(line string) string {
	const limit = 75
	var b strings.Builder
	n := 0
	for _, r := range line {
		size := len(string(r))
		if n+size > limit {
			b.WriteString("\r\n ")
			n = 1
		}
		b.WriteRune(r)
		n += size
	}
	b.WriteString("\r\n")
	return b.String()
}
