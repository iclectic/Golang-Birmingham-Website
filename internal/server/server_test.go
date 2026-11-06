package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/iclectic/golang-birmingham-website/internal/config"
	"github.com/iclectic/golang-birmingham-website/internal/content"
	"github.com/iclectic/golang-birmingham-website/web"
)

// now sits between the repository's past and upcoming events.
var now = time.Date(2026, 11, 6, 12, 0, 0, 0, time.UTC)

func newServer(t *testing.T) *httptest.Server {
	t.Helper()
	c, err := content.Load(filepath.Join("..", "..", "content"))
	if err != nil {
		t.Fatal(err)
	}
	h, err := Handler(Options{
		Config:    &config.Config{Port: "0", BaseURL: "https://golangbirmingham.test", ContentDir: "content"},
		Content:   c,
		Templates: web.Templates,
		Static:    web.Static,
		Now:       func() time.Time { return now },
	})
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return srv
}

func get(t *testing.T, srv *httptest.Server, path string) (*http.Response, string) {
	t.Helper()
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	resp, err := client.Get(srv.URL + path)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return resp, string(body)
}

func TestPages(t *testing.T) {
	srv := newServer(t)
	tests := []struct {
		path   string
		status int
		want   []string
	}{
		{"/", 200, []string{"Where Birmingham builds with", "Next meetup", "Building Microservices with Go", "Join us on Meetup", `id="contact"`}},
		{"/events", 200, []string{"Upcoming", "New Year Lightning Talks", "Past events", "Advanced Concurrency Patterns in Go"}},
		{"/events/go-microservices-workshop", 200, []string{"Upcoming event", "RSVP on Meetup", "Sarah Chen", "/events/go-microservices-workshop/calendar.ics", `"@type":"Event"`, "<li>Packaging and deploying with containers</li>"}},
		{"/events/go-concurrency-patterns", 200, []string{"Past event", "Watch the recording"}},
		{"/speaking", 200, []string{"Speak at Golang Birmingham", "Lightning talk", "Mike Johnson"}},
		{"/sponsors", 200, []string{"Venue Partners", "Birmingham Tech Hub", "Become a sponsor"}},
		{"/organisers", 200, []string{"Ibim Braide"}},
		{"/code-of-conduct", 200, []string{"Reporting a problem", "go.dev/conduct"}},
		{"/privacy", 200, []string{"doesn't use cookies"}},
		{"/does-not-exist", 404, []string{"This page has wandered off"}},
		{"/events/does-not-exist", 404, []string{"This page has wandered off"}},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			resp, body := get(t, srv, tt.path)
			if resp.StatusCode != tt.status {
				t.Fatalf("status = %d, want %d", resp.StatusCode, tt.status)
			}
			if ct := resp.Header.Get("Content-Type"); ct != "text/html; charset=utf-8" {
				t.Errorf("Content-Type = %q", ct)
			}
			for _, w := range tt.want {
				if !strings.Contains(body, w) {
					t.Errorf("body missing %q", w)
				}
			}
			if strings.Contains(body, "/gallery") {
				t.Error("page links to removed gallery")
			}
		})
	}
}

func TestSecurityHeaders(t *testing.T) {
	resp, _ := get(t, newServer(t), "/")
	for _, h := range []string{"Content-Security-Policy", "X-Content-Type-Options", "X-Frame-Options", "Referrer-Policy"} {
		if resp.Header.Get(h) == "" {
			t.Errorf("missing %s header", h)
		}
	}
}

func TestRedirects(t *testing.T) {
	srv := newServer(t)
	for path, want := range map[string]string{
		"/about":    "/",
		"/team":     "/organisers",
		"/speakers": "/speaking",
		"/contact":  "/#contact",
		"/meetup":   "https://www.meetup.com/golang-birmingham/",
		"/gallery":  "https://www.meetup.com/golang-birmingham/",
	} {
		resp, _ := get(t, srv, path)
		if resp.StatusCode/100 != 3 || resp.Header.Get("Location") != want {
			t.Errorf("%s: got %d -> %q, want redirect to %q", path, resp.StatusCode, resp.Header.Get("Location"), want)
		}
	}
}

func TestCalendarDownload(t *testing.T) {
	resp, body := get(t, newServer(t), "/events/go-microservices-workshop/calendar.ics")
	if resp.StatusCode != 200 || !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/calendar") {
		t.Fatalf("got %d %s", resp.StatusCode, resp.Header.Get("Content-Type"))
	}
	if !strings.Contains(body, "SUMMARY:Building Microservices with Go") {
		t.Errorf("unexpected ICS body:\n%s", body)
	}
}

func TestSitemapRobotsHealth(t *testing.T) {
	srv := newServer(t)
	_, sitemap := get(t, srv, "/sitemap.xml")
	for _, w := range []string{"https://golangbirmingham.test/sponsors", "https://golangbirmingham.test/events/go-web-services"} {
		if !strings.Contains(sitemap, w) {
			t.Errorf("sitemap missing %s", w)
		}
	}
	if _, robots := get(t, srv, "/robots.txt"); !strings.Contains(robots, "Sitemap: https://golangbirmingham.test/sitemap.xml") {
		t.Errorf("robots.txt: %s", robots)
	}
	if resp, body := get(t, srv, "/health"); resp.StatusCode != 200 || !strings.Contains(body, `"ok"`) {
		t.Errorf("health: %d %s", resp.StatusCode, body)
	}
}

func TestStaticAssetsAreFingerprinted(t *testing.T) {
	srv := newServer(t)
	_, body := get(t, srv, "/")
	i := strings.Index(body, "/static/css/style.css?v=")
	if i < 0 {
		t.Fatal("stylesheet is not fingerprinted")
	}
	href := body[i : i+strings.IndexByte(body[i:], '"')]
	resp, _ := get(t, srv, href)
	if resp.StatusCode != 200 || !strings.Contains(resp.Header.Get("Cache-Control"), "immutable") {
		t.Errorf("%s: %d %q", href, resp.StatusCode, resp.Header.Get("Cache-Control"))
	}
}
