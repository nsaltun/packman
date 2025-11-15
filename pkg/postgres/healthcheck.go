package postgres

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type HealthStatus struct {
	Status       string        `json:"status"`
	ResponseTime time.Duration `json:"response_time_ms"`
	Error        string        `json:"error,omitempty"`
}

// CheckHealth performs a database health check
func CheckHealth(ctx context.Context, db *sqlx.DB) HealthStatus {
	start := time.Now()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var result int
	err := db.GetContext(ctx, &result, "SELECT 1")

	elapsed := time.Since(start)

	if err != nil {
		return HealthStatus{
			Status:       "unhealthy",
			ResponseTime: elapsed,
			Error:        err.Error(),
		}
	}

	return HealthStatus{
		Status:       "healthy",
		ResponseTime: elapsed,
	}
}
