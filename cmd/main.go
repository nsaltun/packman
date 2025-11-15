package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/nsaltun/packman/config"
	"github.com/nsaltun/packman/internal/handler"
	"github.com/nsaltun/packman/internal/repository"
	"github.com/nsaltun/packman/internal/service"
)

func main() {
	//TODO: implement log LEVEL from config
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	//TODO: load config from file/env
	cfg := config.NewConfig()

	//TODO: initialize postgres client
	packRepo := repository.NewPostgresRepo(nil)
	packService := service.NewPackService(packRepo)

	// Register routes and create HTTP handler
	handlr := handler.NewHTTPHandler(packService)

	// Start server and handle graceful shutdown
	server := handler.NewServer(handlr, cfg.HTTP)
	if err := server.Run(); err != nil {
		log.Fatalf("Server error: %v", err)
	}

	//TODO: Implement graceful shutdown
	//TODO: implement app interface for a common graceful shutdown
}
