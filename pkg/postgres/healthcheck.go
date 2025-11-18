package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthStatus struct {
	Status       string        `json:"status"`
	ResponseTime time.Duration `json:"response_time_ms"`
	Error        string        `json:"error,omitempty"`
	PoolStats    *pgxpool.Stat `json:"pool_stats,omitempty"`
}

// CheckHealth performs a health check on the PostgreSQL database
func (c *Client) CheckHealth(ctx context.Context) HealthStatus {
	start := time.Now()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// ping the database to check connectivity
	err := c.Pool.Ping(ctx)

	elapsed := time.Since(start)

	if err != nil {
		return HealthStatus{
			Status:       "unhealthy",
			ResponseTime: elapsed,
			Error:        err.Error(),
			PoolStats:    c.Pool.Stat(),
		}
	}

	return HealthStatus{
		Status:       "healthy",
		ResponseTime: elapsed,
		PoolStats:    c.Pool.Stat(),
	}
}
