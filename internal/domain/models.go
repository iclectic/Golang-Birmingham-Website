package domain

import (
	"time"
)

type Event struct {
	ID              string    `json:"id"`
	Slug            string    `json:"slug"`
	Title           string    `json:"title"`
	Summary         string    `json:"summary"`
	Description     string    `json:"description"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	VenueName       string    `json:"venue_name"`
	VenueAddress    string    `json:"venue_address"`
	Image           string    `json:"image"`
	Topics          []string  `json:"topics"`
	Speakers        []Speaker `json:"speakers"`
	RegistrationURL string    `json:"registration_url"`
	VideoURL        string    `json:"video_url"`
	SlidesURL       string    `json:"slides_url"`
	Featured        bool      `json:"featured"`
}

type Speaker struct {
	ID           string `json:"id"`
	Slug         string `json:"slug"`
	Name         string `json:"name"`
	Role         string `json:"role"`
	Organisation string `json:"organisation"`
	Biography    string `json:"biography"`
	Image        string `json:"image"`
	LinkedInURL  string `json:"linkedin_url"`
	GitHubURL    string `json:"github_url"`
	WebsiteURL   string `json:"website_url"`
}

type TeamMember struct {
	Name        string `json:"name"`
	Role        string `json:"role"`
	Biography   string `json:"biography"`
	Image       string `json:"image"`
	LinkedInURL string `json:"linkedin_url"`
	GitHubURL   string `json:"github_url"`
	WebsiteURL  string `json:"website_url"`
}

type Sponsor struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Logo        string `json:"logo"`
	WebsiteURL  string `json:"website_url"`
	Description string `json:"description"`
}

type GalleryItem struct {
	ID        string `json:"id"`
	EventSlug string `json:"event_slug"`
	Year      int    `json:"year"`
	Category  string `json:"category"`
	Image     string `json:"image"`
	Thumbnail string `json:"thumbnail"`
	Caption   string `json:"caption"`
	AltText   string `json:"alt_text"`
}

type ContactForm struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	EnquiryType string `json:"enquiry_type"`
	Subject     string `json:"subject"`
	Message     string `json:"message"`
}

type SiteStats struct {
	Members   int `json:"members"`
	Events    int `json:"events"`
	Speakers  int `json:"speakers"`
	Talks     int `json:"talks"`
}

type SiteConfig struct {
	Title             string   `json:"title"`
	Description       string   `json:"description"`
	BaseURL           string   `json:"base_url"`
	MeetupURL         string   `json:"meetup_url"`
	LinkedInURL       string   `json:"linkedin_url"`
	YouTubeURL        string   `json:"youtube_url"`
	GitHubURL         string   `json:"github_url"`
	SpeakerFormURL    string   `json:"speaker_form_url"`
	ContactEmail      string   `json:"contact_email"`
	Stats             SiteStats `json:"stats"`
	SocialLinks       []SocialLink `json:"social_links"`
}

type SocialLink struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Icon string `json:"icon"`
}
