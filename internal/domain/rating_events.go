package domain

import (
	"time"

	"blog/pkg/ddd"
)

const (
	RatingCreatedEventType EventType = "RatingCreated"
	RatingChangedEventType EventType = "RatingChanged"
	RatingRemovedEventType EventType = "RatingRemoved"
)

type RatingCreatedEvent struct {
	RatingID   RatingID
	PostID     PostID
	UserID     UserID
	RatingType RatingType
	CreatedAt  time.Time
	UpdatedAt  *time.Time
	occurredOn time.Time
}

func NewRatingCreatedEvent(
	id RatingID,
	postID PostID,
	userID UserID,
	ratingType RatingType,
	createdAt time.Time,
	updatedAt *time.Time,
) *RatingCreatedEvent {
	return &RatingCreatedEvent{
		RatingID:   id,
		PostID:     postID,
		UserID:     userID,
		RatingType: ratingType,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
		occurredOn: time.Now(),
	}
}

func (e RatingCreatedEvent) OccurredOn() time.Time { return e.occurredOn }
func (e RatingCreatedEvent) EventType() string     { return string(RatingCreatedEventType) }

type RatingChangedEvent struct {
	RatingID      RatingID
	NewRatingType RatingType
	UpdatedAt     time.Time
	occurredOn    time.Time
}

func NewRatingChangedEvent(
	id RatingID,
	newRatingType RatingType,
	updatedAt time.Time,
) *RatingChangedEvent {
	return &RatingChangedEvent{
		RatingID:      id,
		NewRatingType: newRatingType,
		UpdatedAt:     updatedAt,
		occurredOn:    time.Now(),
	}
}

func (e RatingChangedEvent) OccurredOn() time.Time { return e.occurredOn }
func (e RatingChangedEvent) EventType() string     { return string(RatingChangedEventType) }

type RatingRemovedEvent struct {
	RatingID   RatingID
	occurredOn time.Time
}

func NewRatingRemovedEvent(
	id RatingID,
) *RatingRemovedEvent {
	return &RatingRemovedEvent{
		RatingID:   id,
		occurredOn: time.Now(),
	}
}

func (e RatingRemovedEvent) OccurredOn() time.Time { return e.occurredOn }
func (e RatingRemovedEvent) EventType() string     { return string(RatingRemovedEventType) }

func init() {
	ddd.EventRegistry.Register(
		RatingCreatedEvent{},
		"Raised when a new rating is created",
	)

	ddd.EventRegistry.Register(
		RatingChangedEvent{},
		"Raised when a rating is changed",
	)

	ddd.EventRegistry.Register(
		RatingRemovedEvent{},
		"Raised when a rating is removed",
	)
}
