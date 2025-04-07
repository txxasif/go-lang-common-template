package validation

import (
	"regexp"
	"strings"

	"myapp/internal/model"
	"myapp/internal/pkg/response"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// ValidateRegisterRequest validates a registration request
func ValidateRegisterRequest(req *model.RegisterRequest) *response.ValidationErrors {
	errors := &response.ValidationErrors{}

	// Required fields
	if strings.TrimSpace(req.Email) == "" {
		errors.Errors = append(errors.Errors, response.ValidationError{
			Field:   "email",
			Message: "Email is required",
		})
	} else if !emailRegex.MatchString(req.Email) {
		errors.Errors = append(errors.Errors, response.ValidationError{
			Field:   "email",
			Message: "Invalid email format",
		})
	}

	if strings.TrimSpace(req.Username) == "" {
		errors.Errors = append(errors.Errors, response.ValidationError{
			Field:   "username",
			Message: "Username is required",
		})
	} else if !usernameRegex.MatchString(req.Username) {
		errors.Errors = append(errors.Errors, response.ValidationError{
			Field:   "username",
			Message: "Username must be 3-20 characters long and contain only letters, numbers, and underscores",
		})
	}

	if strings.TrimSpace(req.Password) == "" {
		errors.Errors = append(errors.Errors, response.ValidationError{
			Field:   "password",
			Message: "Password is required",
		})
	} else if len(req.Password) < 8 {
		errors.Errors = append(errors.Errors, response.ValidationError{
			Field:   "password",
			Message: "Password must be at least 8 characters long",
		})
	}

	if strings.TrimSpace(req.FirstName) == "" {
		errors.Errors = append(errors.Errors, response.ValidationError{
			Field:   "firstName",
			Message: "First name is required",
		})
	} else if len(req.FirstName) < 2 {
		errors.Errors = append(errors.Errors, response.ValidationError{
			Field:   "firstName",
			Message: "First name must be at least 2 characters long",
		})
	}

	if strings.TrimSpace(req.LastName) == "" {
		errors.Errors = append(errors.Errors, response.ValidationError{
			Field:   "lastName",
			Message: "Last name is required",
		})
	} else if len(req.LastName) < 2 {
		errors.Errors = append(errors.Errors, response.ValidationError{
			Field:   "lastName",
			Message: "Last name must be at least 2 characters long",
		})
	}

	return errors
}

// ValidateLoginRequest validates a login request
func ValidateLoginRequest(req *model.LoginRequest) *response.ValidationErrors {
	errors := &response.ValidationErrors{}

	if strings.TrimSpace(req.Email) == "" {
		errors.Errors = append(errors.Errors, response.ValidationError{
			Field:   "email",
			Message: "Email is required",
		})
	} else if !emailRegex.MatchString(req.Email) {
		errors.Errors = append(errors.Errors, response.ValidationError{
			Field:   "email",
			Message: "Invalid email format",
		})
	}

	if strings.TrimSpace(req.Password) == "" {
		errors.Errors = append(errors.Errors, response.ValidationError{
			Field:   "password",
			Message: "Password is required",
		})
	}

	return errors
}
