package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nsaltun/packman/internal/apperror"
	"github.com/nsaltun/packman/internal/model"
	"github.com/nsaltun/packman/internal/response"
	"github.com/nsaltun/packman/internal/service"
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
	registerRoutes(r *gin.Engine)
}

// httpHandler is the concrete implementation of HttpHandler
type httpHandler struct {
	packService service.PackService
}

// NewHTTPHandler creates a new HTTP handler with the given services
func NewHTTPHandler(packService service.PackService) HttpHandler {
	return &httpHandler{
		packService: packService,
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
}

// CalculatePacks handles the calculation of packs for a given quantity
func (h *httpHandler) CalculatePacks(c *gin.Context) {
	var req model.PackCalculationRequest
	// bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.BadRequestError("Invalid request format", err))
		return
	}

	// validate request
	if err := validateCalculatePacksRequest(&req); err != nil {
		_ = c.Error(apperror.ValidationError(err.Error(), err))
		return
	}

	// call service to calculate packs
	res, err := h.packService.CalculatePacks(c.Request.Context(), req.Quantity)
	if err != nil {
		_ = c.Error(err) // Pass through AppError from service/repo
		return
	}

	// return response
	response.Success(c, http.StatusOK, res)
}

// GetPackSizes handles retrieving the current pack sizes
func (h *httpHandler) GetPackSizes(c *gin.Context) {
	// call service to get pack sizes
	res, err := h.packService.GetPackSizes(c.Request.Context())
	if err != nil {
		_ = c.Error(err) // Pass through AppError from service/repo
		return
	}
	response.Success(c, http.StatusOK, res)
}

// UpdatePackSizes handles updating the pack sizes
func (h *httpHandler) UpdatePackSizes(c *gin.Context) {
	var req model.UpdatePackSizesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.BadRequestError("Invalid request format", err))
		return
	}

	// validate request
	if err := validateUpdatePackSizesRequest(&req); err != nil {
		_ = c.Error(apperror.ValidationError(err.Error(), err))
		return
	}

	// call service to update pack sizes
	err := h.packService.UpdatePackSizes(c.Request.Context(), req.PackSizes, req.UpdatedBy)
	if err != nil {
		_ = c.Error(err) // Pass through AppError from service/repo
		return
	}
	response.Success(c, http.StatusOK, gin.H{"message": "Pack sizes updated successfully"})
}
