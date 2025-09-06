package domain

import (
	"time"

	"blog/pkg/ddd"

	"github.com/google/uuid"
)

type Comment struct {
	*ddd.AggregateBase
	postID        PostID
	commenterID   UserID
	content       string
	createdAt     time.Time
	lastUpdatedAt *time.Time
	archivedAt    *time.Time
}

func NewComment(postID PostID, commenterID UserID, content string) (*Comment, error) {
	if content == "" {
		return nil, ErrCommentCannotBeEmpty
	}

	now := time.Now()

	comment := &Comment{
		AggregateBase: &ddd.AggregateBase{},
		postID:        postID,
		commenterID:   commenterID,
		content:       content,
		createdAt:     now,
		lastUpdatedAt: nil,
		archivedAt:    nil,
	}

	newID := NewCommentID(uuid.New().String())
	comment.SetID(newID)

	event := NewCommentCreatedEvent(comment.GetID(), postID, commenterID, content, now, nil, nil)
	comment.RecordEvent(event)

	return comment, nil
}

func (a Comment) GetID() CommentID {
	return CommentID(a.AggregateBase.GetID())
}

func (a *Comment) SetID(id CommentID) {
	if id == "" {
		return
	}
	a.AggregateBase.SetID(string(id))
}

func (a Comment) PostID() PostID            { return a.postID }
func (a Comment) CommenterID() UserID       { return a.commenterID }
func (a Comment) Content() string           { return a.content }
func (a Comment) CreatedAt() time.Time      { return a.createdAt }
func (a Comment) LastUpdatedAt() *time.Time { return a.lastUpdatedAt }
func (a Comment) Archived() bool            { return a.archivedAt != nil }
func (a Comment) ArchivedAt() *time.Time    { return a.archivedAt }

func (a *Comment) Edit(content string) error {
	if content == "" {
		return ErrCommentCannotBeEmpty
	}

	now := time.Now()
	a.content = content
	a.lastUpdatedAt = &now

	event := NewCommentEditedEvent(a.GetID(), content, now)
	a.RecordEvent(event)

	return nil
}

func (a *Comment) Archive() {
	now := time.Now()
	a.archivedAt = &now

	event := NewCommentArchivedEvent(a.GetID(), now)
	a.RecordEvent(event)
}

func RebuildComment(
	id CommentID,
	postID PostID,
	commenterID UserID,
	content string,
	createdAt time.Time,
	lastUpdatedAt *time.Time,
	archivedAt *time.Time,
) *Comment {
	comment := &Comment{
		postID:        postID,
		commenterID:   commenterID,
		content:       content,
		createdAt:     createdAt,
		lastUpdatedAt: lastUpdatedAt,
		archivedAt:    archivedAt,
	}
	comment.SetID(id)
	return comment
}
