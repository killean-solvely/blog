package events

import (
	"errors"
	"log"

	"blog/internal/domain"
	"blog/pkg/ddd"
)

type PostEventHandler struct{}

func NewPostEventHandler() *PostEventHandler {
	return &PostEventHandler{}
}

func (h PostEventHandler) Register(dispatcher ddd.EventDispatcher) {
	dispatcher.Subscribe(
		domain.PostCreatedEventType.String(),
		h.HandlePostCreated,
	)

	dispatcher.Subscribe(
		domain.PostTitleEditedEventType.String(),
		h.HandlePostTitleEdited,
	)

	dispatcher.Subscribe(
		domain.PostContentEditedEventType.String(),
		h.HandlePostContentEdited,
	)

	dispatcher.Subscribe(
		domain.PostArchivedEventType.String(),
		h.HandlePostArchived,
	)
}

func (h PostEventHandler) HandlePostCreated(event ddd.DomainEvent) error {
	e, ok := event.(*domain.PostCreatedEvent)
	if !ok {
		return errors.New("invalid event type")
	}

	log.Printf(
		"PostCreatedEvent handled for ID: %s",
		e.PostID.String(),
	)

	return nil
}

func (h PostEventHandler) HandlePostTitleEdited(event ddd.DomainEvent) error {
	e, ok := event.(*domain.PostTitleEditedEvent)
	if !ok {
		return errors.New("invalid event type")
	}

	log.Printf(
		"PostTitleEditedEvent handled for ID: %s",
		e.PostID.String(),
	)

	return nil
}

func (h PostEventHandler) HandlePostContentEdited(event ddd.DomainEvent) error {
	e, ok := event.(*domain.PostContentEditedEvent)
	if !ok {
		return errors.New("invalid event type")
	}

	log.Printf(
		"PostContentEditedEvent handled for ID: %s",
		e.PostID.String(),
	)

	return nil
}

func (h PostEventHandler) HandlePostArchived(event ddd.DomainEvent) error {
	e, ok := event.(*domain.PostArchivedEvent)
	if !ok {
		return errors.New("invalid event type")
	}

	log.Printf(
		"PostArchivedEvent handled for ID: %s",
		e.PostID.String(),
	)

	return nil
}
