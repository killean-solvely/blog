package events

import (
	"errors"
	"log"

	"blog/internal/domain"
	"blog/pkg/ddd"
)

type RatingEventHandler struct{}

func NewRatingEventHandler() *RatingEventHandler {
	return &RatingEventHandler{}
}

func (h RatingEventHandler) Register(dispatcher ddd.EventDispatcher) {
	dispatcher.Subscribe(
		domain.RatingCreatedEventType.String(),
		h.HandleRatingCreated,
	)

	dispatcher.Subscribe(
		domain.RatingChangedEventType.String(),
		h.HandleRatingChanged,
	)

	dispatcher.Subscribe(
		domain.RatingRemovedEventType.String(),
		h.HandleRatingRemoved,
	)
}

func (h RatingEventHandler) HandleRatingCreated(event ddd.DomainEvent) error {
	e, ok := event.(*domain.RatingCreatedEvent)
	if !ok {
		return errors.New("invalid event type")
	}

	log.Printf(
		"RatingCreatedEvent handled for ID: %s",
		e.RatingID.String(),
	)

	return nil
}

func (h RatingEventHandler) HandleRatingChanged(event ddd.DomainEvent) error {
	e, ok := event.(*domain.RatingChangedEvent)
	if !ok {
		return errors.New("invalid event type")
	}

	log.Printf(
		"RatingChangedEvent handled for ID: %s",
		e.RatingID.String(),
	)

	return nil
}

func (h RatingEventHandler) HandleRatingRemoved(event ddd.DomainEvent) error {
	e, ok := event.(*domain.RatingRemovedEvent)
	if !ok {
		return errors.New("invalid event type")
	}

	log.Printf(
		"RatingRemovedEvent handled for ID: %s",
		e.RatingID.String(),
	)

	return nil
}
