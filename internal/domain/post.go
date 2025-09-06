package domain

import (
	"time"

	"blog/pkg/ddd"

	"github.com/google/uuid"
)

type Post struct {
	*ddd.AggregateBase
	authorID     UserID
	title        string
	content      string
	createdAt    time.Time
	lastEditedAt *time.Time
	archivedAt   *time.Time
}

func NewPost(authorID UserID, title string, content string) (*Post, error) {
	if title == "" {
		return nil, ErrTitleCannotBeEmpty
	}

	if content == "" {
		return nil, ErrContentCannotBeEmpty
	}

	now := time.Now()

	post := &Post{
		AggregateBase: &ddd.AggregateBase{},
		title:         title,
		content:       content,
		createdAt:     now,
		lastEditedAt:  nil,
		archivedAt:    nil,
	}

	newID := NewPostID(uuid.New().String())
	post.SetID(newID)

	event := NewPostCreatedEvent(post.GetID(), title, content, now, nil, nil)
	post.RecordEvent(event)

	return post, nil
}

func (a Post) GetID() PostID {
	return PostID(a.AggregateBase.GetID())
}

func (a *Post) SetID(id PostID) {
	if id == "" {
		return
	}
	a.AggregateBase.SetID(string(id))
}

func (a Post) AuthorID() UserID         { return a.authorID }
func (a Post) Title() string            { return a.title }
func (a Post) Content() string          { return a.content }
func (a Post) CreatedAt() time.Time     { return a.createdAt }
func (a Post) LastEditedAt() *time.Time { return a.lastEditedAt }
func (a Post) Archived() bool           { return a.archivedAt != nil }
func (a Post) ArchivedAt() *time.Time   { return a.archivedAt }

func (a *Post) EditTitle(title string) error {
	if title == "" {
		return ErrTitleCannotBeEmpty
	}

	a.title = title

	event := NewPostTitleEditedEvent(a.GetID(), title)
	a.RecordEvent(event)

	return nil
}

func (a *Post) EditContent(content string) error {
	if content == "" {
		return ErrContentCannotBeEmpty
	}

	a.content = content

	event := NewPostContentEditedEvent(a.GetID(), content)
	a.RecordEvent(event)

	return nil
}

func (a *Post) Archive() {
	now := time.Now()
	a.archivedAt = &now

	event := NewPostArchivedEvent(a.GetID(), now)
	a.RecordEvent(event)
}

func RebuildPost(
	id PostID,
	authorID UserID,
	title string,
	content string,
	createdAt time.Time,
	lastEditedAt *time.Time,
	archivedAt *time.Time,
) *Post {
	post := &Post{
		AggregateBase: &ddd.AggregateBase{},
		authorID:      authorID,
		title:         title,
		content:       content,
		createdAt:     createdAt,
		lastEditedAt:  lastEditedAt,
		archivedAt:    archivedAt,
	}
	post.SetID(id)
	return post
}
