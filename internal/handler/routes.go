package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nsaltun/packman/internal/model"
	"github.com/nsaltun/packman/internal/service"
)

type PackHttpHandler interface {
	CalculatePacks(c *gin.Context)
	GetPackSizes(c *gin.Context)
	UpdatePackSizes(c *gin.Context)
}

type HttpHandler interface {
	PackHttpHandler
	Health(c *gin.Context)
	registerRoutes(r *gin.Engine)
}

type httpHandler struct {
	packService service.PackService
}

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
	r.GET("/health", h.Health)
}

func (h *httpHandler) CalculatePacks(c *gin.Context) {
	var req model.PackCalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.packService.CalculatePacks(c.Request.Context(), req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *httpHandler) GetPackSizes(c *gin.Context) {
	res, err := h.packService.GetPackSizes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *httpHandler) UpdatePackSizes(c *gin.Context) {
	var req model.UpdatePackSizesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.packService.UpdatePackSizes(c.Request.Context(), nil, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *httpHandler) Health(c *gin.Context) {
	// TODO: Implement
	// 1. Check database connection (optional)
	// 2. Return health status

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
