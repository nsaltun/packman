package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/nsaltun/packman/config"
	"github.com/nsaltun/packman/internal/app"
	"github.com/nsaltun/packman/internal/handler"
	"github.com/nsaltun/packman/internal/repository"
	"github.com/nsaltun/packman/internal/service"
	"github.com/nsaltun/packman/migrations"
	"github.com/nsaltun/packman/pkg/postgres"
)

func main() {
	// Load config from file/env
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	//Initialize logger
	//TODO: implement log LEVEL from config
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Initialize app with lifecycle management (for graceful shutdown)
	application := app.New()

	// Postgress DB client with connection pool
	pgClient, err := postgres.NewClient(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	application.Register(pgClient)

	// DB migrations
	if err := migrations.RunMigrations(pgClient.Pool, cfg.Database.URL); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create repositories and services
	packRepo := repository.NewPostgresRepo(pgClient.Pool)
	packService := service.NewPackService(packRepo)

	// Create handlers
	packHandler := handler.NewPackHTTPHandler(packService)
	healthHandler := handler.NewHealthHandler(pgClient)

	// Create server
	server := handler.NewServer(packHandler, healthHandler, cfg.HTTP)
	application.Register(server)

	// Start all components and wait for shutdown signal
	application.Run()
}
