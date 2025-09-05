package memory

import (
	"blog/pkg/ddd"

	"go.uber.org/zap"
)

// InMemoryEventDispatcher is a simple in-memory implementation of EventDispatcher
type InMemoryEventDispatcher struct {
	handlers map[string][]ddd.EventHandlerFunc
	logger   *zap.SugaredLogger
}

// NewInMemoryEventDispatcher creates a new in-memory event dispatcher
func NewInMemoryEventDispatcher(log *zap.SugaredLogger) *InMemoryEventDispatcher {
	return &InMemoryEventDispatcher{
		handlers: make(map[string][]ddd.EventHandlerFunc),
		logger:   log,
	}
}

// Subscribe registers a handler for a specific event type
func (d *InMemoryEventDispatcher) Subscribe(eventType string, handler ddd.EventHandlerFunc) {
	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

// Dispatch sends an event to all registered handlers for its type
func (d *InMemoryEventDispatcher) Dispatch(event ddd.DomainEvent) error {
	// Validate event is registered if using the global registry
	registry := ddd.GetEventRegistry()
	if err := registry.ValidateEvent(event); err != nil {
		d.logger.Warnw("Event not registered in global registry",
			"event_type", event.EventType(),
			"error", err,
		)
		// In production, you might want to handle this differently
		// For now, we'll log and continue
	}

	for _, handler := range d.handlers[event.EventType()] {
		if err := handler(event); err != nil {
			d.logger.Errorw("Error handling event",
				"event_type", event.EventType(),
				"error", err,
			)

			return err
		}
	}

	return nil
}
