package validation

import "strings"

// FormatValidationErrors formats validation errors into a string
func FormatValidationErrors(errs ValidationErrors) string {
	var messages []string
	for _, err := range errs.Errors {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, ", ")
}
