package domain

import (
	"time"

	"blog/pkg/ddd"
)

type EventType string

func (e EventType) String() string {
	return string(e)
}

const (
	UserCreatedEventType EventType = "UserCreated"
)

type UserCreatedEvent struct {
	UserID      UserID
	Email       string
	Username    string
	Description string
	JoinDate    time.Time
	occurredOn  time.Time
}

func NewUserCreatedEvent(
	id UserID,
	email, username, description string,
	joinDate time.Time,
) *UserCreatedEvent {
	return &UserCreatedEvent{
		UserID:      id,
		Email:       email,
		Username:    username,
		Description: description,
		JoinDate:    joinDate,
		occurredOn:  time.Now(),
	}
}

func (e UserCreatedEvent) OccurredOn() time.Time { return e.occurredOn }
func (e UserCreatedEvent) EventType() string     { return string(UserCreatedEventType) }

func init() {
	ddd.EventRegistry.Register(
		UserCreatedEvent{},
		"Raised when a new user is created",
	)
}
