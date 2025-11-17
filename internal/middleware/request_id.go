package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate or use existing request ID from header
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Store in context for handlers and middleware
		c.Set("request_id", requestID)

		// Set in response header so client can reference it
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}
