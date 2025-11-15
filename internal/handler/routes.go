package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/nsaltun/packman/internal/model"
	"github.com/nsaltun/packman/internal/service"
	"github.com/nsaltun/packman/pkg/postgres"
)

// PackHttpHandler defines the interface for pack-related HTTP handlers
type PackHttpHandler interface {
	CalculatePacks(c *gin.Context)
	GetPackSizes(c *gin.Context)
	UpdatePackSizes(c *gin.Context)
}

// HttpHandler defines the interface for HTTP handlers
type HttpHandler interface {
	PackHttpHandler
	Health(c *gin.Context)
	registerRoutes(r *gin.Engine)
}

// httpHandler is the concrete implementation of HttpHandler
type httpHandler struct {
	packService service.PackService
	db          *sqlx.DB
}

// NewHTTPHandler creates a new HTTP handler with the given services
func NewHTTPHandler(packService service.PackService, db *sqlx.DB) HttpHandler {
	return &httpHandler{
		packService: packService,
		db:          db,
	}
}

// registerRoutes registers all routes for the HTTP handler
func (h *httpHandler) registerRoutes(r *gin.Engine) {
	packs := r.Group("/api/v1")
	{
		packs.POST("/calculate", h.CalculatePacks)
		packs.GET("/pack-sizes", h.GetPackSizes)
		packs.PUT("/pack-sizes", h.UpdatePackSizes)
	}
	r.GET("/health", h.Health)
}

// CalculatePacks handles the calculation of packs for a given quantity
func (h *httpHandler) CalculatePacks(c *gin.Context) {
	var req model.PackCalculationRequest
	// bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validate request
	if err := validateCalculatePacksRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// call service to calculate packs
	res, err := h.packService.CalculatePacks(c.Request.Context(), req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return response
	c.JSON(http.StatusOK, res)
}

// GetPackSizes handles retrieving the current pack sizes
func (h *httpHandler) GetPackSizes(c *gin.Context) {
	// call service to get pack sizes
	res, err := h.packService.GetPackSizes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// UpdatePackSizes handles updating the pack sizes
func (h *httpHandler) UpdatePackSizes(c *gin.Context) {
	var req model.UpdatePackSizesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validate request
	if err := validateUpdatePackSizesRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// call service to update pack sizes
	err := h.packService.UpdatePackSizes(c.Request.Context(), req.PackSizes, req.UpdatedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// Health handles the health check endpoint
func (h *httpHandler) Health(c *gin.Context) {
	dbHealth := postgres.CheckHealth(c.Request.Context(), h.db)

	statusCode := http.StatusOK
	if dbHealth.Status != "healthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, gin.H{
		"status": dbHealth.Status,
		"database": gin.H{
			"status":           dbHealth.Status,
			"response_time_ms": dbHealth.ResponseTime.Milliseconds(),
			"error":            dbHealth.Error,
		},
	})
}
