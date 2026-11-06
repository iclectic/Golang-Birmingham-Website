// Command web serves the Golang Birmingham website.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/iclectic/golang-birmingham-website/internal/config"
	"github.com/iclectic/golang-birmingham-website/internal/content"
	"github.com/iclectic/golang-birmingham-website/internal/server"
	"github.com/iclectic/golang-birmingham-website/web"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if cfg.IsProduction() {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	}

	c, err := content.Load(cfg.ContentDir)
	if err != nil {
		return err
	}
	slog.Info("content loaded", "events", len(c.Events), "speakers", len(c.Speakers), "team", len(c.Team), "sponsors", len(c.Sponsors))

	handler, err := server.Handler(server.Options{
		Config:    cfg,
		Content:   c,
		Templates: web.Templates,
		Static:    web.Static,
	})
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return server.Run(ctx, cfg, handler)
}
