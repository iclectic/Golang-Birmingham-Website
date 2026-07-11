package main

import (
	"fmt"
	"log"
	"os"

	"github.com/placeholder/golang-birmingham-website/internal/config"
	"github.com/placeholder/golang-birmingham-website/internal/content"
	"github.com/placeholder/golang-birmingham-website/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Load content
	loader := content.NewLoader("content")
	contentData, err := loader.Load()
	if err != nil {
		log.Fatalf("Failed to load content: %v", err)
	}

	// Validate content
	if err := validateContent(contentData); err != nil {
		log.Fatalf("Content validation failed: %v", err)
	}

	// Create and start server
	srv := server.New(cfg, contentData)

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("Server error: %v", err)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	srv.WaitForShutdown()
}

func validateContent(content *content.Content) error {
	// Basic validation to ensure required content exists
	if content.Site.Title == "" {
		return fmt.Errorf("site title is required")
	}

	if len(content.Events) == 0 {
		log.Printf("Warning: No events found in content")
	}

	if len(content.Speakers) == 0 {
		log.Printf("Warning: No speakers found in content")
	}

	if len(content.Team) == 0 {
		log.Printf("Warning: No team members found in content")
	}

	return nil
}
