// Package content loads and validates the YAML/JSON content that powers the site.
package content

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/iclectic/golang-birmingham-website/internal/domain"
)

// Content is the full, validated set of site content.
type Content struct {
	Site     domain.SiteConfig
	Events   []domain.Event // sorted by start time, oldest first
	Speakers []domain.Speaker
	Team     []domain.TeamMember
	Sponsors []domain.Sponsor
}

// Load reads all content from dir and validates it.
func Load(dir string) (*Content, error) {
	c := &Content{}

	raw, err := os.ReadFile(filepath.Join(dir, "site.json"))
	if err != nil {
		return nil, fmt.Errorf("reading site config: %w", err)
	}
	if err := json.Unmarshal(raw, &c.Site); err != nil {
		return nil, fmt.Errorf("parsing site config: %w", err)
	}

	if c.Events, err = loadDir[domain.Event](filepath.Join(dir, "events")); err != nil {
		return nil, err
	}
	if c.Speakers, err = loadDir[domain.Speaker](filepath.Join(dir, "speakers")); err != nil {
		return nil, err
	}
	if c.Team, err = loadDir[domain.TeamMember](filepath.Join(dir, "team")); err != nil {
		return nil, err
	}
	if c.Sponsors, err = loadDir[domain.Sponsor](filepath.Join(dir, "sponsors")); err != nil {
		return nil, err
	}

	if err := c.resolve(); err != nil {
		return nil, err
	}
	return c, nil
}

// loadDir decodes every .yaml/.yml file in dir, in filename order.
func loadDir[T any](dir string) ([]T, error) {
	entries, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", dir, err)
	}

	var items []T
	for _, entry := range entries {
		ext := filepath.Ext(entry.Name())
		if entry.IsDir() || (ext != ".yaml" && ext != ".yml") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", path, err)
		}
		var item T
		dec := yaml.NewDecoder(strings.NewReader(string(raw)))
		dec.KnownFields(true)
		if err := dec.Decode(&item); err != nil {
			return nil, fmt.Errorf("parsing %s: %w", path, err)
		}
		items = append(items, item)
	}
	return items, nil
}

// resolve validates content and links events to speakers.
func (c *Content) resolve() error {
	var errs []error

	if c.Site.Title == "" {
		errs = append(errs, errors.New("site.json: title is required"))
	}
	if c.Site.MeetupURL == "" {
		errs = append(errs, errors.New("site.json: meetup_url is required"))
	}

	speakers := make(map[string]domain.Speaker, len(c.Speakers))
	for _, s := range c.Speakers {
		if s.Slug == "" || s.Name == "" {
			errs = append(errs, fmt.Errorf("speaker %q: slug and name are required", s.Name))
			continue
		}
		if _, dup := speakers[s.Slug]; dup {
			errs = append(errs, fmt.Errorf("speaker %q: duplicate slug", s.Slug))
		}
		speakers[s.Slug] = s
	}

	slugs := make(map[string]bool, len(c.Events))
	for i := range c.Events {
		e := &c.Events[i]
		switch {
		case e.Slug == "" || e.Title == "":
			errs = append(errs, fmt.Errorf("event %q: slug and title are required", e.Title))
			continue
		case slugs[e.Slug]:
			errs = append(errs, fmt.Errorf("event %q: duplicate slug", e.Slug))
		case e.StartTime.IsZero() || e.EndTime.IsZero():
			errs = append(errs, fmt.Errorf("event %q: start_time and end_time are required", e.Slug))
		case !e.EndTime.After(e.StartTime):
			errs = append(errs, fmt.Errorf("event %q: end_time must be after start_time", e.Slug))
		}
		slugs[e.Slug] = true
		if e.ID == "" {
			e.ID = e.Slug
		}

		e.Speakers = e.Speakers[:0]
		for _, slug := range e.SpeakerSlugs {
			s, ok := speakers[slug]
			if !ok {
				errs = append(errs, fmt.Errorf("event %q: unknown speaker %q", e.Slug, slug))
				continue
			}
			e.Speakers = append(e.Speakers, s)
		}
	}

	sort.SliceStable(c.Events, func(i, j int) bool {
		return c.Events[i].StartTime.Before(c.Events[j].StartTime)
	})

	return errors.Join(errs...)
}

// Upcoming returns events that have not yet finished, soonest first.
func (c *Content) Upcoming(now time.Time) []domain.Event {
	var out []domain.Event
	for _, e := range c.Events {
		if !e.IsPast(now) {
			out = append(out, e)
		}
	}
	return out
}

// Past returns finished events, most recent first.
func (c *Content) Past(now time.Time) []domain.Event {
	var out []domain.Event
	for i := len(c.Events) - 1; i >= 0; i-- {
		if c.Events[i].IsPast(now) {
			out = append(out, c.Events[i])
		}
	}
	return out
}

// EventBySlug finds an event by its URL slug.
func (c *Content) EventBySlug(slug string) (domain.Event, bool) {
	for _, e := range c.Events {
		if e.Slug == slug {
			return e, true
		}
	}
	return domain.Event{}, false
}

// EventsBySpeaker returns the events a speaker has presented at, most recent first.
func (c *Content) EventsBySpeaker(slug string) []domain.Event {
	var out []domain.Event
	for i := len(c.Events) - 1; i >= 0; i-- {
		for _, s := range c.Events[i].Speakers {
			if s.Slug == slug {
				out = append(out, c.Events[i])
				break
			}
		}
	}
	return out
}

// SponsorGroups groups sponsors by category, in order of first appearance.
func (c *Content) SponsorGroups() []domain.SponsorGroup {
	var groups []domain.SponsorGroup
	index := map[string]int{}
	for _, s := range c.Sponsors {
		i, ok := index[s.Category]
		if !ok {
			i = len(groups)
			index[s.Category] = i
			groups = append(groups, domain.SponsorGroup{Category: s.Category})
		}
		groups[i].Sponsors = append(groups[i].Sponsors, s)
	}
	return groups
}
