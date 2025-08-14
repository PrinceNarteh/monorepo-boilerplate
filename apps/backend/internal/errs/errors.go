package errs

import (
	"fmt"
	"net/http"
)

// AppError represents an application error
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// Common error codes
const (
	ErrCodeValidation      = "VALIDATION_ERROR"
	ErrCodeNotFound        = "NOT_FOUND"
	ErrCodeUnauthorized    = "UNAUTHORIZED"
	ErrCodeForbidden       = "FORBIDDEN"
	ErrCodeInternal        = "INTERNAL_ERROR"
	ErrCodeBadRequest      = "BAD_REQUEST"
	ErrCodeConflict        = "CONFLICT"
	ErrCodeTooManyRequests = "TOO_MANY_REQUESTS"
)

// Predefined errors
var (
	ErrValidation      = &AppError{Code: ErrCodeValidation, Message: "Validation failed", Status: http.StatusBadRequest}
	ErrNotFound        = &AppError{Code: ErrCodeNotFound, Message: "Resource not found", Status: http.StatusNotFound}
	ErrUnauthorized    = &AppError{Code: ErrCodeUnauthorized, Message: "Unauthorized", Status: http.StatusUnauthorized}
	ErrForbidden       = &AppError{Code: ErrCodeForbidden, Message: "Forbidden", Status: http.StatusForbidden}
	ErrInternal        = &AppError{Code: ErrCodeInternal, Message: "Internal server error", Status: http.StatusInternalServerError}
	ErrBadRequest      = &AppError{Code: ErrCodeBadRequest, Message: "Bad request", Status: http.StatusBadRequest}
	ErrConflict        = &AppError{Code: ErrCodeConflict, Message: "Resource conflict", Status: http.StatusConflict}
	ErrTooManyRequests = &AppError{Code: ErrCodeTooManyRequests, Message: "Too many requests", Status: http.StatusTooManyRequests}
)

// New creates a new AppError
func New(code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// NewValidation creates a validation error with custom message
func NewValidation(message string) *AppError {
	return &AppError{
		Code:    ErrCodeValidation,
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

// NewNotFound creates a not found error with custom message
func NewNotFound(resource string) *AppError {
	return &AppError{
		Code:    ErrCodeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
		Status:  http.StatusNotFound,
	}
}

// NewInternal creates an internal error with custom message
func NewInternal(message string) *AppError {
	return &AppError{
		Code:    ErrCodeInternal,
		Message: message,
		Status:  http.StatusInternalServerError,
	}
}
