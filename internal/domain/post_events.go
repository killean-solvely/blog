package domain

import (
	"time"

	"blog/pkg/ddd"
)

const (
	PostCreatedEventType       EventType = "PostCreated"
	PostTitleEditedEventType   EventType = "PostTitleEdited"
	PostContentEditedEventType EventType = "PostContentEdited"
	PostArchivedEventType      EventType = "PostArchived"
)

type PostCreatedEvent struct {
	PostID       PostID
	Title        string
	Content      string
	CreatedAt    time.Time
	LastEditedAt *time.Time
	ArchivedAt   *time.Time
	occurredOn   time.Time
}

func NewPostCreatedEvent(
	id PostID,
	title, content string,
	created time.Time,
	lastEdited *time.Time,
	archivedAt *time.Time,
) *PostCreatedEvent {
	return &PostCreatedEvent{
		PostID:       id,
		Title:        title,
		Content:      content,
		CreatedAt:    created,
		LastEditedAt: lastEdited,
		ArchivedAt:   archivedAt,
		occurredOn:   time.Now(),
	}
}

func (e PostCreatedEvent) OccurredOn() time.Time { return e.occurredOn }
func (e PostCreatedEvent) EventType() string     { return string(PostCreatedEventType) }

type PostTitleEditedEvent struct {
	PostID     PostID
	NewTitle   string
	occurredOn time.Time
}

func NewPostTitleEditedEvent(id PostID, newTitle string) *PostTitleEditedEvent {
	return &PostTitleEditedEvent{
		PostID:     id,
		NewTitle:   newTitle,
		occurredOn: time.Now(),
	}
}

func (e PostTitleEditedEvent) OccurredOn() time.Time { return e.occurredOn }
func (e PostTitleEditedEvent) EventType() string     { return string(PostTitleEditedEventType) }

type PostContentEditedEvent struct {
	PostID     PostID
	NewContent string
	occurredOn time.Time
}

func NewPostContentEditedEvent(id PostID, newContent string) *PostContentEditedEvent {
	return &PostContentEditedEvent{
		PostID:     id,
		NewContent: newContent,
		occurredOn: time.Now(),
	}
}

func (e PostContentEditedEvent) OccurredOn() time.Time { return e.occurredOn }
func (e PostContentEditedEvent) EventType() string     { return string(PostContentEditedEventType) }

type PostArchivedEvent struct {
	PostID     PostID
	ArchivedAt time.Time
	occurredOn time.Time
}

func NewPostArchivedEvent(id PostID, archivedAt time.Time) *PostArchivedEvent {
	return &PostArchivedEvent{
		PostID:     id,
		ArchivedAt: archivedAt,
		occurredOn: time.Now(),
	}
}

func (e PostArchivedEvent) OccurredOn() time.Time { return e.occurredOn }
func (e PostArchivedEvent) EventType() string     { return string(PostArchivedEventType) }

func init() {
	ddd.EventRegistry.Register(
		PostCreatedEvent{},
		"Raised when a new post is created",
	)

	ddd.EventRegistry.Register(
		PostTitleEditedEvent{},
		"Raised when a post's title is edited",
	)

	ddd.EventRegistry.Register(
		PostContentEditedEvent{},
		"Raised when a post's content is edited",
	)

	ddd.EventRegistry.Register(
		PostArchivedEvent{},
		"Raised when a post is archived",
	)
}
