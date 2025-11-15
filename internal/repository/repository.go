package repository

import (
	"context"

	"github.com/nsaltun/packman/internal/model"
)

// PackRepository defines the interface for data access operations
type PackRepository interface {
	// GetPackSizes returns the current active pack sizes
	GetPackSizes(ctx context.Context) ([]int, error)

	// GetPackConfiguration returns the current active pack sizes
	GetPackConfiguration(ctx context.Context) (*model.PackConfiguration, error)

	// UpdatePackSizes updates the pack size configuration
	UpdatePackSizes(ctx context.Context, sizes []int, updatedBy string) error

	// GetConfigurationHistory returns historical configurations
	GetPackConfigurationHistory(ctx context.Context, limit int) ([]*model.PackConfiguration, error)
}
