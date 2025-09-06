package domain

import (
	"time"

	"blog/pkg/ddd"
)

const (
	CommentCreatedEventType  EventType = "CommentCreated"
	CommentEditedEventType   EventType = "CommentEdited"
	CommentArchivedEventType EventType = "CommentArchived"
)

type CommentCreatedEvent struct {
	CommentID     CommentID
	PostID        PostID
	CommenterID   UserID
	Content       string
	CreatedAt     time.Time
	LastUpdatedAt *time.Time
	ArchivedAt    *time.Time
	occurredOn    time.Time
}

func NewCommentCreatedEvent(
	id CommentID,
	postID PostID,
	commenterID UserID,
	content string,
	createdAt time.Time,
	lastUpdatedAt *time.Time,
	archivedAt *time.Time,
) *CommentCreatedEvent {
	return &CommentCreatedEvent{
		CommentID:     id,
		PostID:        postID,
		CommenterID:   commenterID,
		Content:       content,
		CreatedAt:     createdAt,
		LastUpdatedAt: lastUpdatedAt,
		ArchivedAt:    archivedAt,
		occurredOn:    time.Now(),
	}
}

func (e CommentCreatedEvent) OccurredOn() time.Time { return e.occurredOn }
func (e CommentCreatedEvent) EventType() string     { return string(CommentCreatedEventType) }

type CommentEditedEvent struct {
	CommentID     CommentID
	Content       string
	LastUpdatedAt time.Time
	occurredOn    time.Time
}

func NewCommentEditedEvent(
	commentID CommentID,
	content string,
	lastUpdatedAt time.Time,
) *CommentEditedEvent {
	return &CommentEditedEvent{
		CommentID:     commentID,
		Content:       content,
		LastUpdatedAt: lastUpdatedAt,
		occurredOn:    time.Now(),
	}
}

func (e CommentEditedEvent) OccurredOn() time.Time { return e.occurredOn }
func (e CommentEditedEvent) EventType() string     { return string(CommentEditedEventType) }

type CommentArchivedEvent struct {
	CommentID  CommentID
	ArchivedAt time.Time
	occurredOn time.Time
}

func NewCommentArchivedEvent(
	commentID CommentID,
	archivedAt time.Time,
) *CommentArchivedEvent {
	return &CommentArchivedEvent{
		CommentID:  commentID,
		ArchivedAt: archivedAt,
		occurredOn: time.Now(),
	}
}

func (e CommentArchivedEvent) OccurredOn() time.Time { return e.occurredOn }
func (e CommentArchivedEvent) EventType() string     { return string(CommentArchivedEventType) }

func init() {
	ddd.EventRegistry.Register(
		CommentCreatedEvent{},
		"Raised when a new comment is created",
	)

	ddd.EventRegistry.Register(
		CommentEditedEvent{},
		"Raised when a comment is edited",
	)

	ddd.EventRegistry.Register(
		CommentArchivedEvent{},
		"Raised when a comment is archived",
	)
}
