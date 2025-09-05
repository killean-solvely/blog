package testing

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"blog/pkg/ddd"
)

// EventRecorder records all dispatched events for testing purposes
type EventRecorder struct {
	events []ddd.DomainEvent
	mu     sync.RWMutex
}

// NewEventRecorder creates a new event recorder
func NewEventRecorder() *EventRecorder {
	return &EventRecorder{
		events: make([]ddd.DomainEvent, 0),
	}
}

// Dispatch records the event
func (er *EventRecorder) Dispatch(event ddd.DomainEvent) {
	er.mu.Lock()
	defer er.mu.Unlock()
	er.events = append(er.events, event)
}

// DispatchMultiple is a helper for dispatching multiple events
func (er *EventRecorder) DispatchMultiple(events ...ddd.DomainEvent) {
	for _, event := range events {
		er.Dispatch(event)
	}
}

// Subscribe implements EventDispatcher interface (no-op for recorder)
func (er *EventRecorder) Subscribe(eventType string, handler ddd.EventHandlerFunc) {
	// Not implemented for recorder
}

// GetEvents returns a copy of all recorded events
func (er *EventRecorder) GetEvents() []ddd.DomainEvent {
	er.mu.RLock()
	defer er.mu.RUnlock()
	return append([]ddd.DomainEvent{}, er.events...)
}

// GetEventsByType returns all events of a specific type
func (er *EventRecorder) GetEventsByType(eventType string) []ddd.DomainEvent {
	er.mu.RLock()
	defer er.mu.RUnlock()

	var filtered []ddd.DomainEvent
	for _, event := range er.events {
		if event.EventType() == eventType {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

// Clear removes all recorded events
func (er *EventRecorder) Clear() {
	er.mu.Lock()
	defer er.mu.Unlock()
	er.events = []ddd.DomainEvent{}
}

// Count returns the total number of recorded events
func (er *EventRecorder) Count() int {
	er.mu.RLock()
	defer er.mu.RUnlock()
	return len(er.events)
}

// EventAssertion provides fluent assertions for recorded events
type EventAssertion struct {
	t        *testing.T
	recorder *EventRecorder
}

// AssertEvents creates a new event assertion helper
func AssertEvents(t *testing.T, recorder *EventRecorder) *EventAssertion {
	return &EventAssertion{
		t:        t,
		recorder: recorder,
	}
}

// HasEventType asserts that at least one event of the given type was recorded
func (ea *EventAssertion) HasEventType(eventType string) *EventAssertion {
	events := ea.recorder.GetEventsByType(eventType)
	if len(events) == 0 {
		ea.t.Errorf("expected event of type %s, but none found", eventType)
	}
	return ea
}

// HasEventCount asserts that exactly the expected number of events of the given type were recorded
func (ea *EventAssertion) HasEventCount(eventType string, expectedCount int) *EventAssertion {
	events := ea.recorder.GetEventsByType(eventType)
	if len(events) != expectedCount {
		ea.t.Errorf(
			"expected %d events of type %s, but found %d",
			expectedCount,
			eventType,
			len(events),
		)
	}
	return ea
}

// HasTotalEventCount asserts the total number of recorded events
func (ea *EventAssertion) HasTotalEventCount(expectedCount int) *EventAssertion {
	actualCount := ea.recorder.Count()
	if actualCount != expectedCount {
		ea.t.Errorf("expected %d total events, but found %d", expectedCount, actualCount)
	}
	return ea
}

// EventMatches asserts that at least one event of the given type matches the provided matcher
func (ea *EventAssertion) EventMatches(
	eventType string,
	matcher func(event ddd.DomainEvent) bool,
) *EventAssertion {
	events := ea.recorder.GetEventsByType(eventType)
	for _, event := range events {
		if matcher(event) {
			return ea
		}
	}
	ea.t.Errorf("no event of type %s matched the provided criteria", eventType)
	return ea
}

// TestEventDispatcher is a test implementation of EventDispatcher with failure simulation
type TestEventDispatcher struct {
	handlers      map[string][]ddd.EventHandlerFunc
	failOnEvent   map[string]error
	delayOnEvent  map[string]time.Duration
	dispatchCount map[string]int
	mu            sync.RWMutex
}

// NewTestEventDispatcher creates a new test event dispatcher
func NewTestEventDispatcher() *TestEventDispatcher {
	return &TestEventDispatcher{
		handlers:      make(map[string][]ddd.EventHandlerFunc),
		failOnEvent:   make(map[string]error),
		delayOnEvent:  make(map[string]time.Duration),
		dispatchCount: make(map[string]int),
	}
}

// RegisterHandler registers a handler for a specific event type
func (ted *TestEventDispatcher) RegisterHandler(eventType string, handler ddd.EventHandlerFunc) {
	ted.mu.Lock()
	defer ted.mu.Unlock()
	ted.handlers[eventType] = append(ted.handlers[eventType], handler)
}

// SimulateFailure configures the dispatcher to fail when handling a specific event type
func (ted *TestEventDispatcher) SimulateFailure(eventType string, err error) {
	ted.mu.Lock()
	defer ted.mu.Unlock()
	ted.failOnEvent[eventType] = err
}

// SimulateDelay configures the dispatcher to delay when handling a specific event type
func (ted *TestEventDispatcher) SimulateDelay(eventType string, delay time.Duration) {
	ted.mu.Lock()
	defer ted.mu.Unlock()
	ted.delayOnEvent[eventType] = delay
}

// Dispatch handles the event, applying any configured failures or delays
func (ted *TestEventDispatcher) Dispatch(event ddd.DomainEvent) {
	ted.mu.Lock()
	ted.dispatchCount[event.EventType()]++
	ted.mu.Unlock()

	ted.mu.RLock()
	if err, exists := ted.failOnEvent[event.EventType()]; exists {
		ted.mu.RUnlock()
		panic(err) // Simulate error in test
	}

	if delay, exists := ted.delayOnEvent[event.EventType()]; exists {
		ted.mu.RUnlock()
		time.Sleep(delay)
	} else {
		ted.mu.RUnlock()
	}

	ted.mu.RLock()
	handlers := ted.handlers[event.EventType()]
	ted.mu.RUnlock()

	for _, handler := range handlers {
		if err := handler(event); err != nil {
			panic(err) // Simulate error in test
		}
	}
}

// Subscribe implements EventDispatcher interface
func (ted *TestEventDispatcher) Subscribe(eventType string, handler ddd.EventHandlerFunc) {
	ted.RegisterHandler(eventType, handler)
}

// GetDispatchCount returns the number of times events of the given type have been dispatched
func (ted *TestEventDispatcher) GetDispatchCount(eventType string) int {
	ted.mu.RLock()
	defer ted.mu.RUnlock()
	return ted.dispatchCount[eventType]
}

// EventScenario represents a test scenario for event-driven behavior
type EventScenario struct {
	Name         string
	Given        func() []ddd.DomainEvent
	When         func() error
	Then         func(t *testing.T, recorder *EventRecorder)
	ShouldError  bool
	ErrorMessage string
}

// RunEventScenarios executes a series of event scenarios
func RunEventScenarios(t *testing.T, scenarios []EventScenario) {
	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			recorder := NewEventRecorder()

			if scenario.Given != nil {
				initialEvents := scenario.Given()
				recorder.DispatchMultiple(initialEvents...)
			}

			err := scenario.When()

			if scenario.ShouldError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if scenario.ErrorMessage != "" && err.Error() != scenario.ErrorMessage {
					t.Errorf("expected error message '%s', got '%s'", scenario.ErrorMessage, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if scenario.Then != nil {
				scenario.Then(t, recorder)
			}
		})
	}
}

// AggregateTestHelper provides testing utilities for domain aggregates
type AggregateTestHelper struct {
	t         *testing.T
	aggregate interface{}
}

// NewAggregateTestHelper creates a new aggregate test helper
func NewAggregateTestHelper(t *testing.T, aggregate interface{}) *AggregateTestHelper {
	return &AggregateTestHelper{
		t:         t,
		aggregate: aggregate,
	}
}

// AssertEventRaised asserts that the aggregate has raised an event of the expected type
func (ath *AggregateTestHelper) AssertEventRaised(expectedEvent ddd.DomainEvent) {
	if eventAggregate, ok := ath.aggregate.(interface{ GetUncommittedEvents() []ddd.DomainEvent }); ok {
		events := eventAggregate.GetUncommittedEvents()
		for _, event := range events {
			if reflect.TypeOf(event) == reflect.TypeOf(expectedEvent) {
				return
			}
		}
		ath.t.Errorf("expected event of type %T to be raised, but it was not", expectedEvent)
	} else {
		ath.t.Errorf("aggregate does not implement GetUncommittedEvents")
	}
}

// AssertNoEvents asserts that the aggregate has no uncommitted events
func (ath *AggregateTestHelper) AssertNoEvents() {
	if eventAggregate, ok := ath.aggregate.(interface{ GetUncommittedEvents() []ddd.DomainEvent }); ok {
		events := eventAggregate.GetUncommittedEvents()
		if len(events) > 0 {
			ath.t.Errorf("expected no events, but found %d", len(events))
		}
	} else {
		ath.t.Errorf("aggregate does not implement GetUncommittedEvents")
	}
}

// MockEvent is a simple implementation of DomainEvent for testing
type MockEvent struct {
	occurredOn  time.Time
	Type        string
	AggregateId string
	Data        map[string]interface{}
}

// NewMockEvent creates a new mock event
func NewMockEvent(eventType, aggregateID string, data map[string]interface{}) *MockEvent {
	return &MockEvent{
		occurredOn:  time.Now(),
		Type:        eventType,
		AggregateId: aggregateID,
		Data:        data,
	}
}

// OccurredOn returns when the event occurred
func (e MockEvent) OccurredOn() time.Time {
	return e.occurredOn
}

// EventType returns the event type
func (e MockEvent) EventType() string {
	return e.Type
}

// EventTestBuilder helps build collections of events for testing
type EventTestBuilder struct {
	events []ddd.DomainEvent
}

// NewEventTestBuilder creates a new event test builder
func NewEventTestBuilder() *EventTestBuilder {
	return &EventTestBuilder{
		events: make([]ddd.DomainEvent, 0),
	}
}

// WithEvent adds an event to the builder
func (etb *EventTestBuilder) WithEvent(event ddd.DomainEvent) *EventTestBuilder {
	etb.events = append(etb.events, event)
	return etb
}

// WithMockEvent adds a mock event to the builder
func (etb *EventTestBuilder) WithMockEvent(eventType, aggregateID string) *EventTestBuilder {
	etb.events = append(etb.events, NewMockEvent(eventType, aggregateID, nil))
	return etb
}

// Build returns the built event collection
func (etb *EventTestBuilder) Build() []ddd.DomainEvent {
	return etb.events
}

// ExpectEvent finds and returns an event of the expected type, failing the test if not found
func ExpectEvent(
	t *testing.T,
	events []ddd.DomainEvent,
	expectedType string,
) ddd.DomainEvent {
	for _, event := range events {
		if event.EventType() == expectedType {
			return event
		}
	}
	t.Fatalf("expected event of type %s not found", expectedType)
	return nil
}

// ExpectNoEvent asserts that no event of the given type exists
func ExpectNoEvent(t *testing.T, events []ddd.DomainEvent, unexpectedType string) {
	for _, event := range events {
		if event.EventType() == unexpectedType {
			t.Fatalf("unexpected event of type %s found", unexpectedType)
		}
	}
}

// EventMatcher provides a fluent interface for matching events
type EventMatcher struct {
	matchers []func(ddd.DomainEvent) error
}

// NewEventMatcher creates a new event matcher
func NewEventMatcher() *EventMatcher {
	return &EventMatcher{
		matchers: make([]func(ddd.DomainEvent) error, 0),
	}
}

// WithType adds a type matcher
func (em *EventMatcher) WithType(expectedType string) *EventMatcher {
	em.matchers = append(em.matchers, func(event ddd.DomainEvent) error {
		if event.EventType() != expectedType {
			return fmt.Errorf("expected event type %s, got %s", expectedType, event.EventType())
		}
		return nil
	})
	return em
}

// Matches tests if an event matches all configured matchers
func (em *EventMatcher) Matches(event ddd.DomainEvent) error {
	for _, matcher := range em.matchers {
		if err := matcher(event); err != nil {
			return err
		}
	}
	return nil
}
