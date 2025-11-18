package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Component represents a lifecycle-aware component with start and stop capabilities
type Component interface {
	Run() error
	Close(ctx context.Context) error
}

// AbstractComponent provides default no-op implementations for Component interface
// Embed this in your components and override only the methods you need
type AbstractComponent struct{}

// Run is a no-op default implementation
func (AbstractComponent) Run() error {
	return nil
}

// Close is a no-op default implementation
func (AbstractComponent) Close(ctx context.Context) error {
	return nil
}

// App manages application lifecycle and graceful shutdown
type App struct {
	components []Component
}

// New creates a new App instance
func New() *App {
	return &App{
		components: make([]Component, 0),
	}
}

// Register adds a component to be managed by the app
func (a *App) Register(component Component) {
	a.components = append(a.components, component)
}

// Run starts all components and waits for shutdown signal
func (a *App) Run() {
	// Start all components in goroutines
	// errChan is for capturing component run errors
	errChan := make(chan error, len(a.components))
	for _, component := range a.components {
		go func(c Component) {
			if err := c.Run(); err != nil {
				errChan <- err
			}
		}(component)
	}

	// Wait for interrupt signal or component error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal or an error
	select {
	case sig := <-quit:
		slog.Info("shutdown signal received", slog.String("signal", sig.String()))
	case err := <-errChan:
		slog.Error("component error", slog.Any("error", err))
	}

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	a.Shutdown(ctx)
}

// Shutdown gracefully closes all registered components
func (a *App) Shutdown(ctx context.Context) {
	slog.Info("initiating graceful shutdown...")

	// Close components in reverse order (LIFO)
	for i := len(a.components) - 1; i >= 0; i-- {
		if err := a.components[i].Close(ctx); err != nil {
			slog.Error("error closing component", slog.Any("error", err))
		}
	}

	slog.Info("graceful shutdown completed")
}
