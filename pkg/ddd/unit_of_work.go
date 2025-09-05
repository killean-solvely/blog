package ddd

// UnitOfWork represents a business transaction that can span multiple aggregates
// It collects all changes and events, ensuring they are committed atomically
type UnitOfWork interface {
	// Register an aggregate for tracking within this unit of work
	Register(aggregate EventAggregate)

	// Commit persists all tracked aggregates and dispatches their events
	// Returns error if any operation fails, rolling back all changes
	Commit() error

	// Rollback discards all pending changes and events
	Rollback()

	// Events returns all events from all registered aggregates
	Events() []DomainEvent
}

// UnitOfWorkFactory creates new unit of work instances
type UnitOfWorkFactory interface {
	Create() UnitOfWork
}
