// Package validation provides generic validation utilities for Domain-Driven Design applications.
//
// This package offers a type-safe, manual validation approach that aligns with DDD principles:
// - Explicit validation logic for better maintainability
// - Structured error reporting with field-specific messages
// - Composable validation rules that can be combined
// - Framework-agnostic design for use across different layers
//
// Key features:
// - Validatable interface for consistent validation patterns
// - Common validation utilities (required, length, email, etc.)
// - Structured error collection with field mapping
// - Support for nested validation with error prefixing
//
// Usage example:
//
//	type CreateUserRequest struct {
//		Name  string `json:"name"`
//		Email string `json:"email"`
//		Age   int    `json:"age"`
//	}
//
//	func (req CreateUserRequest) Validate() error {
//		validator := validation.New()
//		errors := validation.NewErrors()
//
//		if err := validator.Required(req.Name, "name"); err != nil {
//			errors.ValidationErrors = append(errors.ValidationErrors, *err)
//		}
//
//		if err := validator.Email(req.Email, "email"); err != nil {
//			errors.ValidationErrors = append(errors.ValidationErrors, *err)
//		}
//
//		if errors.HasErrors() {
//			return *errors
//		}
//		return nil
//	}
package validation

