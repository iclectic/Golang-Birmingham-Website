package content

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/placeholder/golang-birmingham-website/internal/domain"
	"gopkg.in/yaml.v3"
)

type Loader struct {
	contentDir string
}

func NewLoader(contentDir string) *Loader {
	return &Loader{
		contentDir: contentDir,
	}
}

type Content struct {
	Site     domain.SiteConfig `json:"site"`
	Events   []domain.Event    `json:"events"`
	Speakers []domain.Speaker  `json:"speakers"`
	Team     []domain.TeamMember `json:"team"`
	Sponsors []domain.Sponsor  `json:"sponsors"`
	Gallery  []domain.GalleryItem `json:"gallery"`
}

func (l *Loader) Load() (*Content, error) {
	content := &Content{}

	// Load site configuration
	siteConfig, err := l.loadSiteConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load site config: %w", err)
	}
	content.Site = *siteConfig

	// Load events
	events, err := l.loadEvents()
	if err != nil {
		return nil, fmt.Errorf("failed to load events: %w", err)
	}
	content.Events = events

	// Load speakers
	speakers, err := l.loadSpeakers()
	if err != nil {
		return nil, fmt.Errorf("failed to load speakers: %w", err)
	}
	content.Speakers = speakers

	// Load team
	team, err := l.loadTeam()
	if err != nil {
		return nil, fmt.Errorf("failed to load team: %w", err)
	}
	content.Team = team

	// Load sponsors
	sponsors, err := l.loadSponsors()
	if err != nil {
		return nil, fmt.Errorf("failed to load sponsors: %w", err)
	}
	content.Sponsors = sponsors

	// Load gallery
	gallery, err := l.loadGallery()
	if err != nil {
		return nil, fmt.Errorf("failed to load gallery: %w", err)
	}
	content.Gallery = gallery

	return content, nil
}

func (l *Loader) loadSiteConfig() (*domain.SiteConfig, error) {
	configPath := filepath.Join(l.contentDir, "site.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config domain.SiteConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (l *Loader) loadEvents() ([]domain.Event, error) {
	eventsDir := filepath.Join(l.contentDir, "events")
	entries, err := os.ReadDir(eventsDir)
	if err != nil {
		return nil, err
	}

	var events []domain.Event
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" && filepath.Ext(entry.Name()) != ".yml" {
			continue
		}

		eventPath := filepath.Join(eventsDir, entry.Name())
		data, err := os.ReadFile(eventPath)
		if err != nil {
			return nil, err
		}

		var event domain.Event
		if err := yaml.Unmarshal(data, &event); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (l *Loader) loadSpeakers() ([]domain.Speaker, error) {
	speakersDir := filepath.Join(l.contentDir, "speakers")
	entries, err := os.ReadDir(speakersDir)
	if err != nil {
		return nil, err
	}

	var speakers []domain.Speaker
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" && filepath.Ext(entry.Name()) != ".yml" {
			continue
		}

		speakerPath := filepath.Join(speakersDir, entry.Name())
		data, err := os.ReadFile(speakerPath)
		if err != nil {
			return nil, err
		}

		var speaker domain.Speaker
		if err := yaml.Unmarshal(data, &speaker); err != nil {
			return nil, err
		}

		speakers = append(speakers, speaker)
	}

	return speakers, nil
}

func (l *Loader) loadTeam() ([]domain.TeamMember, error) {
	teamDir := filepath.Join(l.contentDir, "team")
	entries, err := os.ReadDir(teamDir)
	if err != nil {
		return nil, err
	}

	var team []domain.TeamMember
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" && filepath.Ext(entry.Name()) != ".yml" {
			continue
		}

		teamPath := filepath.Join(teamDir, entry.Name())
		data, err := os.ReadFile(teamPath)
		if err != nil {
			return nil, err
		}

		var member domain.TeamMember
		if err := yaml.Unmarshal(data, &member); err != nil {
			return nil, err
		}

		team = append(team, member)
	}

	return team, nil
}

func (l *Loader) loadSponsors() ([]domain.Sponsor, error) {
	sponsorsDir := filepath.Join(l.contentDir, "sponsors")
	entries, err := os.ReadDir(sponsorsDir)
	if err != nil {
		return nil, err
	}

	var sponsors []domain.Sponsor
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" && filepath.Ext(entry.Name()) != ".yml" {
			continue
		}

		sponsorPath := filepath.Join(sponsorsDir, entry.Name())
		data, err := os.ReadFile(sponsorPath)
		if err != nil {
			return nil, err
		}

		var sponsor domain.Sponsor
		if err := yaml.Unmarshal(data, &sponsor); err != nil {
			return nil, err
		}

		sponsors = append(sponsors, sponsor)
	}

	return sponsors, nil
}

func (l *Loader) loadGallery() ([]domain.GalleryItem, error) {
	galleryDir := filepath.Join(l.contentDir, "gallery")
	entries, err := os.ReadDir(galleryDir)
	if err != nil {
		return nil, err
	}

	var gallery []domain.GalleryItem
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" && filepath.Ext(entry.Name()) != ".yml" {
			continue
		}

		galleryPath := filepath.Join(galleryDir, entry.Name())
		data, err := os.ReadFile(galleryPath)
		if err != nil {
			return nil, err
		}

		var item domain.GalleryItem
		if err := yaml.Unmarshal(data, &item); err != nil {
			return nil, err
		}

		gallery = append(gallery, item)
	}

	return gallery, nil
}
