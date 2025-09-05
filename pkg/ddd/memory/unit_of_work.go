package memory

import (
	"fmt"
	"sync"

	"blog/pkg/ddd"
)

// Repository interface for saving aggregates
type Repository interface {
	Save(aggregate interface{}) error
}

// UnitOfWork provides a generic implementation of the UnitOfWork pattern
type UnitOfWork struct {
	aggregates      []ddd.EventAggregate
	repositories    map[string]Repository
	eventDispatcher ddd.EventDispatcher
	mu              sync.Mutex
	committed       bool
}

// NewUnitOfWork creates a new unit of work instance
func NewUnitOfWork(dispatcher ddd.EventDispatcher) *UnitOfWork {
	return &UnitOfWork{
		aggregates:      make([]ddd.EventAggregate, 0),
		repositories:    make(map[string]Repository),
		eventDispatcher: dispatcher,
	}
}

// Register implements ddd.UnitOfWork interface
func (uow *UnitOfWork) Register(aggregate ddd.EventAggregate) {
	uow.mu.Lock()
	defer uow.mu.Unlock()

	if uow.committed {
		panic("cannot register aggregate after commit")
	}

	uow.aggregates = append(uow.aggregates, aggregate)
}

// RegisterRepository adds a repository for a specific aggregate type
func (uow *UnitOfWork) RegisterRepository(typeName string, repo Repository) {
	uow.mu.Lock()
	defer uow.mu.Unlock()

	if uow.committed {
		panic("cannot register repository after commit")
	}

	uow.repositories[typeName] = repo
}

// Commit implements ddd.UnitOfWork interface
func (uow *UnitOfWork) Commit() error {
	uow.mu.Lock()
	defer uow.mu.Unlock()

	if uow.committed {
		return fmt.Errorf("unit of work already committed")
	}

	// First, save all aggregates to their repositories
	for _, aggregate := range uow.aggregates {
		if err := uow.saveAggregate(aggregate); err != nil {
			return fmt.Errorf("failed to save aggregate: %w", err)
		}
	}

	// Then collect all events
	allEvents := make([]ddd.DomainEvent, 0)
	for _, aggregate := range uow.aggregates {
		events := aggregate.GetUncommittedEvents()
		allEvents = append(allEvents, events...)
	}

	// Dispatch events after successful persistence
	if uow.eventDispatcher != nil && len(allEvents) > 0 {
		for _, event := range allEvents {
			uow.eventDispatcher.Dispatch(event)
		}
	}

	// Mark events as committed
	for _, aggregate := range uow.aggregates {
		aggregate.MarkEventsAsCommitted()
	}

	uow.committed = true
	return nil
}

// saveAggregate saves an aggregate using its registered repository
func (uow *UnitOfWork) saveAggregate(aggregate ddd.EventAggregate) error {
	typeName := fmt.Sprintf("%T", aggregate)
	if repo, exists := uow.repositories[typeName]; exists {
		return repo.Save(aggregate)
	}
	return fmt.Errorf("no repository configured for aggregate type: %s", typeName)
}

// Rollback implements ddd.UnitOfWork interface
func (uow *UnitOfWork) Rollback() {
	uow.mu.Lock()
	defer uow.mu.Unlock()

	for _, aggregate := range uow.aggregates {
		aggregate.MarkEventsAsCommitted()
	}

	uow.committed = true
}

// Events implements ddd.UnitOfWork interface
func (uow *UnitOfWork) Events() []ddd.DomainEvent {
	uow.mu.Lock()
	defer uow.mu.Unlock()

	allEvents := make([]ddd.DomainEvent, 0)
	for _, aggregate := range uow.aggregates {
		events := aggregate.GetUncommittedEvents()
		allEvents = append(allEvents, events...)
	}
	return allEvents
}

// IsCommitted returns whether the unit of work has been committed
func (uow *UnitOfWork) IsCommitted() bool {
	uow.mu.Lock()
	defer uow.mu.Unlock()
	return uow.committed
}

// EventOnlyUnitOfWork handles only events without persistence
type EventOnlyUnitOfWork struct {
	*UnitOfWork
}

// NewEventOnlyUnitOfWork creates a UoW that only handles events, not persistence
func NewEventOnlyUnitOfWork(dispatcher ddd.EventDispatcher) *EventOnlyUnitOfWork {
	return &EventOnlyUnitOfWork{
		UnitOfWork: &UnitOfWork{
			aggregates:      make([]ddd.EventAggregate, 0),
			repositories:    make(map[string]Repository),
			eventDispatcher: dispatcher,
		},
	}
}

// Commit only dispatches events without persistence
func (euow *EventOnlyUnitOfWork) Commit() error {
	euow.mu.Lock()
	defer euow.mu.Unlock()

	if euow.committed {
		return fmt.Errorf("unit of work already committed")
	}

	// Only collect and dispatch events, skip persistence
	allEvents := make([]ddd.DomainEvent, 0)
	for _, aggregate := range euow.aggregates {
		events := aggregate.GetUncommittedEvents()
		allEvents = append(allEvents, events...)
	}

	// Dispatch events
	if euow.eventDispatcher != nil && len(allEvents) > 0 {
		for _, event := range allEvents {
			euow.eventDispatcher.Dispatch(event)
		}
	}

	// Mark events as committed
	for _, aggregate := range euow.aggregates {
		aggregate.MarkEventsAsCommitted()
	}

	euow.committed = true
	return nil
}

// UnitOfWorkFactory creates unit of work instances
type UnitOfWorkFactory struct {
	dispatcher ddd.EventDispatcher
}

// NewUnitOfWorkFactory creates a new unit of work factory
func NewUnitOfWorkFactory(dispatcher ddd.EventDispatcher) *UnitOfWorkFactory {
	return &UnitOfWorkFactory{
		dispatcher: dispatcher,
	}
}

// Create implements ddd.UnitOfWorkFactory interface
func (f *UnitOfWorkFactory) Create() ddd.UnitOfWork {
	return NewUnitOfWork(f.dispatcher)
}

// TransactionalUnitOfWork extends UnitOfWork with transaction support
type TransactionalUnitOfWork struct {
	*UnitOfWork
	tx interface{}
}

// NewTransactionalUnitOfWork creates a new transactional unit of work
func NewTransactionalUnitOfWork(
	dispatcher ddd.EventDispatcher,
	tx interface{},
) *TransactionalUnitOfWork {
	return &TransactionalUnitOfWork{
		UnitOfWork: NewUnitOfWork(dispatcher),
		tx:         tx,
	}
}

// Commit extends the base commit with transaction handling
func (tuow *TransactionalUnitOfWork) Commit() error {
	if err := tuow.UnitOfWork.Commit(); err != nil {
		return err
	}

	// Transaction-specific commit logic would go here
	return nil
}

// GetTransaction returns the underlying transaction
func (tuow *TransactionalUnitOfWork) GetTransaction() interface{} {
	return tuow.tx
}

// ScopedUnitOfWork provides scoped execution with cleanup
type ScopedUnitOfWork struct {
	uow        *UnitOfWork
	onComplete func() error
}

// NewScopedUnitOfWork creates a new scoped unit of work
func NewScopedUnitOfWork(uow *UnitOfWork, onComplete func() error) *ScopedUnitOfWork {
	return &ScopedUnitOfWork{
		uow:        uow,
		onComplete: onComplete,
	}
}

// Register delegates to the underlying unit of work
func (s *ScopedUnitOfWork) Register(aggregate ddd.EventAggregate) {
	s.uow.Register(aggregate)
}

// Complete commits the unit of work and executes the completion callback
func (s *ScopedUnitOfWork) Complete() error {
	if err := s.uow.Commit(); err != nil {
		return err
	}

	if s.onComplete != nil {
		return s.onComplete()
	}

	return nil
}
