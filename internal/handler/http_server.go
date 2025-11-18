package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-contrib/cors"
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

	// Configure CORS middleware
	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORS.AllowOrigins,
		AllowMethods:     cfg.CORS.AllowMethods,
		AllowHeaders:     cfg.CORS.AllowHeaders,
		ExposeHeaders:    cfg.CORS.ExposeHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
		MaxAge:           cfg.CORS.MaxAge,
	}

	// Add middleware (order matters!)
	router.Use(cors.New(corsConfig))      // 1. CORS should be first
	router.Use(middleware.RequestID())    // 2. Generate request ID
	router.Use(gin.Logger())              // 3. Use Gin's built-in logger
	router.Use(gin.Recovery())            // 4. Recover from panics
	router.Use(middleware.ErrorHandler()) // 5. Handle errors and format responses

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
