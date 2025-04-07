package response

import (
	"encoding/json"
	"net/http"
)

// Response represents a standard API response
type Response struct {
	Status  int         `json:"-"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

// Error represents an API error
type Error struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors represents a collection of validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// HasErrors returns true if there are any validation errors
func (v *ValidationErrors) HasErrors() bool {
	return len(v.Errors) > 0
}

// NewSuccess creates a new success response
func NewSuccess(status int, data interface{}) *Response {
	return &Response{
		Status:  status,
		Success: true,
		Data:    data,
	}
}

// NewError creates a new error response
func NewError(status int, code, message string, details interface{}) *Response {
	return &Response{
		Status:  status,
		Success: false,
		Error: &Error{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

// NewValidationError creates a new validation error response
func NewValidationError(errors []ValidationError) *Response {
	return NewError(
		http.StatusBadRequest,
		"VALIDATION_ERROR",
		"Validation failed",
		ValidationErrors{Errors: errors},
	)
}

// NewServiceError creates a new service error response
func NewServiceError(status int, code, message string) *Response {
	return NewError(status, code, message, nil)
}

// Write writes the response to the http.ResponseWriter
func (r *Response) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Status)
	json.NewEncoder(w).Encode(r)
}
