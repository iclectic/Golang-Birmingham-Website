package rendering

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type Renderer struct {
	templateDir string
	templates   *template.Template
}

func NewRenderer(templateDir string) *Renderer {
	// Create template functions
	funcMap := template.FuncMap{
		"now":   time.Now,
		"title": strings.Title,
	}

	var templates *template.Template

	// Parse all template files
	templateFiles := []string{}

	// Get all HTML files in subdirectories
	patterns := []string{
		filepath.Join(templateDir, "*.html"),
		filepath.Join(templateDir, "layouts", "*.html"),
		filepath.Join(templateDir, "components", "*.html"),
	}

	for _, pattern := range patterns {
		files, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		templateFiles = append(templateFiles, files...)
	}

	if len(templateFiles) == 0 {
		// Create empty template if no files found
		templates = template.New("").Funcs(funcMap)
	} else {
		templates = template.Must(template.New("").Funcs(funcMap).ParseFiles(templateFiles...))
	}

	return &Renderer{
		templateDir: templateDir,
		templates:   templates,
	}
}

func (r *Renderer) Render(w http.ResponseWriter, templateName string, data interface{}) error {
	// Create a new template for each request to avoid concurrent access issues
	tmpl, err := r.templates.Clone()
	if err != nil {
		return err
	}

	if _, err := tmpl.ParseFiles(filepath.Join(r.templateDir, "pages", templateName+".html")); err != nil {
		return err
	}

	// Execute the template
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		return err
	}

	return nil
}
