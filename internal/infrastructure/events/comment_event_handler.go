package events

import (
	"errors"
	"log"

	"blog/internal/domain"
	"blog/pkg/ddd"
)

type CommentEventHandler struct{}

func NewCommentEventHandler() *CommentEventHandler {
	return &CommentEventHandler{}
}

func (h CommentEventHandler) Register(dispatcher ddd.EventDispatcher) {
	dispatcher.Subscribe(
		domain.CommentCreatedEventType.String(),
		h.HandleCommentCreated,
	)

	dispatcher.Subscribe(
		domain.CommentEditedEventType.String(),
		h.HandleCommentEdited,
	)

	dispatcher.Subscribe(
		domain.CommentArchivedEventType.String(),
		h.HandleCommentArchived,
	)
}

func (h CommentEventHandler) HandleCommentCreated(event ddd.DomainEvent) error {
	e, ok := event.(*domain.CommentCreatedEvent)
	if !ok {
		return errors.New("invalid event type")
	}

	log.Printf(
		"CommentCreatedEvent handled for ID: %s",
		e.CommentID.String(),
	)

	return nil
}

func (h CommentEventHandler) HandleCommentEdited(event ddd.DomainEvent) error {
	e, ok := event.(*domain.CommentEditedEvent)
	if !ok {
		return errors.New("invalid event type")
	}

	log.Printf(
		"CommentEditedEvent handled for ID: %s",
		e.CommentID.String(),
	)

	return nil
}

func (h CommentEventHandler) HandleCommentArchived(event ddd.DomainEvent) error {
	e, ok := event.(*domain.CommentArchivedEvent)
	if !ok {
		return errors.New("invalid event type")
	}

	log.Printf(
		"CommentArchivedEvent handled for ID: %s",
		e.CommentID.String(),
	)

	return nil
}
