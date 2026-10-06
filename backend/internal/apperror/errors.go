package apperror

import (
	"errors"
)

// ErrNotFound indicates a requested resource does not exist.
var ErrNotFound = errors.New("resource not found")

// ValidationError represents validation failures (both JSON parsing and business rules).
type ValidationError struct {
	Message string
	Details []FieldDetail
}

// FieldDetail describes an error for a specific field.
type FieldDetail struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError creates a ValidationError with details.
func NewValidationError(msg string, details []FieldDetail) *ValidationError {
	return &ValidationError{
		Message: msg,
		Details: details,
	}
}
