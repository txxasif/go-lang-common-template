package validation

import (
	"fmt"
	"net/http"
	"strings"
)

// ErrorCode represents a validation error code
type ErrorCode string

const (
	// Common validation errors
	ErrRequired         ErrorCode = "required"
	ErrInvalidFormat    ErrorCode = "invalid_format"
	ErrInvalidLength    ErrorCode = "invalid_length"
	ErrInvalidValue     ErrorCode = "invalid_value"
	ErrDuplicateValue   ErrorCode = "duplicate_value"
	ErrReservedValue    ErrorCode = "reserved_value"
	ErrProfaneContent   ErrorCode = "profane_content"
	ErrInvalidChars     ErrorCode = "invalid_chars"
	ErrConsecutiveChars ErrorCode = "consecutive_chars"
	ErrInvalidBoundary  ErrorCode = "invalid_boundary"
)

// ValidationError represents a validation error with detailed information
type ValidationError struct {
	Code       ErrorCode `json:"code"`
	Field      string    `json:"field"`
	Message    string    `json:"message"`
	StatusCode int       `json:"status_code"`
	Details    any       `json:"details,omitempty"`
}

// ValidationErrors is a collection of validation errors
type ValidationErrors struct {
	Errors     []ValidationError `json:"errors"`
	StatusCode int               `json:"status_code"`
}

// Error implements the error interface
func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve.Errors {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(messages, ", ")
}

// NewValidationError creates a new validation error
func NewValidationError(code ErrorCode, field string, message string, statusCode int) ValidationError {
	return ValidationError{
		Code:       code,
		Field:      field,
		Message:    message,
		StatusCode: statusCode,
	}
}

// NewValidationErrors creates a new collection of validation errors
func NewValidationErrors(errors ...ValidationError) ValidationErrors {
	statusCode := http.StatusBadRequest
	for _, err := range errors {
		if err.StatusCode > statusCode {
			statusCode = err.StatusCode
		}
	}
	return ValidationErrors{
		Errors:     errors,
		StatusCode: statusCode,
	}
}

// Common validation errors
var (
	ErrRequiredField = func(field string) ValidationError {
		return NewValidationError(
			ErrRequired,
			field,
			fmt.Sprintf("%s is required", field),
			http.StatusBadRequest,
		)
	}

	ErrInvalidEmail = func(field string) ValidationError {
		return NewValidationError(
			ErrInvalidFormat,
			field,
			fmt.Sprintf("%s must be a valid email address", field),
			http.StatusUnprocessableEntity,
		)
	}

	ErrInvalidPassword = func(field string) ValidationError {
		return NewValidationError(
			ErrInvalidValue,
			field,
			"Password must be at least 8 characters long and contain uppercase, lowercase, number, and special character",
			http.StatusUnprocessableEntity,
		)
	}

	ErrInvalidUsername = func(field string) ValidationError {
		return NewValidationError(
			ErrInvalidValue,
			field,
			"Username must be 3-20 characters long, start with a letter, and contain only alphanumeric characters and underscores",
			http.StatusUnprocessableEntity,
		)
	}
)
