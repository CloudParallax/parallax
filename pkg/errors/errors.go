package errors

import (
	"fmt"
	"net/http"
)

// AppError represents a custom application error
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Status  int    `json:"status"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewAppError creates a new application error
func NewAppError(code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// NewAppErrorWithDetails creates a new application error with details
func NewAppErrorWithDetails(code, message, details string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
		Status:  status,
	}
}

// Predefined error types
var (
	// Validation errors
	ErrValidationFailed = NewAppError("VALIDATION_FAILED", "Validation failed", http.StatusBadRequest)
	ErrInvalidInput     = NewAppError("INVALID_INPUT", "Invalid input provided", http.StatusBadRequest)
	ErrMissingField     = NewAppError("MISSING_FIELD", "Required field is missing", http.StatusBadRequest)

	// Authentication and authorization errors
	ErrUnauthorized     = NewAppError("UNAUTHORIZED", "Authentication required", http.StatusUnauthorized)
	ErrForbidden        = NewAppError("FORBIDDEN", "Access denied", http.StatusForbidden)
	ErrInvalidToken     = NewAppError("INVALID_TOKEN", "Invalid or expired token", http.StatusUnauthorized)

	// Resource errors
	ErrNotFound         = NewAppError("NOT_FOUND", "Resource not found", http.StatusNotFound)
	ErrAlreadyExists    = NewAppError("ALREADY_EXISTS", "Resource already exists", http.StatusConflict)
	ErrConflict         = NewAppError("CONFLICT", "Resource conflict", http.StatusConflict)

	// Server errors
	ErrInternalServer   = NewAppError("INTERNAL_SERVER_ERROR", "Internal server error", http.StatusInternalServerError)
	ErrServiceUnavailable = NewAppError("SERVICE_UNAVAILABLE", "Service temporarily unavailable", http.StatusServiceUnavailable)
	ErrDatabaseError    = NewAppError("DATABASE_ERROR", "Database operation failed", http.StatusInternalServerError)

	// Business logic errors
	ErrBusinessRule     = NewAppError("BUSINESS_RULE_VIOLATION", "Business rule violation", http.StatusUnprocessableEntity)
	ErrInvalidOperation = NewAppError("INVALID_OPERATION", "Invalid operation", http.StatusBadRequest)
	ErrOperationFailed  = NewAppError("OPERATION_FAILED", "Operation failed", http.StatusInternalServerError)
)

// Wrap wraps an error with additional context
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	
	if appErr, ok := err.(*AppError); ok {
		return NewAppErrorWithDetails(appErr.Code, appErr.Message, message, appErr.Status)
	}
	
	return fmt.Errorf("%s: %w", message, err)
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError extracts AppError from error
func GetAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// Error implements the error interface
func (ve *ValidationErrors) Error() string {
	if len(ve.Errors) == 0 {
		return "validation failed"
	}
	return fmt.Sprintf("validation failed: %s", ve.Errors[0].Message)
}

// Add adds a validation error
func (ve *ValidationErrors) Add(field, message string, value interface{}) {
	ve.Errors = append(ve.Errors, ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	})
}

// HasErrors returns true if there are validation errors
func (ve *ValidationErrors) HasErrors() bool {
	return len(ve.Errors) > 0
}

// NewValidationErrors creates a new ValidationErrors instance
func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Errors: make([]ValidationError, 0),
	}
}