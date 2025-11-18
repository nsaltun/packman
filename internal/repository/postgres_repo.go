package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nsaltun/packman/internal/model"
)

var (
	// ErrNotFound indicates the requested resource was not found
	ErrNotFound = errors.New("resource not found")
)

// postgresRepo implements the PackRepository interface using PostgreSQL
type postgresRepo struct {
	pool *pgxpool.Pool
}

// NewPostgresRepo creates a new PostgreSQL repository
func NewPostgresRepo(pool *pgxpool.Pool) PackRepository {
	return &postgresRepo{
		pool: pool,
	}
}

// GetPackSizes returns the current active pack sizes
func (s *postgresRepo) GetPackSizes(ctx context.Context) ([]int, error) {
	var sizes []int
	err := s.pool.QueryRow(ctx, `
		SELECT pack_sizes 
		FROM pack_configuration 
		WHERE id = 1`).Scan(&sizes)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return sizes, nil
}

// GetPackConfiguration returns the full configuration with metadata
func (s *postgresRepo) GetPackConfiguration(ctx context.Context) (*model.PackConfiguration, error) {
	var cfg model.PackConfiguration
	var updatedAt pgtype.Timestamp

	err := s.pool.QueryRow(ctx, `
		SELECT id, version, pack_sizes, updated_at, COALESCE(updated_by, '') 
		FROM pack_configuration 
		WHERE id = 1`).Scan(
		&cfg.ID,
		&cfg.Version,
		&cfg.PackSizes,
		&updatedAt,
		&cfg.UpdatedBy,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	cfg.UpdatedAt = updatedAt.Time
	return &cfg, nil
}

// UpdatePackSizes updates the pack size configuration with ACID guarantees
// Uses pessimistic locking (FOR UPDATE) to prevent lost updates caused by concurrent transactions
// Returns the updated configuration immediately after the update
func (s *postgresRepo) UpdatePackSizes(ctx context.Context, sizes []int, updatedBy string) (*model.PackConfiguration, error) {
	// Begin transaction with serializable isolation
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return nil, err
	}

	// Ensure transaction is rolled back only on error
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				slog.ErrorContext(ctx, "failed to rollback transaction",
					slog.String("error", rbErr.Error()),
				)
			}
		}
	}()

	// Lock row to prevent concurrent modifications (pessimistic locking)
	var currentVersion int
	err = tx.QueryRow(ctx, `
		SELECT version 
		FROM pack_configuration 
		WHERE id = 1 
		FOR UPDATE`).Scan(&currentVersion)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Archive current configuration before updating
	_, err = tx.Exec(ctx, `
		INSERT INTO pack_configuration_history (version, pack_sizes, created_by)
		SELECT version, pack_sizes, updated_by
		FROM pack_configuration
		WHERE id = 1`)
	if err != nil {
		return nil, err
	}

	// Update configuration and return the updated row using RETURNING clause
	var cfg model.PackConfiguration
	var updatedAt pgtype.Timestamp

	err = tx.QueryRow(ctx, `
		UPDATE pack_configuration
		SET pack_sizes = $1,
		    version = version + 1,
		    updated_at = CURRENT_TIMESTAMP,
		    updated_by = $2
		WHERE id = 1
		RETURNING id, version, pack_sizes, updated_at, updated_by`, sizes, updatedBy).Scan(
		&cfg.ID,
		&cfg.Version,
		&cfg.PackSizes,
		&updatedAt,
		&cfg.UpdatedBy,
	)
	if err != nil {
		return nil, err
	}

	cfg.UpdatedAt = updatedAt.Time

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// GetPackConfigurationHistory returns historical configurations with pagination
func (s *postgresRepo) GetPackConfigurationHistory(ctx context.Context, limit int) ([]*model.PackConfiguration, error) {

	// Validate and cap limit to prevent resource exhaustion
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// Query historical configurations ordered by creation time descending
	rows, err := s.pool.Query(ctx, `
		SELECT id, version, pack_sizes, created_at, COALESCE(created_by, '') 
		FROM pack_configuration_history 
		ORDER BY created_at DESC 
		LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []*model.PackConfiguration
	// Iterate over rows and scan into structs
	for rows.Next() {
		var cfg model.PackConfiguration
		var createdAt pgtype.Timestamp

		err := rows.Scan(
			&cfg.ID,
			&cfg.Version,
			&cfg.PackSizes,
			&createdAt,
			&cfg.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}

		cfg.UpdatedAt = createdAt.Time
		configs = append(configs, &cfg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return configs, nil
}
