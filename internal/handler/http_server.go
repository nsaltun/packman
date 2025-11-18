package handler

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nsaltun/packman/config"
	"github.com/nsaltun/packman/internal/middleware"
)

// Server wraps the HTTP server with Gin router
type Server struct {
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

// Run starts the server and handles graceful shutdown
func (s *Server) Run() error {
	// Start server in goroutine
	go func() {
		slog.Info("Starting server on ", slog.String("addr", s.httpServer.Addr))
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.Shutdown(ctx)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	slog.Info("Shutting down server...")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}
	slog.Info("Server exited gracefully")
	return nil
}
