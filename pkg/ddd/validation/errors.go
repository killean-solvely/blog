package validation

import (
	"fmt"
	"strings"
)

// Error represents a validation error for a specific field
type Error struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// Error implements the error interface
func (ve Error) Error() string {
	return fmt.Sprintf("%s: %s", ve.Field, ve.Message)
}

// Errors represents a collection of validation errors
type Errors struct {
	ValidationErrors []Error `json:"errors"`
}

// Error implements the error interface
func (ve Errors) Error() string {
	if len(ve.ValidationErrors) == 0 {
		return "validation failed"
	}

	var messages []string
	for _, err := range ve.ValidationErrors {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// HasErrors returns true if there are any validation errors
func (ve Errors) HasErrors() bool {
	return len(ve.ValidationErrors) > 0
}

// Add adds a validation error to the collection
func (ve *Errors) Add(field, message string) {
	ve.ValidationErrors = append(ve.ValidationErrors, Error{
		Field:   field,
		Message: message,
	})
}

// AddWithCode adds a validation error with a specific error code
func (ve *Errors) AddWithCode(field, message, code string) {
	ve.ValidationErrors = append(ve.ValidationErrors, Error{
		Field:   field,
		Message: message,
		Code:    code,
	})
}

// NewErrors creates a new Errors instance
func NewErrors() *Errors {
	return &Errors{
		ValidationErrors: make([]Error, 0),
	}
}

// NewError creates a single validation error
func NewError(field, message string) Error {
	return Error{
		Field:   field,
		Message: message,
	}
}

// NewErrorWithCode creates a validation error with a code
func NewErrorWithCode(field, message, code string) Error {
	return Error{
		Field:   field,
		Message: message,
		Code:    code,
	}
}

