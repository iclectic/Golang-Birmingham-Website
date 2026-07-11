package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/placeholder/golang-birmingham-website/internal/config"
	"github.com/placeholder/golang-birmingham-website/internal/content"
	"github.com/placeholder/golang-birmingham-website/internal/handlers"
	"github.com/placeholder/golang-birmingham-website/internal/middleware"
)

type Server struct {
	config      *config.Config
	content     *content.Content
	httpServer  *http.Server
	handlers    *handlers.Handlers
}

func New(cfg *config.Config, content *content.Content) *Server {
	srv := &Server{
		config:  cfg,
		content: content,
	}

	// Initialize handlers
	srv.handlers = handlers.New(cfg, content)

	// Setup routes
	handler := srv.setupRoutes()

	// Create HTTP server
	srv.httpServer = &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return srv
}

func (s *Server) setupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Static files
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Health check
	mux.HandleFunc("/health", s.handlers.Health)

	// Robots.txt
	mux.HandleFunc("/robots.txt", s.handlers.Robots)

	// Sitemap
	mux.HandleFunc("/sitemap.xml", s.handlers.Sitemap)

	// Home page
	mux.HandleFunc("/", s.handlers.Home)

	// Events
	mux.HandleFunc("/events", s.handlers.Events)
	mux.HandleFunc("/events/", s.handlers.EventDetail)

	// About
	mux.HandleFunc("/about", s.handlers.About)

	// Speakers
	mux.HandleFunc("/speakers", s.handlers.Speakers)

	// Team
	mux.HandleFunc("/team", s.handlers.Team)

	// Gallery
	mux.HandleFunc("/gallery", s.handlers.Gallery)

	// Code of Conduct
	mux.HandleFunc("/code-of-conduct", s.handlers.CodeOfConduct)

	// Contact
	mux.HandleFunc("/contact", s.handlers.Contact)
	mux.HandleFunc("/contact/", s.handlers.ContactSubmit)

	// Privacy
	mux.HandleFunc("/privacy", s.handlers.Privacy)

	// Apply middleware
	var handler http.Handler = mux
	handler = middleware.SecurityHeaders(handler)
	handler = middleware.RequestLogging(handler)
	handler = middleware.RateLimiting(handler)

	return handler
}

func (s *Server) Start() error {
	log.Printf("Starting server on port %s", s.config.Port)
	log.Printf("Server URL: %s", s.config.BaseURL)

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed to start: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")

	return s.httpServer.Shutdown(ctx)
}

func (s *Server) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
