package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nsaltun/packman/pkg/postgres"
)

// HealthHandler defines the interface for health check handlers
type HealthHandler interface {
	Check(c *gin.Context)
}

// healthHandler is the concrete implementation of HealthHandler
type healthHandler struct {
	pgClient *postgres.Client
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(pgClient *postgres.Client) HealthHandler {
	return &healthHandler{
		pgClient: pgClient,
	}
}

// Check handles the health check endpoint
func (h *healthHandler) Check(c *gin.Context) {
	dbHealth := h.pgClient.CheckHealth(c.Request.Context())
	statusCode := http.StatusOK
	if dbHealth.Status != "healthy" {
		statusCode = http.StatusServiceUnavailable

		// Log health check failure
		requestID, _ := c.Get("request_id")
		slog.Error("health check failed",
			slog.String("request_id", fmt.Sprintf("%v", requestID)),
			slog.String("error", dbHealth.Error),
			slog.Int64("response_time_ms", dbHealth.ResponseTime.Milliseconds()),
		)
	}

	var poolStats gin.H
	if dbHealth.PoolStats != nil {
		poolStats = gin.H{
			"total_conns":    dbHealth.PoolStats.TotalConns(),
			"acquired_conns": dbHealth.PoolStats.AcquiredConns(),
			"idle_conns":     dbHealth.PoolStats.IdleConns(),
			"max_conns":      dbHealth.PoolStats.MaxConns(),
		}
	}

	c.JSON(statusCode, gin.H{
		"status": dbHealth.Status,
		"database": gin.H{
			"status":           dbHealth.Status,
			"response_time_ms": dbHealth.ResponseTime.Milliseconds(),
			"error":            dbHealth.Error,
		},
		"connection_pool": poolStats,
	})
}
