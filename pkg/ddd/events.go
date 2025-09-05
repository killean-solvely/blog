package ddd

import "time"

// DomainEvent represents something important that happened in the domain
// All domain events must implement this interface
type DomainEvent interface {
	OccurredOn() time.Time
	EventType() string
}

// EventDispatcher handles the dispatching of domain events to registered handlers
type EventDispatcher interface {
	Dispatch(event DomainEvent) error
	Subscribe(eventType string, handler EventHandlerFunc)
}

// EventHandler processes specific types of domain events
type EventHandlerFunc func(event DomainEvent) error

// EventAggregate represents an aggregate that maintains a list of uncommitted events
// Aggregates should embed AggregateBase to implement this interface
type EventAggregate interface {
	GetUncommittedEvents() []DomainEvent // Returns uncommitted events
	MarkEventsAsCommitted()              // Clears event list after dispatch
	RecordEvent(event DomainEvent)       // Records a new event
}
