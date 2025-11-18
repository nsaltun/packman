package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nsaltun/packman/config"
	"github.com/nsaltun/packman/internal/app"
	"github.com/nsaltun/packman/internal/middleware"
)

// Server wraps the HTTP server with Gin router
type Server struct {
	app.AbstractComponent
	cfg        config.HttpConfig
	httpServer *http.Server
}

// NewServer creates and configures a new HTTP server
func NewServer(packHandler PackHTTPHandler, healthHandler HealthHandler, cfg config.HttpConfig) *Server {
	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Add middleware (order matters!)
	router.Use(middleware.RequestID())    // 1. Generate request ID first
	router.Use(gin.Logger())              // 2. Use Gin's built-in logger
	router.Use(gin.Recovery())            // 3. Recover from panics
	router.Use(middleware.ErrorHandler()) // 4. Handle errors and format responses

	// Register routes
	packHandler.registerRoutes(router)
	router.GET("/health", healthHandler.Check)

	// Configure HTTP server with timeouts
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &Server{
		cfg:        cfg,
		httpServer: httpServer,
	}
}

// Run starts the HTTP server (blocks until shutdown)
func (s *Server) Run() error {
	slog.Info("starting HTTP server", slog.String("addr", s.httpServer.Addr))
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

// Close gracefully shuts down the server
func (s *Server) Close(ctx context.Context) error {
	slog.Info("Shutting down HTTP server...")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}
	slog.Info("HTTP server stopped")
	return nil
}
