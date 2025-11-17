package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nsaltun/packman/config"
)

type PostgresClient struct {
	DB *sqlx.DB
}

// NewClient creates a new PostgreSQL connection with pooling
func NewClient(ctx context.Context, cfg config.DatabaseConfig) (*PostgresClient, error) {
	db, err := sqlx.ConnectContext(ctx, "postgres", cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Verify connection
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresClient{DB: db}, nil
}

// Close closes the database connection
func Close(db *sqlx.DB) error {
	return db.Close()
}
