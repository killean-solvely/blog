package events

import (
	"errors"
	"log"

	"blog/internal/domain"
	"blog/pkg/ddd"
)

type UserEventHandler struct{}

func NewUserEventHandler() *UserEventHandler {
	return &UserEventHandler{}
}

func (h UserEventHandler) Register(dispatcher ddd.EventDispatcher) {
	dispatcher.Subscribe(
		domain.UserCreatedEventType.String(),
		h.HandleUserCreated,
	)

	dispatcher.Subscribe(
		domain.UserRoleAddedEventType.String(),
		h.HandleUserRoleAdded,
	)

	dispatcher.Subscribe(
		domain.UserRoleRemovedEventType.String(),
		h.HandleUserRoleRemoved,
	)

	dispatcher.Subscribe(
		domain.UserDescriptionUpdatedEventType.String(),
		h.HandleUserDescriptionUpdated,
	)

	dispatcher.Subscribe(
		domain.UserPasswordUpdatedEventType.String(),
		h.HandleUserPasswordUpdated,
	)
}

func (h UserEventHandler) HandleUserCreated(event ddd.DomainEvent) error {
	e, ok := event.(*domain.UserCreatedEvent)
	if !ok {
		return errors.New("invalid event type")
	}

	log.Printf(
		"UserCreatedEvent handled for ID: %s",
		e.UserID.String(),
	)

	return nil
}

func (h UserEventHandler) HandleUserRoleAdded(event ddd.DomainEvent) error {
	e, ok := event.(*domain.UserRoleAddedEvent)
	if !ok {
		return errors.New("invalid event type")
	}

	log.Printf(
		"UserRoleAddedEvent handled for ID: %s",
		e.UserID.String(),
	)

	return nil
}

func (h UserEventHandler) HandleUserRoleRemoved(event ddd.DomainEvent) error {
	e, ok := event.(*domain.UserRoleRemovedEvent)
	if !ok {
		return errors.New("invalid event type")
	}

	log.Printf(
		"UserRoleRemovedEvent handled for ID: %s",
		e.UserID.String(),
	)

	return nil
}

func (h UserEventHandler) HandleUserDescriptionUpdated(event ddd.DomainEvent) error {
	e, ok := event.(*domain.UserDescriptionUpdatedEvent)
	if !ok {
		return errors.New("invalid event type")
	}

	log.Printf(
		"UserDescriptionUpdatedEvent handled for ID: %s",
		e.UserID.String(),
	)

	return nil
}

func (h UserEventHandler) HandleUserPasswordUpdated(event ddd.DomainEvent) error {
	e, ok := event.(*domain.UserPasswordUpdatedEvent)
	if !ok {
		return errors.New("invalid event type")
	}

	log.Printf(
		"UserPasswordUpdatedEvent handled for ID: %s",
		e.UserID.String(),
	)

	return nil
}
