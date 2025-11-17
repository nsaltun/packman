package response

import (
	"github.com/gin-gonic/gin"
	"github.com/nsaltun/packman/internal/apperror"
) // APIResponse represents a standardized API response structure
type APIResponse struct {
	Data      interface{} `json:"data,omitempty"`
	Error     *APIError   `json:"error,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

// APIError represents error information in the API response
type APIError struct {
	Code    apperror.ErrorCode     `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Success sends a successful JSON response
func Success(c *gin.Context, statusCode int, data interface{}) {
	requestID, _ := c.Get("request_id")
	reqID := ""
	if requestID != nil {
		reqID = requestID.(string)
	}

	c.JSON(statusCode, APIResponse{
		Data:      data,
		RequestID: reqID,
	})
}

// Error sends an error JSON response
func Error(c *gin.Context, statusCode int, code apperror.ErrorCode, message string, details map[string]interface{}) {
	requestID, _ := c.Get("request_id")
	reqID := ""
	if requestID != nil {
		reqID = requestID.(string)
	}

	c.JSON(statusCode, APIResponse{
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
		RequestID: reqID,
	})
}

// FromAppError sends an error response from an AppError
func FromAppError(c *gin.Context, err *apperror.AppError) {
	requestID, _ := c.Get("request_id")
	reqID := ""
	if requestID != nil {
		reqID = requestID.(string)
	}

	c.JSON(err.StatusCode, APIResponse{
		Error: &APIError{
			Code:    err.Code,
			Message: err.Message,
			Details: err.Details,
		},
		RequestID: reqID,
	})
}
