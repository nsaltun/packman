package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/nsaltun/packman/internal/model"
)

// postgresStore implements the Store interface using PostgreSQL
type postgresStore struct {
	db *sqlx.DB
}

// NewPostgresStore creates a new PostgreSQL store
func NewPostgresRepo(db *sqlx.DB) PackRepository {
	return &postgresStore{db: db}
}

// GetPackSizes returns the current active pack sizes
func (s *postgresStore) GetPackSizes(ctx context.Context) ([]int, error) {
	//TODO: Implement - query pack_configuration table for active pack sizes
	return nil, nil
}

// GetPackConfiguration returns the full configuration with metadata
func (s *postgresStore) GetPackConfiguration(ctx context.Context) (*model.PackConfiguration, error) {
	// TODO: Implement - query pack_configuration table with all fields
	return nil, nil
}

// UpdatePackSizes updates the pack size configuration
func (s *postgresStore) UpdatePackSizes(ctx context.Context, sizes []int, updatedBy string) error {
	// TODO: Implement - update pack_configuration table and insert into history
	return nil
}

// GetPackConfigurationHistory returns historical configurations
func (s *postgresStore) GetPackConfigurationHistory(ctx context.Context, limit int) ([]*model.PackConfiguration, error) {
	// TODO: Implement - query pack_configuration_history table
	return nil, nil
}
