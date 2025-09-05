package domain

import (
	"time"

	"blog/pkg/ddd"

	"github.com/google/uuid"
)

type Rating struct {
	*ddd.AggregateBase
	postID     PostID
	userID     UserID
	ratingType RatingType
	createdAt  time.Time
	updatedAt  *time.Time
}

func NewRating(postID PostID, userID UserID, ratingType RatingType) *Rating {
	now := time.Now()

	rating := &Rating{
		AggregateBase: &ddd.AggregateBase{},
		postID:        postID,
		userID:        userID,
		ratingType:    ratingType,
		createdAt:     now,
		updatedAt:     nil,
	}

	newID := NewRatingID(uuid.New().String())
	rating.SetID(newID)

	event := NewRatingCreatedEvent(rating.GetID(), postID, userID, ratingType, now, nil)
	rating.RecordEvent(event)

	return rating
}

func (a Rating) GetID() RatingID {
	return RatingID(a.AggregateBase.GetID())
}

func (a *Rating) SetID(id RatingID) {
	if id == "" {
		return
	}
	a.AggregateBase.SetID(string(id))
}

func (a *Rating) ChangeRating(ratingType RatingType) {
	now := time.Now()
	a.ratingType = ratingType
	a.updatedAt = &now

	event := NewRatingChangedEvent(a.GetID(), ratingType, now)
	a.RecordEvent(event)
}

func (a *Rating) RemoveRating() {
	event := NewRatingRemovedEvent(a.GetID())
	a.RecordEvent(event)
}
