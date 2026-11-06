// Package handlers implements the website's HTTP handlers.
package handlers

import (
	"encoding/json"
	"encoding/xml"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/iclectic/golang-birmingham-website/internal/calendar"
	"github.com/iclectic/golang-birmingham-website/internal/config"
	"github.com/iclectic/golang-birmingham-website/internal/content"
	"github.com/iclectic/golang-birmingham-website/internal/domain"
	"github.com/iclectic/golang-birmingham-website/internal/rendering"
)

// Handlers serves every page on the site.
type Handlers struct {
	cfg      *config.Config
	content  *content.Content
	renderer *rendering.Renderer
	now      func() time.Time
}

// New creates the handlers. now may be nil, in which case time.Now is used.
func New(cfg *config.Config, c *content.Content, r *rendering.Renderer, now func() time.Time) *Handlers {
	if now == nil {
		now = time.Now
	}
	return &Handlers{cfg: cfg, content: c, renderer: r, now: now}
}

// Page is the data passed to every template.
type Page struct {
	Title        string
	Description  string
	Path         string
	Active       string
	CanonicalURL string
	BaseURL      string
	Year         int
	Site         domain.SiteConfig

	NextEvent     *domain.Event
	Upcoming      []domain.Event
	Past          []domain.Event
	Event         *domain.Event
	EventIsPast   bool
	Speakers      []SpeakerProfile
	Team          []domain.TeamMember
	Sponsors      []domain.Sponsor
	SponsorGroups []domain.SponsorGroup

	StructuredData template.JS
}

// SpeakerProfile is a speaker together with the events they spoke at.
type SpeakerProfile struct {
	domain.Speaker
	Events []domain.Event
}

func (h *Handlers) page(r *http.Request, active, title, description string) Page {
	return Page{
		Title:        title,
		Description:  description,
		Path:         r.URL.Path,
		Active:       active,
		CanonicalURL: h.cfg.BaseURL + r.URL.Path,
		BaseURL:      h.cfg.BaseURL,
		Year:         h.now().Year(),
		Site:         h.content.Site,
	}
}

func (h *Handlers) render(w http.ResponseWriter, r *http.Request, status int, name string, p Page) {
	if err := h.renderer.Render(w, status, name, p); err != nil {
		slog.Error("render failed", "page", name, "path", r.URL.Path, "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// Home renders the landing page.
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	now := h.now()
	p := h.page(r, "home", h.content.Site.Title+" | "+h.content.Site.Tagline, h.content.Site.Description)
	p.Upcoming = h.content.Upcoming(now)
	if len(p.Upcoming) > 0 {
		p.NextEvent = &p.Upcoming[0]
	}
	p.Past = limit(h.content.Past(now), 3)
	p.Sponsors = h.content.Sponsors
	p.StructuredData = jsonLD(map[string]any{
		"@context":    "https://schema.org",
		"@type":       "Organization",
		"name":        h.content.Site.Title,
		"url":         h.cfg.BaseURL,
		"description": h.content.Site.Description,
		"sameAs":      socialURLs(h.content.Site),
	})
	h.render(w, r, http.StatusOK, "home", p)
}

// Events lists upcoming and past events.
func (h *Handlers) Events(w http.ResponseWriter, r *http.Request) {
	now := h.now()
	p := h.page(r, "events", "Events | "+h.content.Site.Title,
		"Upcoming and past Golang Birmingham meetups, workshops and talks.")
	p.Upcoming = h.content.Upcoming(now)
	p.Past = h.content.Past(now)
	h.render(w, r, http.StatusOK, "events", p)
}

// EventDetail renders a single event.
func (h *Handlers) EventDetail(w http.ResponseWriter, r *http.Request) {
	event, ok := h.content.EventBySlug(r.PathValue("slug"))
	if !ok {
		h.NotFound(w, r)
		return
	}
	p := h.page(r, "events", event.Title+" | "+h.content.Site.Title, event.Summary)
	p.Event = &event
	p.EventIsPast = event.IsPast(h.now())
	p.StructuredData = jsonLD(eventLD(event, h.cfg.BaseURL))
	h.render(w, r, http.StatusOK, "event-detail", p)
}

// EventCalendar serves an .ics file for an event.
func (h *Handlers) EventCalendar(w http.ResponseWriter, r *http.Request) {
	event, ok := h.content.EventBySlug(r.PathValue("slug"))
	if !ok {
		h.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="`+event.Slug+`.ics"`)
	_, _ = w.Write(calendar.GenerateICS(event, h.cfg.BaseURL, h.now()))
}

// Speaking explains how to give a talk and lists past speakers.
func (h *Handlers) Speaking(w http.ResponseWriter, r *http.Request) {
	p := h.page(r, "speaking", "Speak at "+h.content.Site.Title,
		"Share what you've learned with Birmingham's Go community. First-time speakers are very welcome.")
	for _, s := range h.content.Speakers {
		p.Speakers = append(p.Speakers, SpeakerProfile{Speaker: s, Events: h.content.EventsBySpeaker(s.Slug)})
	}
	h.render(w, r, http.StatusOK, "speaking", p)
}

// Sponsors lists sponsors and explains how to support the meetup.
func (h *Handlers) Sponsors(w http.ResponseWriter, r *http.Request) {
	p := h.page(r, "sponsors", "Sponsors | "+h.content.Site.Title,
		"The organisations that keep Golang Birmingham free, and how your company can support the community.")
	p.SponsorGroups = h.content.SponsorGroups()
	h.render(w, r, http.StatusOK, "sponsors", p)
}

// Organisers lists the team.
func (h *Handlers) Organisers(w http.ResponseWriter, r *http.Request) {
	p := h.page(r, "organisers", "Organisers | "+h.content.Site.Title,
		"Meet the volunteers who run Golang Birmingham.")
	p.Team = h.content.Team
	h.render(w, r, http.StatusOK, "organisers", p)
}

// CodeOfConduct renders the code of conduct.
func (h *Handlers) CodeOfConduct(w http.ResponseWriter, r *http.Request) {
	p := h.page(r, "code-of-conduct", "Code of Conduct | "+h.content.Site.Title,
		"How we keep Golang Birmingham a safe, friendly and welcoming community for everyone.")
	h.render(w, r, http.StatusOK, "code-of-conduct", p)
}

// Privacy renders the privacy notice.
func (h *Handlers) Privacy(w http.ResponseWriter, r *http.Request) {
	p := h.page(r, "privacy", "Privacy | "+h.content.Site.Title,
		"What data the Golang Birmingham website collects, and what it doesn't.")
	h.render(w, r, http.StatusOK, "privacy", p)
}

// NotFound renders the 404 page.
func (h *Handlers) NotFound(w http.ResponseWriter, r *http.Request) {
	p := h.page(r, "", "Page not found | "+h.content.Site.Title, "The page you were looking for doesn't exist.")
	h.render(w, r, http.StatusNotFound, "404", p)
}

// Meetup redirects to the group's Meetup page.
func (h *Handlers) Meetup(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, h.content.Site.MeetupURL, http.StatusFound)
}

// Health is a liveness probe.
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Robots serves robots.txt.
func (h *Handlers) Robots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("User-agent: *\nAllow: /\n\nSitemap: " + h.cfg.BaseURL + "/sitemap.xml\n"))
}

// StaticPaths are the sitemap entries for non-event pages.
var StaticPaths = []string{"/", "/events", "/speaking", "/sponsors", "/organisers", "/code-of-conduct", "/privacy"}

// Sitemap serves sitemap.xml, including every event page.
func (h *Handlers) Sitemap(w http.ResponseWriter, r *http.Request) {
	type url struct {
		Loc     string `xml:"loc"`
		LastMod string `xml:"lastmod,omitempty"`
	}
	type urlset struct {
		XMLName xml.Name `xml:"urlset"`
		XMLNS   string   `xml:"xmlns,attr"`
		URLs    []url    `xml:"url"`
	}
	set := urlset{XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9"}
	for _, p := range StaticPaths {
		set.URLs = append(set.URLs, url{Loc: h.cfg.BaseURL + p})
	}
	for _, e := range h.content.Events {
		set.URLs = append(set.URLs, url{Loc: h.cfg.BaseURL + "/events/" + e.Slug, LastMod: e.StartTime.Format("2006-01-02")})
	}
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	_, _ = w.Write([]byte(xml.Header))
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	_ = enc.Encode(set)
}

func limit(events []domain.Event, n int) []domain.Event {
	if len(events) > n {
		return events[:n]
	}
	return events
}

func socialURLs(site domain.SiteConfig) []string {
	urls := []string{site.MeetupURL}
	for _, l := range site.SocialLinks {
		if l.URL != site.MeetupURL {
			urls = append(urls, l.URL)
		}
	}
	return urls
}

func eventLD(e domain.Event, baseURL string) map[string]any {
	ld := map[string]any{
		"@context":            "https://schema.org",
		"@type":               "Event",
		"name":                e.Title,
		"description":         e.Summary,
		"startDate":           e.StartTime.Format(time.RFC3339),
		"endDate":             e.EndTime.Format(time.RFC3339),
		"eventStatus":         "https://schema.org/EventScheduled",
		"eventAttendanceMode": "https://schema.org/OfflineEventAttendanceMode",
		"url":                 baseURL + "/events/" + e.Slug,
		"isAccessibleForFree": true,
		"organizer":           map[string]any{"@type": "Organization", "name": "Golang Birmingham", "url": baseURL},
		"location": map[string]any{
			"@type":   "Place",
			"name":    e.VenueName,
			"address": e.VenueAddress,
		},
	}
	if len(e.Speakers) > 0 {
		var performers []map[string]string
		for _, s := range e.Speakers {
			performers = append(performers, map[string]string{"@type": "Person", "name": s.Name})
		}
		ld["performer"] = performers
	}
	return ld
}

// jsonLD marshals structured data for a <script type="application/ld+json">
// block. encoding/json escapes <, > and & so the output cannot break out.
func jsonLD(v any) template.JS {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return template.JS(b)
}
