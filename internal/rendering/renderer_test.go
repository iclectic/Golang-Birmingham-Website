package rendering

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderWrapsRequestedPageInBaseLayout(t *testing.T) {
	templateDir := t.TempDir()
	for name, source := range map[string]string{
		"layouts/base.html":      `<!doctype html><main>{{template "content" .}}</main>`,
		"components/header.html": `{{define "header"}}header{{end}}`,
		"pages/home.html":        `{{define "content"}}<h1>{{.Title}}</h1>{{end}}`,
		"pages/about.html":       `{{define "content"}}<h1>About</h1>{{end}}`,
	} {
		path := filepath.Join(templateDir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(source), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	renderer := NewRenderer(templateDir)
	response := httptest.NewRecorder()
	if err := renderer.Render(response, "home", struct{ Title string }{Title: "Golang Birmingham"}); err != nil {
		t.Fatal(err)
	}

	body := response.Body.String()
	if !strings.Contains(body, "<!doctype html>") {
		t.Fatalf("expected base layout, got %q", body)
	}
	if !strings.Contains(body, "<h1>Golang Birmingham</h1>") {
		t.Fatalf("expected home content, got %q", body)
	}
	if strings.Contains(body, "<h1>About</h1>") {
		t.Fatalf("rendered content from another page: %q", body)
	}
}
