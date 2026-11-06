// Package server wires routes and middleware into an http.Server.
package server

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"github.com/iclectic/golang-birmingham-website/internal/config"
	"github.com/iclectic/golang-birmingham-website/internal/content"
	"github.com/iclectic/golang-birmingham-website/internal/handlers"
	"github.com/iclectic/golang-birmingham-website/internal/middleware"
	"github.com/iclectic/golang-birmingham-website/internal/rendering"
)

// Options configures the HTTP handler.
type Options struct {
	Config    *config.Config
	Content   *content.Content
	Templates fs.FS
	Static    fs.FS
	// Now overrides the clock, for tests.
	Now func() time.Time
}

// Handler builds the site's root http.Handler.
func Handler(opts Options) (http.Handler, error) {
	renderer, err := rendering.New(opts.Templates, opts.Static)
	if err != nil {
		return nil, err
	}
	for _, page := range []string{"home", "events", "event-detail", "speaking", "sponsors", "organisers", "code-of-conduct", "privacy", "404"} {
		if !renderer.Has(page) {
			return nil, fmt.Errorf("missing page template %q", page)
		}
	}
	h := handlers.New(opts.Config, opts.Content, renderer, opts.Now)

	mux := http.NewServeMux()
	mux.Handle("GET /static/", middleware.StaticCache(http.StripPrefix("/static/", http.FileServerFS(opts.Static))))
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /robots.txt", h.Robots)
	mux.HandleFunc("GET /sitemap.xml", h.Sitemap)

	mux.HandleFunc("GET /{$}", h.Home)
	mux.HandleFunc("GET /events", h.Events)
	mux.HandleFunc("GET /events/{slug}", h.EventDetail)
	mux.HandleFunc("GET /events/{slug}/calendar.ics", h.EventCalendar)
	mux.HandleFunc("GET /speaking", h.Speaking)
	mux.HandleFunc("GET /sponsors", h.Sponsors)
	mux.HandleFunc("GET /organisers", h.Organisers)
	mux.HandleFunc("GET /code-of-conduct", h.CodeOfConduct)
	mux.HandleFunc("GET /privacy", h.Privacy)
	mux.HandleFunc("GET /meetup", h.Meetup)

	for from, to := range map[string]string{
		"/about":    "/",
		"/contact":  "/#contact",
		"/team":     "/organisers",
		"/speakers": "/speaking",
		"/sponsor":  "/sponsors",
		"/conduct":  "/code-of-conduct",
	} {
		mux.Handle("GET "+from, http.RedirectHandler(to, http.StatusMovedPermanently))
	}
	mux.HandleFunc("GET /gallery", h.Meetup)

	mux.HandleFunc("/", h.NotFound)

	return middleware.Chain(mux,
		middleware.Recover,
		middleware.Logging,
		middleware.SecurityHeaders(opts.Config.IsProduction()),
	), nil
}

// Run serves until ctx is cancelled, then shuts down gracefully.
func Run(ctx context.Context, cfg *config.Config, handler http.Handler) error {
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errc := make(chan error, 1)
	go func() {
		slog.Info("listening", "addr", srv.Addr, "base_url", cfg.BaseURL, "env", cfg.Environment)
		errc <- srv.ListenAndServe()
	}()

	select {
	case err := <-errc:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
	}

	slog.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
