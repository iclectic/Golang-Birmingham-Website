// Package rendering parses and executes the site's HTML templates.
package rendering

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"time"
	_ "time/tzdata" // ensure Europe/London is available in minimal containers
	"unicode"
)

// London is the time zone all event times are displayed in.
var London = mustLocation("Europe/London")

// Renderer holds one pre-parsed template set per page.
type Renderer struct {
	pages map[string]*template.Template
}

// New parses layouts/*.html and components/*.html once, then clones that base
// for every page in pages/*.html. staticFS is used to fingerprint assets.
func New(templates, staticFS fs.FS) (*Renderer, error) {
	funcs := Funcs(staticFS)

	base := template.New("").Funcs(funcs)
	for _, pattern := range []string{"layouts/*.html", "components/*.html"} {
		var err error
		if base, err = base.ParseFS(templates, pattern); err != nil {
			return nil, fmt.Errorf("parsing %s: %w", pattern, err)
		}
	}

	pageFiles, err := fs.Glob(templates, "pages/*.html")
	if err != nil {
		return nil, err
	}
	r := &Renderer{pages: make(map[string]*template.Template, len(pageFiles))}
	for _, file := range pageFiles {
		t, err := base.Clone()
		if err != nil {
			return nil, err
		}
		if t, err = t.ParseFS(templates, file); err != nil {
			return nil, fmt.Errorf("parsing %s: %w", file, err)
		}
		r.pages[strings.TrimSuffix(path.Base(file), ".html")] = t
	}
	return r, nil
}

// Has reports whether a page template exists.
func (r *Renderer) Has(page string) bool {
	_, ok := r.pages[page]
	return ok
}

// Render executes a page into a buffer first so that template errors never
// produce a half-written response.
func (r *Renderer) Render(w http.ResponseWriter, status int, page string, data any) error {
	t, ok := r.pages[page]
	if !ok {
		return fmt.Errorf("unknown page %q", page)
	}
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "base.html", data); err != nil {
		return fmt.Errorf("rendering %s: %w", page, err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, err := buf.WriteTo(w)
	return err
}

// Funcs returns the template helpers available to every page.
func Funcs(staticFS fs.FS) template.FuncMap {
	versions := map[string]string{}
	return template.FuncMap{
		"asset": func(name string) string {
			v, ok := versions[name]
			if !ok {
				if b, err := fs.ReadFile(staticFS, name); err == nil {
					sum := sha256.Sum256(b)
					v = hex.EncodeToString(sum[:4])
				}
				versions[name] = v
			}
			if v == "" {
				return "/static/" + name
			}
			return "/static/" + name + "?v=" + v
		},
		"date": func(t time.Time, layout string) string {
			return t.In(London).Format(layout)
		},
		"iso": func(t time.Time) string {
			return t.Format(time.RFC3339)
		},
		"initials":   Initials,
		"paragraphs": Paragraphs,
		"join":       strings.Join,
		"hasPrefix":  strings.HasPrefix,
		"lower":      strings.ToLower,
		"mapsURL": func(q string) string {
			return "https://www.google.com/maps/search/?api=1&query=" + template.URLQueryEscaper(q)
		},
	}
}

// Initials returns up to two uppercase initials for a name.
func Initials(name string) string {
	var out []rune
	for _, word := range strings.Fields(name) {
		for _, r := range word {
			if unicode.IsLetter(r) {
				out = append(out, unicode.ToUpper(r))
				break
			}
		}
	}
	if len(out) > 2 {
		out = []rune{out[0], out[len(out)-1]}
	}
	return string(out)
}

// Paragraphs converts plain text into escaped HTML paragraphs. Blank lines
// separate paragraphs and lines starting with "- " become list items.
func Paragraphs(text string) template.HTML {
	var b strings.Builder
	inList := false
	closeList := func() {
		if inList {
			b.WriteString("</ul>")
			inList = false
		}
	}
	var para []string
	flush := func() {
		if len(para) > 0 {
			b.WriteString("<p>" + template.HTMLEscapeString(strings.Join(para, " ")) + "</p>")
			para = nil
		}
	}
	for _, raw := range strings.Split(text, "\n") {
		line := strings.TrimSpace(raw)
		switch {
		case line == "":
			flush()
			closeList()
		case strings.HasPrefix(line, "- "):
			flush()
			if !inList {
				b.WriteString("<ul>")
				inList = true
			}
			b.WriteString("<li>" + template.HTMLEscapeString(strings.TrimPrefix(line, "- ")) + "</li>")
		default:
			closeList()
			para = append(para, line)
		}
	}
	flush()
	closeList()
	return template.HTML(b.String())
}

func mustLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}
