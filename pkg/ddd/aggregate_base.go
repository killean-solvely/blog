package ddd

import (
	"sync"
)

// AggregateBase provides basic aggregate functionality for domain aggregates
// It should be embedded as a pointer in aggregate structs to enable event handling
type AggregateBase struct {
	id     string        // Aggregate ID
	events []DomainEvent // Uncommitted events
	mu     sync.Mutex    // Thread safety
}

// RecordEvent adds a new domain event to the aggregate's uncommitted events
// This method is thread-safe
func (a *AggregateBase) RecordEvent(event DomainEvent) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.events = append(a.events, event)
}

// GetUncommittedEvents returns a copy of all uncommitted events
// This method is thread-safe and returns a defensive copy
func (a *AggregateBase) GetUncommittedEvents() []DomainEvent {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]DomainEvent{}, a.events...)
}

// MarkEventsAsCommitted clears all uncommitted events
// This should be called after events have been successfully persisted/dispatched
func (a *AggregateBase) MarkEventsAsCommitted() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.events = []DomainEvent{}
}

// GetID returns the aggregate's unique identifier
func (a *AggregateBase) GetID() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.id
}

// SetID sets the aggregate's unique identifier
// This is typically called during aggregate creation or reconstitution
func (a *AggregateBase) SetID(id string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.id = id
}

