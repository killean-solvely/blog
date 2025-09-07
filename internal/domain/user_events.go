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
	UserPasswordUpdatedEventType    EventType = "UserPasswordUpdated"
)

type UserCreatedEvent struct {
	UserID       UserID
	Email        string
	PasswordHash string
	Username     string
	Description  string
	UserRoles    []UserRole
	JoinDate     time.Time
	occurredOn   time.Time
}

func NewUserCreatedEvent(
	id UserID,
	email, username, passwordHash, description string,
	userRoles []UserRole,
	joinDate time.Time,
) *UserCreatedEvent {
	return &UserCreatedEvent{
		UserID:       id,
		Email:        email,
		PasswordHash: passwordHash,
		Username:     username,
		Description:  description,
		UserRoles:    userRoles,
		JoinDate:     joinDate,
		occurredOn:   time.Now(),
	}
}

func (e UserCreatedEvent) OccurredOn() time.Time { return e.occurredOn }
func (e UserCreatedEvent) EventType() string     { return string(UserCreatedEventType) }

type UserRoleAddedEvent struct {
	UserID     UserID
	Role       UserRole
	occurredOn time.Time
}

func NewUserRoleAddedEvent(id UserID, role UserRole) *UserRoleAddedEvent {
	return &UserRoleAddedEvent{
		UserID:     id,
		Role:       role,
		occurredOn: time.Now(),
	}
}

func (e UserRoleAddedEvent) OccurredOn() time.Time { return e.occurredOn }
func (e UserRoleAddedEvent) EventType() string     { return string(UserRoleAddedEventType) }

type UserRoleRemovedEvent struct {
	UserID     UserID
	Role       UserRole
	occurredOn time.Time
}

func NewUserRoleRemovedEvent(id UserID, role UserRole) *UserRoleRemovedEvent {
	return &UserRoleRemovedEvent{
		UserID:     id,
		Role:       role,
		occurredOn: time.Now(),
	}
}

func (e UserRoleRemovedEvent) OccurredOn() time.Time { return e.occurredOn }
func (e UserRoleRemovedEvent) EventType() string     { return string(UserRoleRemovedEventType) }

type UserDescriptionUpdatedEvent struct {
	UserID      UserID
	Description string
	occurredOn  time.Time
}

func NewUserDescriptionUpdatedEvent(id UserID, description string) *UserDescriptionUpdatedEvent {
	return &UserDescriptionUpdatedEvent{
		UserID:      id,
		Description: description,
		occurredOn:  time.Now(),
	}
}

func (e UserDescriptionUpdatedEvent) OccurredOn() time.Time { return e.occurredOn }

func (e UserDescriptionUpdatedEvent) EventType() string {
	return string(UserDescriptionUpdatedEventType)
}

type UserPasswordUpdatedEvent struct {
	UserID     UserID
	Password   string
	occurredOn time.Time
}

func NewUserPasswordUpdatedEvent(id UserID, description string) *UserPasswordUpdatedEvent {
	return &UserPasswordUpdatedEvent{
		UserID:     id,
		Password:   description,
		occurredOn: time.Now(),
	}
}

func (e UserPasswordUpdatedEvent) OccurredOn() time.Time { return e.occurredOn }

func (e UserPasswordUpdatedEvent) EventType() string { return string(UserPasswordUpdatedEventType) }

func init() {
	ddd.EventRegistry.Register(
		UserCreatedEvent{},
		"Raised when a new user is created",
	)

	ddd.EventRegistry.Register(
		UserRoleAddedEvent{},
		"Raised when a new role is added to a user",
	)

	ddd.EventRegistry.Register(
		UserRoleRemovedEvent{},
		"Raised when a role is removed from a user",
	)

	ddd.EventRegistry.Register(
		UserDescriptionUpdatedEvent{},
		"Raised when a user's description is updated",
	)

	ddd.EventRegistry.Register(
		UserPasswordUpdatedEvent{},
		"Raised when a user's password is updated",
	)
}
