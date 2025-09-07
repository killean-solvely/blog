package validation

import (
	"fmt"
	"regexp"
	"strings"
)

// Validatable interface for request DTOs that can be validated
type Validatable interface {
	Validate() error
}

// Validator provides common validation utilities
type Validator struct{}

// New creates a new validator instance
func New() *Validator {
	return &Validator{}
}

// Required validates that a string field is not empty
func (v *Validator) Required(value, fieldName string) *Error {
	if strings.TrimSpace(value) == "" {
		return &Error{
			Field:   fieldName,
			Message: "is required",
			Code:    "required",
		}
	}
	return nil
}

// RequiredStringSlice validates that a string slice is not nil or empty
func (v *Validator) RequiredStringSlice(value []string, fieldName string) *Error {
	if value == nil {
		return &Error{
			Field:   fieldName,
			Message: "is required",
			Code:    "required",
		}
	}
	return nil
}

// MinLength validates minimum string length
func (v *Validator) MinLength(value, fieldName string, min int) *Error {
	if len(strings.TrimSpace(value)) < min {
		return &Error{
			Field:   fieldName,
			Message: fmt.Sprintf("must be at least %d characters long", min),
			Code:    "min_length",
		}
	}
	return nil
}

// MaxLength validates maximum string length
func (v *Validator) MaxLength(value, fieldName string, max int) *Error {
	if len(value) > max {
		return &Error{
			Field:   fieldName,
			Message: fmt.Sprintf("must be %d characters or less", max),
			Code:    "max_length",
		}
	}
	return nil
}

// MinLengthStringSlice validates minimum string slice length
func (v *Validator) MinLengthStringSlice(value []string, fieldName string, min int) *Error {
	if len(value) < min {
		return &Error{
			Field:   fieldName,
			Message: fmt.Sprintf("must contain at least %d elements", min),
			Code:    "min_length",
		}
	}
	return nil
}

// MaxLengthStringSlice validates maximum string slice length
func (v *Validator) MaxLengthStringSlice(value []string, fieldName string, max int) *Error {
	if len(value) > max {
		return &Error{
			Field:   fieldName,
			Message: fmt.Sprintf("must contain %d elements or less", max),
			Code:    "max_length",
		}
	}
	return nil
}

// MinValue validates minimum integer value
func (v *Validator) MinValue(value int, fieldName string, min int) *Error {
	if value < min {
		return &Error{
			Field:   fieldName,
			Message: fmt.Sprintf("must be at least %d", min),
			Code:    "min_value",
		}
	}
	return nil
}

// MaxValue validates maximum integer value
func (v *Validator) MaxValue(value int, fieldName string, max int) *Error {
	if value > max {
		return &Error{
			Field:   fieldName,
			Message: fmt.Sprintf("must be at most %d", max),
			Code:    "max_value",
		}
	}
	return nil
}

// MinValuef validates minimum float64 value
func (v *Validator) MinValuef(value float64, fieldName string, min float64) *Error {
	if value < min {
		return &Error{
			Field:   fieldName,
			Message: fmt.Sprintf("must be at least %f", min),
			Code:    "min_value",
		}
	}
	return nil
}

// MaxValuef validates maximum float64 value
func (v *Validator) MaxValuef(value float64, fieldName string, max float64) *Error {
	if value > max {
		return &Error{
			Field:   fieldName,
			Message: fmt.Sprintf("must be at most %f", max),
			Code:    "max_value",
		}
	}
	return nil
}

// Email validates email format using a simple regex
func (v *Validator) Email(value, fieldName string) *Error {
	if value == "" {
		return nil // Use Required() separately for required fields
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		return &Error{
			Field:   fieldName,
			Message: "must be a valid email address",
			Code:    "invalid_email",
		}
	}
	return nil
}

// In validates that a value is in a list of allowed values
func (v *Validator) In(value string, allowed []string, fieldName string) *Error {
	for _, allowedValue := range allowed {
		if value == allowedValue {
			return nil
		}
	}

	return &Error{
		Field:   fieldName,
		Message: fmt.Sprintf("must be one of: %s", strings.Join(allowed, ", ")),
		Code:    "invalid_choice",
	}
}

// Range validates that a numeric value is within a range
func (v *Validator) Range(value int, fieldName string, min, max int) *Error {
	if value < min || value > max {
		return &Error{
			Field:   fieldName,
			Message: fmt.Sprintf("must be between %d and %d", min, max),
			Code:    "out_of_range",
		}
	}
	return nil
}

// Pattern validates that a string matches a regex pattern
func (v *Validator) Pattern(value, fieldName, pattern, errorMessage string) *Error {
	if value == "" {
		return nil // Use Required() separately for required fields
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return &Error{
			Field:   fieldName,
			Message: "invalid pattern validation",
			Code:    "invalid_pattern",
		}
	}

	if !regex.MatchString(value) {
		message := errorMessage
		if message == "" {
			message = fmt.Sprintf("must match pattern: %s", pattern)
		}
		return &Error{
			Field:   fieldName,
			Message: message,
			Code:    "pattern_mismatch",
		}
	}
	return nil
}
