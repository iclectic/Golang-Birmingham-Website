package content

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/iclectic/golang-birmingham-website/internal/domain"
)

func TestLoadRepositoryContent(t *testing.T) {
	c, err := Load(filepath.Join("..", "..", "content"))
	if err != nil {
		t.Fatalf("repository content failed validation: %v", err)
	}
	if c.Site.Title == "" || len(c.Events) == 0 || len(c.Team) == 0 {
		t.Fatalf("expected site title, events and team, got %+v", c.Site)
	}
	for i := 1; i < len(c.Events); i++ {
		if c.Events[i].StartTime.Before(c.Events[i-1].StartTime) {
			t.Errorf("events not sorted: %s before %s", c.Events[i-1].Slug, c.Events[i].Slug)
		}
	}
}

func write(t *testing.T, dir, name, body string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLoadResolvesSpeakersAndValidates(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "site.json", `{"title":"T","meetup_url":"https://meetup.example"}`)
	write(t, dir, "speakers/a.yaml", "slug: ada\nname: Ada Lovelace\n")
	write(t, dir, "events/1.yaml", "slug: one\ntitle: One\nstart_time: 2030-01-02T18:00:00Z\nend_time: 2030-01-02T20:00:00Z\nspeakers: [ada]\n")

	c, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got := c.Events[0].Speakers; len(got) != 1 || got[0].Name != "Ada Lovelace" {
		t.Fatalf("speakers not resolved: %+v", got)
	}
	if c.Events[0].ID != "one" {
		t.Errorf("ID should default to slug, got %q", c.Events[0].ID)
	}

	tests := []struct {
		name, file, body, want string
	}{
		{"unknown speaker", "events/2.yaml", "slug: two\ntitle: Two\nstart_time: 2030-01-02T18:00:00Z\nend_time: 2030-01-02T20:00:00Z\nspeakers: [bob]\n", `unknown speaker "bob"`},
		{"duplicate slug", "events/2.yaml", "slug: one\ntitle: Two\nstart_time: 2030-01-02T18:00:00Z\nend_time: 2030-01-02T20:00:00Z\n", "duplicate slug"},
		{"end before start", "events/2.yaml", "slug: two\ntitle: Two\nstart_time: 2030-01-02T18:00:00Z\nend_time: 2030-01-02T17:00:00Z\n", "end_time must be after"},
		{"unknown field", "events/2.yaml", "slug: two\ntitle: Two\nstarttime: 2030-01-02T18:00:00Z\n", "field starttime not found"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			write(t, dir, tt.file, tt.body)
			defer os.Remove(filepath.Join(dir, tt.file))
			_, err := Load(dir)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("want error containing %q, got %v", tt.want, err)
			}
		})
	}
}

func TestUpcomingAndPast(t *testing.T) {
	now := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	ev := func(slug string, start time.Time) domain.Event {
		return domain.Event{Slug: slug, StartTime: start, EndTime: start.Add(2 * time.Hour)}
	}
	c := &Content{Events: []domain.Event{
		ev("old", now.AddDate(0, -2, 0)),
		ev("recent", now.AddDate(0, -1, 0)),
		ev("happening", now.Add(-time.Hour)),
		ev("soon", now.AddDate(0, 0, 7)),
		ev("later", now.AddDate(0, 1, 0)),
	}}

	slugs := func(es []domain.Event) string {
		var s []string
		for _, e := range es {
			s = append(s, e.Slug)
		}
		return strings.Join(s, ",")
	}
	if got := slugs(c.Upcoming(now)); got != "happening,soon,later" {
		t.Errorf("Upcoming = %s", got)
	}
	if got := slugs(c.Past(now)); got != "recent,old" {
		t.Errorf("Past = %s", got)
	}
}

func TestSponsorGroupsKeepOrder(t *testing.T) {
	c := &Content{Sponsors: []domain.Sponsor{
		{Name: "A", Category: "Venue"}, {Name: "B", Category: "Food"}, {Name: "C", Category: "Venue"},
	}}
	g := c.SponsorGroups()
	if len(g) != 2 || g[0].Category != "Venue" || len(g[0].Sponsors) != 2 || g[1].Category != "Food" {
		t.Fatalf("unexpected groups: %+v", g)
	}
}
