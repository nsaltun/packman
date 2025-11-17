package apperror

import (
	"fmt"
	"net/http"
)

// ErrorCode represents a unique error code for categorizing errors
type ErrorCode string

const (
	// Client errors (4xx)
	ErrCodeValidation   ErrorCode = "VALIDATION_ERROR"
	ErrCodeNotFound     ErrorCode = "NOT_FOUND"
	ErrCodeConflict     ErrorCode = "CONFLICT"
	ErrCodeBadRequest   ErrorCode = "BAD_REQUEST"
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden    ErrorCode = "FORBIDDEN"

	// Server errors (5xx)
	ErrCodeInternal       ErrorCode = "INTERNAL_ERROR"
	ErrCodeServiceUnavail ErrorCode = "SERVICE_UNAVAILABLE"
)

// AppError is the base error type that includes metadata
type AppError struct {
	Code       ErrorCode              // Machine-readable error code
	Message    string                 // User-friendly message
	Internal   error                  // Internal error (for logging, not exposed to client)
	StatusCode int                    // HTTP status code
	Details    map[string]interface{} // Additional context
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %s (internal: %v)", e.Code, e.Message, e.Internal)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the internal error for error unwrapping
func (e *AppError) Unwrap() error {
	return e.Internal
}

// NewAppError creates a new application error
func NewAppError(code ErrorCode, message string, statusCode int, internal error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Internal:   internal,
		StatusCode: statusCode,
		Details:    make(map[string]interface{}),
	}
}

// WithDetails adds additional details to the error
func (e *AppError) WithDetails(key string, value interface{}) *AppError {
	e.Details[key] = value
	return e
}

// NotFoundError creates a 404 error
func NotFoundError(message string, internal error) *AppError {
	if message == "" {
		message = "Resource not found"
	}
	return NewAppError(ErrCodeNotFound, message, http.StatusNotFound, internal)
}

// ValidationError creates a 400 validation error
func ValidationError(message string, internal error) *AppError {
	if message == "" {
		message = "Validation failed"
	}
	return NewAppError(ErrCodeValidation, message, http.StatusBadRequest, internal)
}

// BadRequestError creates a 400 bad request error
func BadRequestError(message string, internal error) *AppError {
	if message == "" {
		message = "Bad request"
	}
	return NewAppError(ErrCodeBadRequest, message, http.StatusBadRequest, internal)
}

// ConflictError creates a 409 conflict error
func ConflictError(message string, internal error) *AppError {
	if message == "" {
		message = "Resource conflict"
	}
	return NewAppError(ErrCodeConflict, message, http.StatusConflict, internal)
}

// InternalError creates a 500 internal server error
func InternalError(message string, internal error) *AppError {
	if message == "" {
		message = "An internal error occurred"
	}
	return NewAppError(ErrCodeInternal, message, http.StatusInternalServerError, internal)
}

// ServiceUnavailableError creates a 503 service unavailable error
func ServiceUnavailableError(message string, internal error) *AppError {
	if message == "" {
		message = "Service temporarily unavailable"
	}
	return NewAppError(ErrCodeServiceUnavail, message, http.StatusServiceUnavailable, internal)
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// AsAppError attempts to cast an error to AppError
func AsAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}
