package ddd

import (
	"fmt"
	"reflect"
	"sync"
)

// EventRegistryMetadata contains information about a registered domain event
type EventRegistryMetadata struct {
	Type        string
	Description string
	Example     interface{}
}

// eventRegistry maintains a thread-safe registry of all domain events
type eventRegistry struct {
	events map[string]EventRegistryMetadata
	mu     sync.RWMutex
}

// Global event registry instance
var EventRegistry = &eventRegistry{
	events: make(map[string]EventRegistryMetadata),
}

// GetEventRegistry returns the global event registry instance
func GetEventRegistry() *eventRegistry {
	return EventRegistry
}

// Register adds a new event type to the registry
// Panics if the event type is already registered
func (r *eventRegistry) Register(event DomainEvent, description string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	eventType := event.EventType()
	if _, exists := r.events[eventType]; exists {
		panic(fmt.Sprintf("event type %s already registered", eventType))
	}

	r.events[eventType] = EventRegistryMetadata{
		Type:        eventType,
		Description: description,
		Example:     event,
	}
}

// MustRegister is an alias for Register for consistency
func (r *eventRegistry) MustRegister(event DomainEvent, description string) {
	r.Register(event, description)
}

// IsRegistered checks if an event type is registered
func (r *eventRegistry) IsRegistered(eventType string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.events[eventType]
	return exists
}

// GetMetadata returns metadata for a registered event type
func (r *eventRegistry) GetMetadata(eventType string) (EventRegistryMetadata, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	metadata, exists := r.events[eventType]
	return metadata, exists
}

// ListEvents returns all registered event metadata
func (r *eventRegistry) ListEvents() []EventRegistryMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()

	events := make([]EventRegistryMetadata, 0, len(r.events))
	for _, metadata := range r.events {
		events = append(events, metadata)
	}
	return events
}

// ValidateEvent checks if an event is registered
func (r *eventRegistry) ValidateEvent(event DomainEvent) error {
	if !r.IsRegistered(event.EventType()) {
		return fmt.Errorf("unregistered event type: %s", event.EventType())
	}
	return nil
}

// MustValidateEvent panics if the event is not registered
func (r *eventRegistry) MustValidateEvent(event DomainEvent) {
	if err := r.ValidateEvent(event); err != nil {
		panic(err)
	}
}

// ValidatingEventDispatcher wraps an EventDispatcher to validate events before dispatch
type ValidatingEventDispatcher struct {
	dispatcher EventDispatcher
}

// NewValidatingEventDispatcher creates a new validating event dispatcher
func NewValidatingEventDispatcher(dispatcher EventDispatcher) *ValidatingEventDispatcher {
	return &ValidatingEventDispatcher{
		dispatcher: dispatcher,
	}
}

// Dispatch validates the event before dispatching it
func (v *ValidatingEventDispatcher) Dispatch(event DomainEvent) {
	if err := EventRegistry.ValidateEvent(event); err != nil {
		panic(fmt.Errorf("invalid event: %w", err))
	}
	v.dispatcher.Dispatch(event)
}

// Subscribe delegates to the wrapped dispatcher
func (v *ValidatingEventDispatcher) Subscribe(eventType string, handler EventHandlerFunc) {
	v.dispatcher.Subscribe(eventType, handler)
}

// EventRegistrar defines an interface for components that can register their events
type EventRegistrar interface {
	RegisterEvents()
}

// MustRegisterAllEvents registers events from multiple registrars
func MustRegisterAllEvents(registrars ...EventRegistrar) {
	for _, registrar := range registrars {
		registrar.RegisterEvents()
	}
}

// CheckEventIntegrity verifies that all events used by aggregates are registered
func CheckEventIntegrity(aggregates ...interface{}) error {
	var unregisteredEvents []string

	for _, aggregate := range aggregates {
		events := extractEventsFromAggregate(aggregate)
		for _, event := range events {
			if !EventRegistry.IsRegistered(event) {
				unregisteredEvents = append(unregisteredEvents, event)
			}
		}
	}

	if len(unregisteredEvents) > 0 {
		return fmt.Errorf("found unregistered events: %v", unregisteredEvents)
	}

	return nil
}

// extractEventsFromAggregate uses reflection to find potential event types in an aggregate
func extractEventsFromAggregate(aggregate interface{}) []string {
	var events []string

	v := reflect.ValueOf(aggregate)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		if method.Type.NumOut() > 0 {
			for j := 0; j < method.Type.NumOut(); j++ {
				outType := method.Type.Out(j)
				if outType.Implements(reflect.TypeOf((*DomainEvent)(nil)).Elem()) {
					events = append(events, outType.Name())
				}
			}
		}
	}

	return events
}
