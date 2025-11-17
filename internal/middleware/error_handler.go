package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nsaltun/packman/internal/apperror"
	"github.com/nsaltun/packman/internal/response"
)

// ErrorHandler handles errors and returns standardized API responses
// It logs internal error details but only exposes safe messages to clients
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there were any errors during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			requestID, _ := c.Get("request_id")

			// Check if it's an AppError (our custom error type)
			if appErr, ok := apperror.AsAppError(err); ok {
				// Log full error details (including internal error)
				if appErr.StatusCode >= 500 {
					slog.Error("server error",
						slog.String("request_id", requestID.(string)),
						slog.String("code", string(appErr.Code)),
						slog.String("message", appErr.Message),
						slog.Any("internal", appErr.Internal), // This is logged but not exposed
					)
				} else {
					slog.Warn("client error",
						slog.String("request_id", requestID.(string)),
						slog.String("code", string(appErr.Code)),
						slog.String("message", appErr.Message),
					)
				}

				// Return sanitized response (no internal details)
				response.FromAppError(c, appErr)
			} else {
				// Unknown error - log and return generic error
				slog.Error("unhandled error",
					slog.String("request_id", requestID.(string)),
					slog.String("error", err.Error()),
				)

				response.Error(c, http.StatusInternalServerError,
					apperror.ErrCodeInternal,
					"An internal error occurred",
					nil)
			}

			c.Abort()
		}
	}
}
