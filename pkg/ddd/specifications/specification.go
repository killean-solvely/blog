package specifications

import "fmt"

// Specification represents a composable business rule that can be evaluated against an entity
type Specification[T any] interface {
	IsSatisfiedBy(entity T) bool
	String() string
}

// AndSpecification combines two specifications using a logical AND operation
type AndSpecification[T any] struct {
	left, right Specification[T]
}

// NewAndSpecification creates a new AND specification
func NewAndSpecification[T any](left, right Specification[T]) *AndSpecification[T] {
	return &AndSpecification[T]{
		left:  left,
		right: right,
	}
}

// IsSatisfiedBy returns true if both specifications are satisfied
func (s *AndSpecification[T]) IsSatisfiedBy(entity T) bool {
	return s.left.IsSatisfiedBy(entity) && s.right.IsSatisfiedBy(entity)
}

// String returns a string representation of the specification
func (s *AndSpecification[T]) String() string {
	return fmt.Sprintf("(%s AND %s)", s.left.String(), s.right.String())
}

// OrSpecification combines two specifications using a logical OR operation
type OrSpecification[T any] struct {
	left, right Specification[T]
}

// NewOrSpecification creates a new OR specification
func NewOrSpecification[T any](left, right Specification[T]) *OrSpecification[T] {
	return &OrSpecification[T]{
		left:  left,
		right: right,
	}
}

// IsSatisfiedBy returns true if either specification is satisfied
func (s *OrSpecification[T]) IsSatisfiedBy(entity T) bool {
	return s.left.IsSatisfiedBy(entity) || s.right.IsSatisfiedBy(entity)
}

// String returns a string representation of the specification
func (s *OrSpecification[T]) String() string {
	return fmt.Sprintf("(%s OR %s)", s.left.String(), s.right.String())
}

// NotSpecification negates the result of a specification
type NotSpecification[T any] struct {
	specification Specification[T]
}

// NewNotSpecification creates a new NOT specification
func NewNotSpecification[T any](specification Specification[T]) *NotSpecification[T] {
	return &NotSpecification[T]{
		specification: specification,
	}
}

// IsSatisfiedBy returns the negation of the wrapped specification
func (s *NotSpecification[T]) IsSatisfiedBy(entity T) bool {
	return !s.specification.IsSatisfiedBy(entity)
}

// String returns a string representation of the specification
func (s *NotSpecification[T]) String() string {
	return fmt.Sprintf("NOT %s", s.specification.String())
}

// And combines two specifications with AND logic
func And[T any](left, right Specification[T]) Specification[T] {
	return NewAndSpecification(left, right)
}

// Or combines two specifications with OR logic
func Or[T any](left, right Specification[T]) Specification[T] {
	return NewOrSpecification(left, right)
}

// Not negates a specification
func Not[T any](spec Specification[T]) Specification[T] {
	return NewNotSpecification(spec)
}
