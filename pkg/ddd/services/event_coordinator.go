package services

import (
	"context"
	"fmt"
	"sync"

	"blog/pkg/ddd"
)

// EventCoordinator coordinates complex event processing workflows
type EventCoordinator struct {
	handlers      map[string][]ddd.EventHandlerFunc
	sagas         map[string][]Saga
	errorHandlers map[string]ErrorHandler
	mu            sync.RWMutex
}

// Saga represents a long-running business process that can handle events
type Saga interface {
	Handle(ctx context.Context, event ddd.DomainEvent) ([]ddd.DomainEvent, error)
	GetHandledEventTypes() []string
}

// ErrorHandler handles errors that occur during event processing
type ErrorHandler func(event ddd.DomainEvent, err error) error

// NewEventCoordinator creates a new event coordinator
func NewEventCoordinator() *EventCoordinator {
	return &EventCoordinator{
		handlers:      make(map[string][]ddd.EventHandlerFunc),
		sagas:         make(map[string][]Saga),
		errorHandlers: make(map[string]ErrorHandler),
	}
}

// RegisterHandler registers an event handler for a specific event type
func (ec *EventCoordinator) RegisterHandler(eventType string, handler ddd.EventHandlerFunc) {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.handlers[eventType] = append(ec.handlers[eventType], handler)
}

// RegisterSaga registers a saga for handling events
func (ec *EventCoordinator) RegisterSaga(saga Saga) {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	for _, eventType := range saga.GetHandledEventTypes() {
		ec.sagas[eventType] = append(ec.sagas[eventType], saga)
	}
}

// RegisterErrorHandler registers an error handler for a specific event type
func (ec *EventCoordinator) RegisterErrorHandler(eventType string, handler ErrorHandler) {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.errorHandlers[eventType] = handler
}

// Coordinate processes events through handlers and sagas, handling errors appropriately
func (ec *EventCoordinator) Coordinate(ctx context.Context, events ...ddd.DomainEvent) error {
	var resultEvents []ddd.DomainEvent

	for _, event := range events {
		if err := ec.handleEvent(ctx, event); err != nil {
			if errHandler, exists := ec.errorHandlers[event.EventType()]; exists {
				if handlerErr := errHandler(event, err); handlerErr != nil {
					return fmt.Errorf("error handler failed: %w", handlerErr)
				}
			} else {
				return err
			}
		}

		sagaEvents, err := ec.handleSagas(ctx, event)
		if err != nil {
			return fmt.Errorf("saga handling failed: %w", err)
		}
		resultEvents = append(resultEvents, sagaEvents...)
	}

	if len(resultEvents) > 0 {
		return ec.Coordinate(ctx, resultEvents...)
	}

	return nil
}

// handleEvent processes an event through all registered handlers
func (ec *EventCoordinator) handleEvent(ctx context.Context, event ddd.DomainEvent) error {
	ec.mu.RLock()
	handlers := ec.handlers[event.EventType()]
	ec.mu.RUnlock()

	for _, handler := range handlers {
		if err := handler(event); err != nil {
			return fmt.Errorf("handler failed for event %s: %w", event.EventType(), err)
		}
	}

	return nil
}

// handleSagas processes an event through all registered sagas
func (ec *EventCoordinator) handleSagas(
	ctx context.Context,
	event ddd.DomainEvent,
) ([]ddd.DomainEvent, error) {
	ec.mu.RLock()
	sagas := ec.sagas[event.EventType()]
	ec.mu.RUnlock()

	var resultEvents []ddd.DomainEvent

	for _, saga := range sagas {
		events, err := saga.Handle(ctx, event)
		if err != nil {
			return nil, fmt.Errorf("saga failed for event %s: %w", event.EventType(), err)
		}
		resultEvents = append(resultEvents, events...)
	}

	return resultEvents, nil
}

// WorkflowCoordinator coordinates complex multi-step workflows with compensation
type WorkflowCoordinator struct {
	steps         []WorkflowStep
	compensations []CompensationStep
	state         WorkflowState
}

// WorkflowStep represents a single step in a workflow
type WorkflowStep struct {
	Name     string
	Execute  func(ctx context.Context, data interface{}) (interface{}, error)
	Validate func(result interface{}) error
}

// CompensationStep represents compensation logic for a workflow step
type CompensationStep struct {
	Name       string
	Compensate func(ctx context.Context, data interface{}) error
}

// WorkflowState tracks the current state of workflow execution
type WorkflowState struct {
	CurrentStep    int
	CompletedSteps []string
	Data           map[string]interface{}
	Failed         bool
}

// NewWorkflowCoordinator creates a new workflow coordinator
func NewWorkflowCoordinator(
	steps []WorkflowStep,
	compensations []CompensationStep,
) *WorkflowCoordinator {
	return &WorkflowCoordinator{
		steps:         steps,
		compensations: compensations,
		state: WorkflowState{
			Data: make(map[string]interface{}),
		},
	}
}

// Execute runs the workflow steps in sequence with compensation on failure
func (wc *WorkflowCoordinator) Execute(ctx context.Context, initialData interface{}) error {
	wc.state.Data["initial"] = initialData

	for i, step := range wc.steps {
		wc.state.CurrentStep = i

		result, err := step.Execute(ctx, wc.state.Data)
		if err != nil {
			wc.state.Failed = true
			if compensateErr := wc.compensate(ctx); compensateErr != nil {
				return fmt.Errorf("compensation failed after step %s: %w", step.Name, compensateErr)
			}
			return fmt.Errorf("step %s failed: %w", step.Name, err)
		}

		if step.Validate != nil {
			if err := step.Validate(result); err != nil {
				wc.state.Failed = true
				if compensateErr := wc.compensate(ctx); compensateErr != nil {
					return fmt.Errorf(
						"compensation failed after validation of step %s: %w",
						step.Name,
						compensateErr,
					)
				}
				return fmt.Errorf("validation failed for step %s: %w", step.Name, err)
			}
		}

		wc.state.Data[step.Name] = result
		wc.state.CompletedSteps = append(wc.state.CompletedSteps, step.Name)
	}

	return nil
}

// compensate runs compensation logic for completed steps in reverse order
func (wc *WorkflowCoordinator) compensate(ctx context.Context) error {
	for i := len(wc.state.CompletedSteps) - 1; i >= 0; i-- {
		stepName := wc.state.CompletedSteps[i]

		for _, compensation := range wc.compensations {
			if compensation.Name == stepName {
				if err := compensation.Compensate(ctx, wc.state.Data[stepName]); err != nil {
					return fmt.Errorf("compensation for step %s failed: %w", stepName, err)
				}
				break
			}
		}
	}

	return nil
}

// EventRouter provides conditional event routing
type EventRouter struct {
	routes map[string][]EventRoute
	mu     sync.RWMutex
}

// EventRoute represents a conditional route for an event
type EventRoute struct {
	Condition func(event ddd.DomainEvent) bool
	Handler   ddd.EventHandlerFunc
}

// NewEventRouter creates a new event router
func NewEventRouter() *EventRouter {
	return &EventRouter{
		routes: make(map[string][]EventRoute),
	}
}

// AddRoute adds a conditional route for an event type
func (er *EventRouter) AddRoute(
	eventType string,
	condition func(ddd.DomainEvent) bool,
	handler ddd.EventHandlerFunc,
) {
	er.mu.Lock()
	defer er.mu.Unlock()

	er.routes[eventType] = append(er.routes[eventType], EventRoute{
		Condition: condition,
		Handler:   handler,
	})
}

// Route processes an event through all matching routes
func (er *EventRouter) Route(event ddd.DomainEvent) error {
	er.mu.RLock()
	routes := er.routes[event.EventType()]
	er.mu.RUnlock()

	for _, route := range routes {
		if route.Condition(event) {
			if err := route.Handler(event); err != nil {
				return fmt.Errorf("route handler failed: %w", err)
			}
		}
	}

	return nil
}
