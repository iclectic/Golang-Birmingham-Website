package calendar

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/placeholder/golang-birmingham-website/internal/domain"
)

func GenerateICS(event *domain.Event, baseURL string) ([]byte, error) {
	var builder strings.Builder
	
	// ICS file header
	builder.WriteString("BEGIN:VCALENDAR\r\n")
	builder.WriteString("VERSION:2.0\r\n")
	builder.WriteString("PRODID:-//Golang Birmingham//Community Events//EN\r\n")
	builder.WriteString("CALSCALE:GREGORIAN\r\n")
	builder.WriteString("METHOD:PUBLISH\r\n")
	
	// Event details
	builder.WriteString("BEGIN:VEVENT\r\n")
	builder.WriteString(fmt.Sprintf("UID:%s@golangbirmingham.com\r\n", event.ID))
	builder.WriteString(fmt.Sprintf("DTSTART:%s\r\n", event.StartTime.UTC().Format("20060102T150405Z")))
	builder.WriteString(fmt.Sprintf("DTEND:%s\r\n", event.EndTime.UTC().Format("20060102T150405Z")))
	builder.WriteString(fmt.Sprintf("DTSTAMP:%s\r\n", time.Now().UTC().Format("20060102T150405Z")))
	
	// Title
	escapedTitle := escapeICSString(event.Title)
	builder.WriteString(fmt.Sprintf("SUMMARY:%s\r\n", escapedTitle))
	
	// Description
	description := fmt.Sprintf("%s\n\n%s", event.Summary, event.Description)
	if len(event.Speakers) > 0 {
		description += "\n\nSpeakers:"
		for _, speaker := range event.Speakers {
			description += fmt.Sprintf("\n- %s", speaker.Name)
			if speaker.Role != "" {
				description += fmt.Sprintf(" (%s)", speaker.Role)
			}
			if speaker.Organisation != "" {
				description += fmt.Sprintf(", %s", speaker.Organisation)
			}
		}
	}
	escapedDescription := escapeICSString(description)
	builder.WriteString(fmt.Sprintf("DESCRIPTION:%s\r\n", escapedDescription))
	
	// Location
	location := fmt.Sprintf("%s, %s", event.VenueName, event.VenueAddress)
	escapedLocation := escapeICSString(location)
	builder.WriteString(fmt.Sprintf("LOCATION:%s\r\n", escapedLocation))
	
	// URL
	eventURL := fmt.Sprintf("%s/events/%s", baseURL, event.Slug)
	builder.WriteString(fmt.Sprintf("URL:%s\r\n", eventURL))
	
	// Status
	builder.WriteString("STATUS:CONFIRMED\r\n")
	
	// End event
	builder.WriteString("END:VEVENT\r\n")
	
	// End calendar
	builder.WriteString("END:VCALENDAR\r\n")
	
	return []byte(builder.String()), nil
}

func ServeICS(w http.ResponseWriter, event *domain.Event, baseURL string) {
	icsData, err := GenerateICS(event, baseURL)
	if err != nil {
		http.Error(w, "Failed to generate calendar file", http.StatusInternalServerError)
		return
	}
	
	// Set headers for ICS file download
	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.ics\"", event.Slug))
	
	w.Write(icsData)
}

func escapeICSString(s string) string {
	// Escape special characters for ICS format
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, ",", "\\,")
	s = strings.ReplaceAll(s, ";", "\\;")
	s = strings.ReplaceAll(s, "\n", "\\n")
	
	// Wrap in quotes if contains spaces or special chars
	if strings.ContainsAny(s, " :") {
		s = strings.ReplaceAll(s, "\r", "")
		s = strings.ReplaceAll(s, "\n", "\\n")
	}
	
	return s
}
