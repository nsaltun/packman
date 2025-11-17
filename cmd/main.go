package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/nsaltun/packman/config"
	"github.com/nsaltun/packman/internal/handler"
	"github.com/nsaltun/packman/internal/repository"
	"github.com/nsaltun/packman/internal/service"
	"github.com/nsaltun/packman/pkg/postgres"
)

func main() {
	//TODO: load config from file/env
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	//Initialize logger
	//TODO: implement log LEVEL from config
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Initialize postgres client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pgClient, err := postgres.NewClient(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer postgres.Close(pgClient.DB)

	slog.Info("Successfully connected to database")

	// TODO: Run migrations when we have migration files
	// if err := postgres.RunMigrations(pgClient.DB, migrations, "migrations"); err != nil {
	// 	log.Fatalf("Failed to run migrations: %v", err)
	// }

	packRepo := repository.NewPostgresRepo(pgClient.DB)
	packService := service.NewPackService(packRepo)

	// Register routes and create HTTP handler
	//TODO: make better healthcheck. Don't send it to DB client in handler
	handlr := handler.NewHTTPHandler(packService, pgClient.DB)

	// Start server and handle graceful shutdown
	server := handler.NewServer(handlr, cfg.HTTP)
	if err := server.Run(); err != nil {
		log.Fatalf("Server error: %v", err)
	}

	//TODO: Implement graceful shutdown
	//TODO: implement app interface for a common graceful shutdown
}
