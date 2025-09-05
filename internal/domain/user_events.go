package domain

import (
	"time"

	"blog/pkg/ddd"
)

const (
	UserCreatedEventType            EventType = "UserCreated"
	UserRoleAddedEventType          EventType = "UserRoleAdded"
	UserRoleRemovedEventType        EventType = "UserRoleRemoved"
	UserDescriptionUpdatedEventType EventType = "UserDescriptionUpdated"
)

type UserCreatedEvent struct {
	UserID      UserID
	Email       string
	Username    string
	Description string
	UserRoles   []UserRole
	JoinDate    time.Time
	occurredOn  time.Time
}

func NewUserCreatedEvent(
	id UserID,
	email, username, description string,
	userRoles []UserRole,
	joinDate time.Time,
) *UserCreatedEvent {
	return &UserCreatedEvent{
		UserID:      id,
		Email:       email,
		Username:    username,
		Description: description,
		UserRoles:   userRoles,
		JoinDate:    joinDate,
		occurredOn:  time.Now(),
	}
}

func (e UserCreatedEvent) OccurredOn() time.Time { return e.occurredOn }
func (e UserCreatedEvent) EventType() string     { return string(UserCreatedEventType) }

type UserRoleAdded struct {
	UserID     UserID
	Role       UserRole
	occurredOn time.Time
}

func NewUserRoleAddedEvent(id UserID, role UserRole) *UserRoleAdded {
	return &UserRoleAdded{
		UserID:     id,
		Role:       role,
		occurredOn: time.Now(),
	}
}

func (e UserRoleAdded) OccurredOn() time.Time { return e.occurredOn }
func (e UserRoleAdded) EventType() string     { return string(UserRoleAddedEventType) }

type UserRoleRemoved struct {
	UserID     UserID
	Role       UserRole
	occurredOn time.Time
}

func NewUserRoleRemovedEvent(id UserID, role UserRole) *UserRoleRemoved {
	return &UserRoleRemoved{
		UserID:     id,
		Role:       role,
		occurredOn: time.Now(),
	}
}

func (e UserRoleRemoved) OccurredOn() time.Time { return e.occurredOn }
func (e UserRoleRemoved) EventType() string     { return string(UserRoleRemovedEventType) }

type UserDescriptionUpdated struct {
	UserID      UserID
	Description string
	occurredOn  time.Time
}

func NewUserDescriptionUpdatedEvent(id UserID, description string) *UserDescriptionUpdated {
	return &UserDescriptionUpdated{
		UserID:      id,
		Description: description,
		occurredOn:  time.Now(),
	}
}

func (e UserDescriptionUpdated) OccurredOn() time.Time { return e.occurredOn }

func (e UserDescriptionUpdated) EventType() string { return string(UserDescriptionUpdatedEventType) }

func init() {
	ddd.EventRegistry.Register(
		UserCreatedEvent{},
		"Raised when a new user is created",
	)

	ddd.EventRegistry.Register(
		UserRoleAdded{},
		"Raised when a new role is added to a user",
	)

	ddd.EventRegistry.Register(
		UserRoleRemoved{},
		"Raised when a role is removed from a user",
	)

	ddd.EventRegistry.Register(
		UserDescriptionUpdated{},
		"Raised when a user's description is updated",
	)
}
