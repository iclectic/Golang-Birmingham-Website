package handlers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/placeholder/golang-birmingham-website/internal/calendar"
	"github.com/placeholder/golang-birmingham-website/internal/config"
	"github.com/placeholder/golang-birmingham-website/internal/content"
	"github.com/placeholder/golang-birmingham-website/internal/domain"
	"github.com/placeholder/golang-birmingham-website/internal/rendering"
)

type Handlers struct {
	config   *config.Config
	content  *content.Content
	renderer *rendering.Renderer
}

type TemplateData struct {
	Title         string
	Description   string
	CanonicalURL  string
	BaseURL       string
	ActivePage    string
	Year          int
	Site          domain.SiteConfig
	UpcomingEvent *domain.Event
	RecentEvents  []domain.Event
	Speakers      []domain.Speaker
	Team          []domain.TeamMember
	Sponsors      []domain.Sponsor
	Gallery       []domain.GalleryItem
	GalleryPreview []domain.GalleryItem
	FormErrors    map[string]string
	FormData      map[string]string
	Success       string
}

func New(cfg *config.Config, content *content.Content) *Handlers {
	renderer := rendering.NewRenderer("web/templates")

	return &Handlers{
		config:   cfg,
		content:  content,
		renderer: renderer,
	}
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.NotFound(w, r)
		return
	}

	data := h.getBaseTemplateData("home", "Golang Birmingham | Go Developer Community",
		"Join Golang Birmingham for Go programming talks, workshops, networking and community events in Birmingham, United Kingdom.")

	// Get upcoming event
	var upcomingEvents []domain.Event
	var pastEvents []domain.Event
	now := time.Now()

	for _, event := range h.content.Events {
		if event.StartTime.After(now) {
			upcomingEvents = append(upcomingEvents, event)
		} else {
			pastEvents = append(pastEvents, event)
		}
	}

	if len(upcomingEvents) > 0 {
		data.UpcomingEvent = &upcomingEvents[0]
	}

	// Get recent events (last 4)
	if len(pastEvents) > 0 {
		if len(pastEvents) > 4 {
			data.RecentEvents = pastEvents[:4]
		} else {
			data.RecentEvents = pastEvents
		}
	}

	// Get gallery preview (first 6)
	if len(h.content.Gallery) > 0 {
		if len(h.content.Gallery) > 6 {
			data.GalleryPreview = h.content.Gallery[:6]
		} else {
			data.GalleryPreview = h.content.Gallery
		}
	}

	data.Speakers = h.content.Speakers
	data.Sponsors = h.content.Sponsors

	h.renderPage(w, "home", data)
}

func (h *Handlers) Events(w http.ResponseWriter, r *http.Request) {
	data := h.getBaseTemplateData("events", "Events - Golang Birmingham",
		"Upcoming and past events for Golang Birmingham community.")

	// Filter events
	var upcomingEvents []domain.Event
	var pastEvents []domain.Event
	now := time.Now()

	for _, event := range h.content.Events {
		if event.StartTime.After(now) {
			upcomingEvents = append(upcomingEvents, event)
		} else {
			pastEvents = append(pastEvents, event)
		}
	}

	// Combine for template (upcoming first)
	var allEvents []domain.Event
	allEvents = append(allEvents, upcomingEvents...)
	allEvents = append(allEvents, pastEvents...)

	data.UpcomingEvent = nil // Will be handled differently in events template
	data.RecentEvents = allEvents

	h.renderPage(w, "events", data)
}

func (h *Handlers) EventDetail(w http.ResponseWriter, r *http.Request) {
	// Extract slug from URL
	path := strings.TrimPrefix(r.URL.Path, "/events/")
	if path == "" {
		h.Events(w, r)
		return
	}

	// Handle calendar download
	if strings.HasPrefix(path, "calendar/") {
		eventSlug := strings.TrimPrefix(path, "calendar/")
		h.EventCalendar(w, r, eventSlug)
		return
	}

	// Find event by slug
	var event *domain.Event
	for _, e := range h.content.Events {
		if e.Slug == path {
			event = &e
			break
		}
	}

	if event == nil {
		h.NotFound(w, r)
		return
	}

	data := h.getBaseTemplateData("events", event.Title+" - Golang Birmingham",
		event.Summary)

	data.UpcomingEvent = event

	h.renderPage(w, "event-detail", data)
}

func (h *Handlers) EventCalendar(w http.ResponseWriter, r *http.Request, eventSlug string) {
	// Find event by slug
	var event *domain.Event
	for _, e := range h.content.Events {
		if e.Slug == eventSlug {
			event = &e
			break
		}
	}

	if event == nil {
		h.NotFound(w, r)
		return
	}

	// Generate and serve ICS file
	calendar.ServeICS(w, event, h.config.BaseURL)
}

func (h *Handlers) About(w http.ResponseWriter, r *http.Request) {
	data := h.getBaseTemplateData("about", "About - Golang Birmingham",
		"Learn about Golang Birmingham, our mission, and what we do for the Go community in Birmingham.")

	h.renderPage(w, "about", data)
}

func (h *Handlers) Speakers(w http.ResponseWriter, r *http.Request) {
	data := h.getBaseTemplateData("speakers", "Speakers - Golang Birmingham",
		"Meet the speakers who have presented at Golang Birmingham events.")

	data.Speakers = h.content.Speakers

	h.renderPage(w, "speakers", data)
}

func (h *Handlers) Team(w http.ResponseWriter, r *http.Request) {
	data := h.getBaseTemplateData("team", "Team - Golang Birmingham",
		"Meet the organisers and volunteers behind Golang Birmingham.")

	data.Team = h.content.Team

	h.renderPage(w, "team", data)
}

func (h *Handlers) Gallery(w http.ResponseWriter, r *http.Request) {
	data := h.getBaseTemplateData("gallery", "Gallery - Golang Birmingham",
		"Photos from Golang Birmingham events and meetups.")

	data.Gallery = h.content.Gallery

	h.renderPage(w, "gallery", data)
}

func (h *Handlers) CodeOfConduct(w http.ResponseWriter, r *http.Request) {
	data := h.getBaseTemplateData("code-of-conduct", "Code of Conduct - Golang Birmingham",
		"Golang Birmingham's code of conduct for community events and interactions.")

	h.renderPage(w, "code-of-conduct", data)
}

func (h *Handlers) Contact(w http.ResponseWriter, r *http.Request) {
	data := h.getBaseTemplateData("contact", "Contact - Golang Birmingham",
		"Get in touch with the Golang Birmingham team.")

	if r.Method == "GET" {
		h.renderPage(w, "contact", data)
		return
	}

	// Handle form submission
	if r.Method == "POST" {
		h.handleContactSubmit(w, r, data)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (h *Handlers) ContactSubmit(w http.ResponseWriter, r *http.Request) {
	// This handles POST /contact/ - redirect to /contact for consistency
	http.Redirect(w, r, "/contact", http.StatusPermanentRedirect)
}

func (h *Handlers) handleContactSubmit(w http.ResponseWriter, r *http.Request, data TemplateData) {
	if err := r.ParseForm(); err != nil {
		h.ServerError(w, r, err)
		return
	}

	// Validate form
	form := domain.ContactForm{
		Name:        r.FormValue("name"),
		Email:       r.FormValue("email"),
		EnquiryType: r.FormValue("enquiry_type"),
		Subject:     r.FormValue("subject"),
		Message:     r.FormValue("message"),
	}

	errors := h.validateContactForm(form)

	if len(errors) > 0 {
		data.FormErrors = errors
		data.FormData = map[string]string{
			"name":         form.Name,
			"email":        form.Email,
			"enquiry_type": form.EnquiryType,
			"subject":      form.Subject,
			"message":      form.Message,
		}
		h.renderPage(w, "contact", data)
		return
	}

	// Log the submission (in production, this would send an email)
	log.Printf("Contact form submission: %s <%s> - %s: %s",
		form.Name, form.Email, form.Subject, form.Message)

	// Show success message
	data.Success = "Thank you for your message. We'll get back to you soon!"

	h.renderPage(w, "contact", data)
}

func (h *Handlers) validateContactForm(form domain.ContactForm) map[string]string {
	errors := make(map[string]string)

	if form.Name == "" {
		errors["name"] = "Name is required"
	}

	if form.Email == "" {
		errors["email"] = "Email is required"
	} else if !isValidEmail(form.Email) {
		errors["email"] = "Please enter a valid email address"
	}

	if form.EnquiryType == "" {
		errors["enquiry_type"] = "Please select an enquiry type"
	}

	if form.Subject == "" {
		errors["subject"] = "Subject is required"
	}

	if form.Message == "" {
		errors["message"] = "Message is required"
	} else if len(form.Message) < 10 {
		errors["message"] = "Message must be at least 10 characters long"
	}

	return errors
}

func isValidEmail(email string) bool {
	// Basic email validation
	return len(email) > 5 && len(email) < 254 &&
		   email[len(email)-4:] == ".com" || email[len(email)-3:] == ".org" ||
		   email[len(email)-4:] == ".co.uk" // Simplified for demo
}

func (h *Handlers) Privacy(w http.ResponseWriter, r *http.Request) {
	data := h.getBaseTemplateData("privacy", "Privacy - Golang Birmingham",
		"Golang Birmingham's privacy policy.")

	h.renderPage(w, "privacy", data)
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
}

func (h *Handlers) Robots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(`User-agent: *
Allow: /
Sitemap: ` + h.config.BaseURL + `/sitemap.xml`))
}

func (h *Handlers) Sitemap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")

	// Basic sitemap - in production, this would be more comprehensive
	sitemap := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>` + h.config.BaseURL + `/</loc>
    <lastmod>` + time.Now().Format("2006-01-02") + `</lastmod>
    <changefreq>weekly</changefreq>
    <priority>1.0</priority>
  </url>
  <url>
    <loc>` + h.config.BaseURL + `/events</loc>
    <lastmod>` + time.Now().Format("2006-01-02") + `</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.8</priority>
  </url>
  <url>
    <loc>` + h.config.BaseURL + `/about</loc>
    <lastmod>` + time.Now().Format("2006-01-02") + `</lastmod>
    <changefreq>monthly</changefreq>
    <priority>0.6</priority>
  </url>
  <url>
    <loc>` + h.config.BaseURL + `/speakers</loc>
    <lastmod>` + time.Now().Format("2006-01-02") + `</lastmod>
    <changefreq>monthly</changefreq>
    <priority>0.6</priority>
  </url>
  <url>
    <loc>` + h.config.BaseURL + `/team</loc>
    <lastmod>` + time.Now().Format("2006-01-02") + `</lastmod>
    <changefreq>monthly</changefreq>
    <priority>0.5</priority>
  </url>
  <url>
    <loc>` + h.config.BaseURL + `/gallery</loc>
    <lastmod>` + time.Now().Format("2006-01-02") + `</lastmod>
    <changefreq>monthly</changefreq>
    <priority>0.5</priority>
  </url>
  <url>
    <loc>` + h.config.BaseURL + `/code-of-conduct</loc>
    <lastmod>` + time.Now().Format("2006-01-02") + `</lastmod>
    <changefreq>yearly</changefreq>
    <priority>0.3</priority>
  </url>
  <url>
    <loc>` + h.config.BaseURL + `/contact</loc>
    <lastmod>` + time.Now().Format("2006-01-02") + `</lastmod>
    <changefreq>yearly</changefreq>
    <priority>0.4</priority>
  </url>
</urlset>`

	w.Write([]byte(sitemap))
}

func (h *Handlers) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	data := h.getBaseTemplateData("404", "Page Not Found - Golang Birmingham",
		"The page you're looking for doesn't exist.")

	h.renderPage(w, "404", data)
}

func (h *Handlers) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Server error: %v", err)

	w.WriteHeader(http.StatusInternalServerError)

	data := h.getBaseTemplateData("500", "Server Error - Golang Birmingham",
		"Something went wrong on our end.")

	h.renderPage(w, "500", data)
}

func (h *Handlers) getBaseTemplateData(activePage, title, description string) TemplateData {
	return TemplateData{
		Title:        title,
		Description:  description,
		CanonicalURL: h.config.BaseURL + "/" + activePage,
		BaseURL:      h.config.BaseURL,
		ActivePage:   activePage,
		Year:         time.Now().Year(),
		Site:         h.content.Site,
		Speakers:     h.content.Speakers,
		Team:         h.content.Team,
		Sponsors:     h.content.Sponsors,
		Gallery:      h.content.Gallery,
		FormErrors:   make(map[string]string),
		FormData:     make(map[string]string),
	}
}

func (h *Handlers) renderPage(w http.ResponseWriter, templateName string, data TemplateData) {
	if err := h.renderer.Render(w, templateName, data); err != nil {
		log.Printf("Template rendering error: %v", err)
		h.ServerError(w, nil, err)
	}
}
