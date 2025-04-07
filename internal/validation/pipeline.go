package validation

import (
	"context"
	"net/http"
	"strings"
	"sync"
)

// Validator defines the interface for validation
type Validator interface {
	Validate(ctx context.Context, value interface{}) error
}

// ValidationRule defines a single validation rule
type ValidationRule func(ctx context.Context, value interface{}) error

// ValidationPipeline represents a chain of validation rules
type ValidationPipeline struct {
	rules []ValidationRule
	mu    sync.RWMutex
}

// NewValidationPipeline creates a new validation pipeline
func NewValidationPipeline() *ValidationPipeline {
	return &ValidationPipeline{
		rules: make([]ValidationRule, 0),
	}
}

// AddRule adds a validation rule to the pipeline
func (p *ValidationPipeline) AddRule(rule ValidationRule) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.rules = append(p.rules, rule)
}

// Validate executes all validation rules in the pipeline
func (p *ValidationPipeline) Validate(ctx context.Context, value interface{}) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var errors ValidationErrors
	for _, rule := range p.rules {
		if err := rule(ctx, value); err != nil {
			if ve, ok := err.(ValidationErrors); ok {
				errors.Errors = append(errors.Errors, ve.Errors...)
			} else {
				errors.Errors = append(errors.Errors, ValidationError{
					Code:       ErrInvalidValue,
					Message:    err.Error(),
					StatusCode: http.StatusBadRequest,
				})
			}
		}
	}

	if len(errors.Errors) > 0 {
		return errors
	}
	return nil
}

// Common validation rules
var (
	// RequiredRule checks if a value is not empty
	RequiredRule = func(field string) ValidationRule {
		return func(ctx context.Context, value interface{}) error {
			if value == nil || value == "" {
				return NewValidationErrors(ErrRequiredField(field))
			}
			return nil
		}
	}

	// EmailRule validates an email address
	EmailRule = func(field string) ValidationRule {
		return func(ctx context.Context, value interface{}) error {
			if str, ok := value.(string); ok {
				if !strings.Contains(str, "@") {
					return NewValidationErrors(ErrInvalidEmail(field))
				}
			}
			return nil
		}
	}

	// PasswordRule validates a password
	PasswordRule = func(field string, config *Config) ValidationRule {
		return func(ctx context.Context, value interface{}) error {
			if str, ok := value.(string); ok {
				if len(str) < config.Password.MinLength {
					return NewValidationErrors(ErrInvalidPassword(field))
				}
				// Add more password validation logic here
			}
			return nil
		}
	}

	// UsernameRule validates a username
	UsernameRule = func(field string, config *Config) ValidationRule {
		return func(ctx context.Context, value interface{}) error {
			if str, ok := value.(string); ok {
				if len(str) < config.Username.MinLength || len(str) > config.Username.MaxLength {
					return NewValidationErrors(ErrInvalidUsername(field))
				}
				// Add more username validation logic here
			}
			return nil
		}
	}
)
