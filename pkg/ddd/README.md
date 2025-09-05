# DDD Framework Package

This package provides reusable Domain-Driven Design (DDD) infrastructure components that can be used across multiple projects. It implements core DDD patterns including domain events, aggregates, repositories, specifications, validation, and unit of work.

## Quick Start

### 1. Define Your Aggregate

```go
package domain

import "github.com/your-org/your-project/pkg/ddd"

type Order struct {
    *ddd.AggregateBase
    customerID string
    items      []OrderItem
    status     OrderStatus
}

func NewOrder(id, customerID string) *Order {
    order := &Order{
        AggregateBase: &ddd.AggregateBase{},
        customerID:    customerID,
        items:         []OrderItem{},
        status:        OrderStatusPending,
    }
    order.SetID(id)
    
    // Raise creation event
    event := NewOrderCreated(id, customerID)
    ddd.EventRegistry.MustValidateEvent(event)
    order.RecordEvent(event)
    
    return order
}
```

### 2. Define and Register Events

```go
package domain

import (
    "time"
    "github.com/your-org/your-project/pkg/ddd"
)

type OrderCreated struct {
    occurredOn time.Time
    OrderID    string
    CustomerID string
}

func NewOrderCreated(orderID, customerID string) OrderCreated {
    return OrderCreated{
        occurredOn: time.Now(),
        OrderID:    orderID,
        CustomerID: customerID,
    }
}

func (e OrderCreated) OccurredOn() time.Time { return e.occurredOn }
func (e OrderCreated) EventType() string     { return "OrderCreated" }

// Register in init()
func init() {
    ddd.EventRegistry.Register(
        NewOrderCreated("", ""),
        "Raised when a new order is created",
    )
}
```

### 3. Create Repository

```go
package infrastructure

import (
    "github.com/your-org/your-project/pkg/ddd"
    "github.com/your-org/your-project/pkg/ddd/memory"
    "github.com/your-org/your-project/internal/domain"
)

type OrderRepository struct {
    orders map[string]*domain.Order
}

func NewOrderRepository() *OrderRepository {
    return &OrderRepository{
        orders: make(map[string]*domain.Order),
    }
}

func (r *OrderRepository) Save(order *domain.Order) error {
    r.orders[order.GetID()] = order
    return nil
}

func (r *OrderRepository) FindByID(id string) (*domain.Order, error) {
    order, exists := r.orders[id]
    if !exists {
        return nil, domain.ErrOrderNotFound
    }
    return order, nil
}
```

### 4. Set Up Event Handling

```go
// Create dispatcher
dispatcher := memory.NewInMemoryEventDispatcher(logger)

// Subscribe handlers using function-based approach
dispatcher.Subscribe("OrderCreated", func(event ddd.DomainEvent) error {
    e, ok := event.(domain.OrderCreated)
    if !ok {
        return fmt.Errorf("unexpected event type")
    }
    // Handle the event
    logger.Infow("Order created", "orderID", e.OrderID)
    return nil
})

// Or use method references
orderHandler := NewOrderEventHandler(logger)
dispatcher.Subscribe("OrderCreated", orderHandler.HandleOrderCreated)
dispatcher.Subscribe("OrderShipped", orderHandler.HandleOrderShipped)
```

### 5. Use Validation

```go
import "github.com/your-org/your-project/pkg/ddd/validation"

type CreateOrderRequest struct {
    CustomerID string
    Items      []OrderItemRequest
}

func (r CreateOrderRequest) Validate() error {
    errors := validation.NewErrors()
    validator := validation.New()
    
    if err := validator.Required(r.CustomerID, "customerID"); err != nil {
        errors.ValidationErrors = append(errors.ValidationErrors, *err)
    }
    
    if len(r.Items) == 0 {
        errors.Add("items", "at least one item is required", "required")
    }
    
    if errors.HasErrors() {
        return errors
    }
    return nil
}
```

## Package Structure

```
pkg/ddd/
├── README.md                    # This file
├── aggregate_base.go           # Generic aggregate base class with event handling
├── events.go                   # Core event interfaces
├── event_registry.go           # Event registration and validation
├── unit_of_work.go            # Unit of Work interface
├── memory/                     # In-memory implementations
│   ├── dispatcher.go          # Function-based event dispatcher
│   └── unit_of_work.go       # In-memory unit of work implementation
├── services/                   # Advanced domain services
│   └── event_coordinator.go   # Event coordination and saga patterns
├── specifications/             # Specification pattern implementation
│   └── specification.go       # Generic specification combinators
├── testing/                    # Testing utilities
│   └── event_helpers.go       # Event testing infrastructure
└── validation/                 # Domain validation framework
    ├── doc.go                 # Package documentation
    ├── errors.go              # Validation error structures
    └── validator.go           # Common validation utilities
```

## Core Components

### 1. AggregateBase

Provides event recording capabilities for domain aggregates with ID tracking.

```go
import "github.com/your-org/your-project/pkg/ddd"

type YourAggregate struct {
    *ddd.AggregateBase
    // your fields...
}

func NewYourAggregate(id string) *YourAggregate {
    aggregate := &YourAggregate{
        AggregateBase: &ddd.AggregateBase{},
    }
    aggregate.SetID(id)
    return aggregate
}

func (a *YourAggregate) DoSomething() {
    // business logic...
    event := NewSomethingHappened(a.GetID())
    ddd.EventRegistry.MustValidateEvent(event)
    a.RecordEvent(event)
}

// Implement AggregateRoot interface
func (a *YourAggregate) GetID() string {
    return a.AggregateBase.GetID()
}
```

### 2. Event Registry

Centralized event registration and validation system.

```go
// Register events during initialization
func init() {
    ddd.EventRegistry.Register(
        NewSomethingHappened("", nil),
        "Description of what this event represents",
    )
}

// Validate events before dispatch
func (a *YourAggregate) DoSomething() {
    event := NewSomethingHappened(a.id)
    ddd.EventRegistry.MustValidateEvent(event) // Panics if not registered
    a.RecordEvent(event)
}
```

### 3. Specifications

Generic, composable business rules using Go generics.

```go
import "github.com/your-org/your-project/pkg/ddd/specifications"

type YourEntitySpecification = specifications.Specification[YourEntity]

type HighPrioritySpec struct{}

func (s HighPrioritySpec) IsSatisfiedBy(entity YourEntity) bool {
    return entity.Priority() == High
}

func (s HighPrioritySpec) String() string {
    return "HighPriority"
}

// Combine specifications
spec := specifications.And(
    HighPrioritySpec{},
    specifications.Not(CompletedSpec{}),
)
```

### 4. Unit of Work

Atomic operations across multiple aggregates with automatic event collection.

```go
import (
    "github.com/your-org/your-project/pkg/ddd"
    "github.com/your-org/your-project/pkg/ddd/memory"
)

// Create factory
factory := memory.NewUnitOfWorkFactory(dispatcher)

// Use in application service
func (s *YourService) MultiAggregateOperation() error {
    uow := s.uowFactory.Create()

    // Load and modify aggregates
    aggregate1 := s.repo1.FindByID(id1)
    aggregate2 := s.repo2.FindByID(id2)

    aggregate1.DoSomething()
    aggregate2.DoSomethingElse()

    // Register for atomic commit
    uow.Register(aggregate1)
    uow.Register(aggregate2)

    // Commit saves all aggregates and dispatches all events atomically
    return uow.Commit()
}
```

### 5. Validation Framework

Type-safe, manual validation approach for domain objects.

```go
import "github.com/your-org/your-project/pkg/ddd/validation"

type CreateUserRequest struct {
    Username string
    Email    string
    Age      int
}

func (r CreateUserRequest) Validate() *validation.Errors {
    errors := validation.NewErrors()
    v := validation.NewValidator(errors)

    v.Required("username", r.Username)
    v.MinLength("username", r.Username, 3)
    v.MaxLength("username", r.Username, 50)

    v.Required("email", r.Email)
    v.Email("email", r.Email)

    v.MinValue("age", r.Age, 18)
    v.MaxValue("age", r.Age, 120)

    return errors
}

// Usage
request := CreateUserRequest{Username: "ab", Email: "invalid", Age: 15}
if errors := request.Validate(); errors.HasErrors() {
    // Handle validation errors
    for _, err := range errors.All() {
        fmt.Printf("Field: %s, Error: %s\n", err.Field, err.Message)
    }
}
```

### 6. Event Dispatcher

The framework provides a function-based event dispatcher for handling domain events:

```go
import "github.com/your-org/your-project/pkg/ddd/memory"

// Create dispatcher
dispatcher := memory.NewInMemoryEventDispatcher(logger)

// Define handler functions
type EventHandlerFunc func(event ddd.DomainEvent) error

// Subscribe handlers - multiple ways:

// 1. Anonymous functions
dispatcher.Subscribe("OrderCreated", func(event ddd.DomainEvent) error {
    orderCreated, ok := event.(domain.OrderCreated)
    if !ok {
        return fmt.Errorf("unexpected event type")
    }
    // Send confirmation email
    return emailService.SendOrderConfirmation(orderCreated.OrderID)
})

// 2. Method references
handler := NewOrderEventHandler(logger, emailService)
dispatcher.Subscribe("OrderCreated", handler.HandleOrderCreated)
dispatcher.Subscribe("OrderShipped", handler.HandleOrderShipped)

// 3. Multiple handlers for same event
dispatcher.Subscribe("OrderCreated", handler.HandleOrderCreated)
dispatcher.Subscribe("OrderCreated", metricsHandler.RecordOrderCreated)
dispatcher.Subscribe("OrderCreated", auditHandler.LogOrderCreated)
```

Event Handler Example:
```go
type OrderEventHandler struct {
    logger       *zap.SugaredLogger
    emailService EmailService
}

func (h *OrderEventHandler) HandleOrderCreated(event ddd.DomainEvent) error {
    e, ok := event.(domain.OrderCreated)
    if !ok {
        return fmt.Errorf("unexpected event type")
    }
    
    h.logger.Infow("Processing order created event", 
        "orderID", e.OrderID,
        "customerID", e.CustomerID,
    )
    
    // Side effects
    if err := h.emailService.SendOrderConfirmation(e.OrderID); err != nil {
        h.logger.Errorw("Failed to send confirmation email", "error", err)
        // Decide whether to fail or continue
        return err
    }
    
    return nil
}
```

### 7. Repository Patterns

Implement domain repositories for aggregate persistence:

```go
package infrastructure

type YourRepository struct {
    aggregates map[string]*domain.YourAggregate
}

func NewYourRepository() *YourRepository {
    return &YourRepository{
        aggregates: make(map[string]*domain.YourAggregate),
    }
}

func (r *YourRepository) Save(aggregate *domain.YourAggregate) error {
    r.aggregates[aggregate.GetID()] = aggregate
    return nil
}

func (r *YourRepository) FindByID(id string) (*domain.YourAggregate, error) {
    aggregate, exists := r.aggregates[id]
    if !exists {
        return nil, domain.ErrAggregateNotFound
    }
    return aggregate, nil
}
```

### 8. Event Testing

Comprehensive testing utilities for event-driven behavior.

```go
import (
    "testing"
    "github.com/your-org/your-project/pkg/ddd/testing"
)

func TestAggregateEvents(t *testing.T) {
    // Record events for testing
    recorder := ddtesting.NewEventRecorder()

    // Perform operations that raise events
    aggregate := NewYourAggregate()
    aggregate.DoSomething()

    // Dispatch events to recorder
    for _, event := range aggregate.GetUncommittedEvents() {
        recorder.Dispatch(event)
    }

    // Assert events were raised
    ddtesting.AssertEvents(t, recorder).
        HasEventType("SomethingHappened").
        HasEventCount("SomethingHappened", 1).
        HasTotalEventCount(1)
}

// Test event handlers
func TestOrderEventHandler(t *testing.T) {
    // Create mock services
    emailService := &MockEmailService{}
    handler := NewOrderEventHandler(logger, emailService)
    
    // Create test event
    event := domain.NewOrderCreated("order-123", "customer-456")
    
    // Test handler
    err := handler.HandleOrderCreated(event)
    assert.NoError(t, err)
    
    // Verify side effects
    assert.True(t, emailService.SendOrderConfirmationCalled)
    assert.Equal(t, "order-123", emailService.LastOrderID)
}
```

## Repository Patterns

Implement domain repositories for aggregate persistence:

```go
package domain

// Repository interface in domain layer
type YourAggregateRepository interface {
    Save(aggregate *YourAggregate) error
    FindByID(id string) (*YourAggregate, error)
    FindAll() ([]*YourAggregate, error)
    Delete(id string) error
}
```

### In-Memory Repository Implementation

```go
package infrastructure

import "github.com/your-org/your-project/internal/domain"

type InMemoryYourAggregateRepository struct {
    aggregates map[string]*domain.YourAggregate
}

func NewInMemoryYourAggregateRepository() *InMemoryYourAggregateRepository {
    return &InMemoryYourAggregateRepository{
        aggregates: make(map[string]*domain.YourAggregate),
    }
}

func (r *InMemoryYourAggregateRepository) Save(aggregate *domain.YourAggregate) error {
    r.aggregates[aggregate.GetID()] = aggregate
    return nil
}

func (r *InMemoryYourAggregateRepository) FindByID(id string) (*domain.YourAggregate, error) {
    aggregate, exists := r.aggregates[id]
    if !exists {
        return nil, domain.ErrYourAggregateNotFound
    }
    return aggregate, nil
}

func (r *InMemoryYourAggregateRepository) FindAll() ([]*domain.YourAggregate, error) {
    var aggregates []*domain.YourAggregate
    for _, aggregate := range r.aggregates {
        aggregates = append(aggregates, aggregate)
    }
    return aggregates, nil
}

func (r *InMemoryYourAggregateRepository) Delete(id string) error {
    delete(r.aggregates, id)
    return nil
}
```

## Advanced Patterns

### Event Coordination

For complex event workflows and saga patterns:

```go
import "github.com/your-org/your-project/pkg/ddd/services"

coordinator := services.NewEventCoordinator()

// Register handlers using function-based approach
coordinator.RegisterHandler("OrderCreated", func(event ddd.DomainEvent) error {
    // Handle order creation logic
    return nil
})

// Or register handler methods
coordinator.RegisterHandler("OrderCreated", orderSaga.HandleOrderCreated)
coordinator.RegisterHandler("PaymentProcessed", orderSaga.HandlePaymentProcessed)

// Process events with coordination
err := coordinator.Coordinate(ctx, events...)
```

### Workflow Coordination

For multi-step processes with compensation:

```go
steps := []services.WorkflowStep{
    {
        Name:     "step1",
        Execute:  step1Func,
        Validate: validateStep1,
    },
    // more steps...
}

compensations := []services.CompensationStep{
    {
        Name:       "step1",
        Compensate: compensateStep1,
    },
    // more compensations...
}

workflow := services.NewWorkflowCoordinator(steps, compensations)
err := workflow.Execute(ctx, initialData)
```

## Integration with Your Domain

### 1. Extend Generic Components

```go
// internal/infrastructure/memory/unit_of_work.go
package memory

import (
    "github.com/your-org/your-project/pkg/ddd/memory"
    "github.com/your-org/your-project/internal/domain"
)

// YourAwareUnitOfWork extends generic UoW with domain-specific functionality
type YourAwareUnitOfWork struct {
    *memory.UnitOfWork
    yourRepo domain.YourRepository
}

func NewYourAwareUnitOfWork(dispatcher domain.EventDispatcher, yourRepo domain.YourRepository) *YourAwareUnitOfWork {
    uow := memory.NewUnitOfWork(dispatcher)
    uow.RegisterRepository("*domain.YourEntity", &yourRepositoryAdapter{yourRepo})

    return &YourAwareUnitOfWork{
        UnitOfWork: uow,
        yourRepo:   yourRepo,
    }
}
```

### 2. Domain-Specific Specifications

```go
// internal/domain/specifications/your_specifications.go
package specifications

import (
    "github.com/your-org/your-project/pkg/ddd/specifications"
    "github.com/your-org/your-project/internal/domain"
)

type YourEntitySpecification = specifications.Specification[domain.YourEntity]

// Your specific business rules
type HighPrioritySpec struct{}
type CompletedSpec struct{}
// etc.

// Convenience functions
func And(left, right YourEntitySpecification) YourEntitySpecification {
    return specifications.And(left, right)
}
```

## Common Patterns & Recipes

### Complete Domain Aggregate

```go
package domain

import (
    "fmt"
    "github.com/your-org/your-project/pkg/ddd"
)

type Product struct {
    *ddd.AggregateBase
    name        string
    price       Money
    inventory   int
    isActive    bool
}

// Factory
func NewProduct(id, name string, price Money) (*Product, error) {
    if name == "" {
        return nil, ErrProductNameRequired
    }
    if price.Amount <= 0 {
        return nil, ErrInvalidPrice
    }

    product := &Product{
        AggregateBase: &ddd.AggregateBase{},
        name:          name,
        price:         price,
        inventory:     0,
        isActive:      true,
    }
    product.SetID(id)

    // Record creation event for side effects
    event := NewProductCreated(id, name, price)
    ddd.EventRegistry.MustValidateEvent(event)
    product.RecordEvent(event)

    return product, nil
}

// Business methods
func (p *Product) UpdatePrice(newPrice Money) error {
    if !p.isActive {
        return ErrProductInactive
    }
    if newPrice.Amount <= 0 {
        return ErrInvalidPrice
    }

    oldPrice := p.price
    p.price = newPrice

    event := NewPriceUpdated(p.GetID(), oldPrice, newPrice)
    ddd.EventRegistry.MustValidateEvent(event)
    p.RecordEvent(event)

    return nil
}

func (p *Product) AddInventory(quantity int) error {
    if quantity <= 0 {
        return ErrInvalidQuantity
    }

    p.inventory += quantity

    event := NewInventoryAdded(p.GetID(), quantity)
    ddd.EventRegistry.MustValidateEvent(event)
    p.RecordEvent(event)

    return nil
}

func (p *Product) Deactivate() error {
    if !p.isActive {
        return ErrProductAlreadyInactive
    }

    p.isActive = false

    event := NewProductDeactivated(p.GetID())
    ddd.EventRegistry.MustValidateEvent(event)
    p.RecordEvent(event)

    return nil
}

// Read-only getters
func (p *Product) GetID() string      { return p.AggregateBase.GetID() }
func (p *Product) Name() string       { return p.name }
func (p *Product) Price() Money       { return p.price }
func (p *Product) Inventory() int     { return p.inventory }
func (p *Product) IsActive() bool     { return p.isActive }
```

### Testing Domain Aggregates

```go
func TestProduct_UpdatePrice(t *testing.T) {
    // Arrange
    product, _ := NewProduct("prod-1", "Widget", Money{100})
    product.MarkEventsAsCommitted() // Clear creation event
    
    recorder := testing.NewEventRecorder()
    
    // Act
    err := product.UpdatePrice(Money{150})
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, Money{150}, product.Price())
    
    // Verify events
    for _, event := range product.GetUncommittedEvents() {
        recorder.Dispatch(event)
    }
    
    testing.AssertEvents(t, recorder).
        HasEventType("PriceUpdated").
        HasEventCount("PriceUpdated", 1).
        HasTotalEventCount(1)
}
```

### Specification Combinations

```go
// Complex business rule using specifications
spec := specifications.And(
    IsActiveSpec{},
    specifications.Or(
        HasInventorySpec{MinQuantity: 10},
        IsFeaturedSpec{},
    ),
    specifications.Not(IsDiscontinuedSpec{}),
)

// Use in repository or service
products := repo.FindAll()
var eligibleProducts []*Product
for _, p := range products {
    if spec.IsSatisfiedBy(*p) {
        eligibleProducts = append(eligibleProducts, p)
    }
}
```

## Migration Guide

### From Traditional CRUD to DDD

1. **Start with State-Based Aggregate**: Begin with traditional aggregate
2. **Add Event Recording**: Record events for side effects and notifications
3. **Implement Business Rules**: Move validation and business logic into aggregates
4. **Create Repository Interfaces**: Define repository contracts in domain layer
5. **Add Application Services**: Orchestrate use cases and coordinate repositories

### From Other DDD Frameworks

| Other Framework | This Package |
|----------------|--------------|
| `IAggregate` | `AggregateRoot` interface |
| `IDomainEvent` | `DomainEvent` interface |
| `IRepository<T>` | Domain-specific repository interfaces |
| `ISpecification<T>` | `specifications.Specification[T]` |
| `IUnitOfWork` | `ddd.UnitOfWork` |
| `ValueObject` | Plain Go types with validation |

## Troubleshooting

### Common Issues

1. **"Event not registered" panic**
   - Ensure events are registered in init()
   - Check import order - domain package must be imported

2. **Anemic domain models**
   - Move business logic from services into aggregates
   - Use "Tell, Don't Ask" principle - methods not getters/setters

3. **Cross-aggregate transactions**
   - Use eventual consistency with domain events
   - Implement saga patterns for complex workflows

4. **Validation errors not user-friendly**
   - Use validation.Errors for structured errors
   - Map to HTTP/gRPC error codes at boundaries

## Benefits

1. **Reusability**: Use the same DDD infrastructure across multiple projects
2. **Type Safety**: Generic specifications provide compile-time type checking
3. **Consistency**: Standardized patterns across your organization
4. **Testing**: Comprehensive testing utilities included
5. **Maintainability**: Centralized DDD infrastructure with clear interfaces
6. **Domain Purity**: Domain-specific wrappers maintain clean architecture
7. **Event-Driven Architecture**: Support for domain events and side effects
8. **In-Memory Options**: Quick prototyping with in-memory implementations

## Best Practices

1. **Event Registration**: Always register events during package initialization
2. **Event Validation**: Use `MustValidateEvent` in aggregates to catch registration issues early
3. **Event Handling**: Use function-based handlers for flexibility - can be anonymous functions, methods, or standalone functions
4. **Handler Organization**: Group related handlers in handler structs in the infrastructure layer
5. **Unit of Work**: Use for multi-aggregate operations to ensure atomicity
6. **Testing**: Use the provided testing utilities for comprehensive event testing
7. **Direct Imports**: Import pkg/ddd directly in domain layer for pragmatic DDD implementation
8. **Repository Interfaces**: Define repository contracts in domain layer, implement in infrastructure
11. **Business Logic**: Keep all business rules within domain aggregates, not in services
12. **Tell Don't Ask**: Use command methods like `order.Ship()` instead of `order.SetStatus("shipped")`
9. **Event Naming**: Use past tense for events (e.g., "OrderCreated" not "CreateOrder")
10. **Aggregate Boundaries**: Keep aggregates small and focused on a single consistency boundary

This package provides a solid foundation for implementing DDD patterns while maintaining flexibility for domain-specific customizations.