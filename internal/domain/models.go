// Package domain defines the content types rendered by the website.
package domain

import "time"

// Event is a single meetup, workshop or social.
type Event struct {
	ID              string    `yaml:"id"`
	Slug            string    `yaml:"slug"`
	Title           string    `yaml:"title"`
	Summary         string    `yaml:"summary"`
	Description     string    `yaml:"description"`
	StartTime       time.Time `yaml:"start_time"`
	EndTime         time.Time `yaml:"end_time"`
	VenueName       string    `yaml:"venue_name"`
	VenueAddress    string    `yaml:"venue_address"`
	Topics          []string  `yaml:"topics"`
	SpeakerSlugs    []string  `yaml:"speakers"`
	RegistrationURL string    `yaml:"registration_url"`
	VideoURL        string    `yaml:"video_url"`
	SlidesURL       string    `yaml:"slides_url"`

	// Speakers is resolved from SpeakerSlugs when content is loaded.
	Speakers []Speaker `yaml:"-"`
}

// IsPast reports whether the event has finished at the given time.
func (e Event) IsPast(now time.Time) bool {
	return !e.EndTime.After(now)
}

// Speaker is someone who has presented at an event.
type Speaker struct {
	Slug         string `yaml:"slug"`
	Name         string `yaml:"name"`
	Role         string `yaml:"role"`
	Organisation string `yaml:"organisation"`
	Biography    string `yaml:"biography"`
	Image        string `yaml:"image"`
	LinkedInURL  string `yaml:"linkedin_url"`
	GitHubURL    string `yaml:"github_url"`
	WebsiteURL   string `yaml:"website_url"`
}

// TeamMember is an organiser or volunteer.
type TeamMember struct {
	Name        string `yaml:"name"`
	Role        string `yaml:"role"`
	Biography   string `yaml:"biography"`
	Image       string `yaml:"image"`
	LinkedInURL string `yaml:"linkedin_url"`
	GitHubURL   string `yaml:"github_url"`
	WebsiteURL  string `yaml:"website_url"`
}

// Sponsor is an organisation supporting the meetup.
type Sponsor struct {
	Name        string `yaml:"name"`
	Category    string `yaml:"category"`
	Logo        string `yaml:"logo"`
	WebsiteURL  string `yaml:"website_url"`
	Description string `yaml:"description"`
}

// SponsorGroup is a set of sponsors sharing a category.
type SponsorGroup struct {
	Category string
	Sponsors []Sponsor
}

// SiteStats are headline numbers shown on the home page.
type SiteStats struct {
	Members  int `json:"members"`
	Events   int `json:"events"`
	Speakers int `json:"speakers"`
}

// SiteConfig holds site-wide settings and links.
type SiteConfig struct {
	Title          string       `json:"title"`
	Tagline        string       `json:"tagline"`
	Description    string       `json:"description"`
	MeetupURL      string       `json:"meetup_url"`
	SpeakerFormURL string       `json:"speaker_form_url"`
	ContactEmail   string       `json:"contact_email"`
	ConductEmail   string       `json:"conduct_email"`
	Stats          SiteStats    `json:"stats"`
	SocialLinks    []SocialLink `json:"social_links"`
}

// SocialLink is an external community profile.
type SocialLink struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Icon string `json:"icon"`
}
