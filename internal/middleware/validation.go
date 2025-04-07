package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"

	httputil "myapp/internal/pkg/http"

	"github.com/go-playground/validator/v10"
)

var (
	validate           = validator.New()
	ErrNoValidatedData = errors.New("no validated data found in context")
)

// validatedKey is a custom type for the context key
type validatedKey string

const validatedContextKey validatedKey = "validated"

// ValidateRequest validates a request body against a struct
func ValidateRequest(next http.HandlerFunc, v interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(v); err != nil {
			httputil.Error(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if err := validate.Struct(v); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				var errorMessages []string
				for _, e := range validationErrors {
					errorMessages = append(errorMessages, formatValidationError(e))
				}
				httputil.Error(w, http.StatusBadRequest, strings.Join(errorMessages, ", "))
				return
			}
			httputil.Error(w, http.StatusBadRequest, "Invalid request")
			return
		}

		next.ServeHTTP(w, r)
	}
}

// formatValidationError formats a validation error into a human-readable message
func formatValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " is required"
	case "email":
		return err.Field() + " must be a valid email"
	case "min":
		return err.Field() + " must be at least " + err.Param() + " characters long"
	case "max":
		return err.Field() + " must be at most " + err.Param() + " characters long"
	default:
		return err.Field() + " is invalid"
	}
}

// GetValidated returns the validated struct from the context
func GetValidated(r *http.Request, s interface{}) error {
	val := r.Context().Value(validatedContextKey)
	if val == nil {
		return ErrNoValidatedData
	}

	// Copy the validated data to the provided struct
	reflect.ValueOf(s).Elem().Set(reflect.ValueOf(val).Elem())
	return nil
}
