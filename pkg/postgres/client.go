package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nsaltun/packman/config"
)

// Client wraps pgxpool with observability and best practices
type Client struct {
	Pool *pgxpool.Pool
}

// NewClient creates a production-ready pgx connection pool with observability
func NewClient(ctx context.Context, cfg config.DatabaseConfig) (*Client, error) {
	start := time.Now()

	// Build pgx pool config with best practices
	poolCfg, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	poolCfg.MaxConns = int32(cfg.MaxOpenConns)
	poolCfg.MinConns = int32(cfg.MaxIdleConns)
	poolCfg.MaxConnLifetime = cfg.ConnMaxLifetime
	poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolCfg.HealthCheckPeriod = cfg.HealthCheckPeriod

	// Add connection lifecycle hooks for observability
	poolCfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		slog.Debug("new database connection established",
			slog.String("remote_addr", conn.Config().Host))
		return nil
	}

	// Create pool
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Verify connectivity
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.Info("PostgreSQL connection pool established",
		slog.Duration("duration", time.Since(start)),
		slog.Int("max_conns", cfg.MaxOpenConns),
		slog.Int("min_conns", cfg.MaxIdleConns),
		slog.Duration("max_lifetime", cfg.ConnMaxLifetime),
	)

	return &Client{Pool: pool}, nil
}

// Close gracefully shuts down the connection pool
func (c *Client) Close() {
	slog.Info("closing database connection pool")
	c.Pool.Close()
}
