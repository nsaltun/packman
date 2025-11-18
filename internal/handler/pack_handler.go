package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nsaltun/packman/internal/apperror"
	"github.com/nsaltun/packman/internal/model"
	"github.com/nsaltun/packman/internal/response"
	"github.com/nsaltun/packman/internal/service"
	"github.com/nsaltun/packman/pkg/sets"
)

// PackHTTPHandler defines the interface for pack-related HTTP handlers
type PackHTTPHandler interface {
	registerRoutes(r *gin.Engine)
	CalculatePacks(c *gin.Context)
	GetPackSizes(c *gin.Context)
	UpdatePackSizes(c *gin.Context)
}

// packHTTPHandler is the concrete implementation of PackHTTPHandler
type packHTTPHandler struct {
	packService service.PackService
}

// NewPackHTTPHandler creates a new HTTP handler with the given services
func NewPackHTTPHandler(packService service.PackService) PackHTTPHandler {
	return &packHTTPHandler{
		packService: packService,
	}
}

// registerRoutes registers all routes for the HTTP handler
func (h *packHTTPHandler) registerRoutes(r *gin.Engine) {
	packs := r.Group("/api/v1")
	{
		packs.POST("/calculate", h.CalculatePacks)
		packs.GET("/pack-sizes", h.GetPackSizes)
		packs.PUT("/pack-sizes", h.UpdatePackSizes)
	}
}

// CalculatePacks handles the calculation of packs for a given quantity
func (h *packHTTPHandler) CalculatePacks(c *gin.Context) {
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
func (h *packHTTPHandler) GetPackSizes(c *gin.Context) {
	// call service to get pack sizes
	res, err := h.packService.GetPackSizes(c.Request.Context())
	if err != nil {
		_ = c.Error(err) // Pass through AppError from service/repo
		return
	}
	response.Success(c, http.StatusOK, res)
}

// UpdatePackSizes handles updating the pack sizes
func (h *packHTTPHandler) UpdatePackSizes(c *gin.Context) {
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

	//deduplicate pack sizes
	req.PackSizes = sets.DeduplicateIntSlice(req.PackSizes)

	// call service to update pack sizes
	err := h.packService.UpdatePackSizes(c.Request.Context(), req.PackSizes, req.UpdatedBy)
	if err != nil {
		_ = c.Error(err) // Pass through AppError from service/repo
		return
	}
	response.Success(c, http.StatusOK, gin.H{"message": "Pack sizes updated successfully"})
}
