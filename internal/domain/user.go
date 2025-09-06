package domain

import (
	"time"

	"blog/pkg/ddd"

	"github.com/google/uuid"
)

type User struct {
	*ddd.AggregateBase
	email        string
	passwordHash string
	username     string
	description  string
	userRoles    map[UserRole]bool
	joinDate     time.Time
}

func NewUser(
	email, username, passwordHash, description string,
	userRoles []UserRole,
) (*User, error) {
	if len(userRoles) == 0 {
		return nil, ErrMissingUserRoles
	}

	now := time.Now()

	setRoles := map[UserRole]bool{}
	for _, role := range userRoles {
		setRoles[role] = true
	}

	user := &User{
		AggregateBase: &ddd.AggregateBase{},
		email:         email,
		passwordHash:  passwordHash,
		username:      username,
		description:   description,
		userRoles:     setRoles,
		joinDate:      now,
	}

	newID := NewUserID(uuid.New().String())
	user.SetID(newID)

	event := NewUserCreatedEvent(
		user.GetID(),
		email,
		passwordHash,
		username,
		description,
		userRoles,
		now,
	)
	user.RecordEvent(event)

	return user, nil
}

func (a User) GetID() UserID {
	return UserID(a.AggregateBase.GetID())
}

func (a *User) SetID(id UserID) {
	if id == "" {
		return
	}
	a.AggregateBase.SetID(string(id))
}

func (a User) Email() string        { return a.email }
func (a User) Username() string     { return a.username }
func (a User) PasswordHash() string { return a.passwordHash }
func (a User) Description() string  { return a.description }
func (a User) JoinDate() time.Time  { return a.joinDate }

func (a User) UserRoles() []UserRole {
	roleSlice := []UserRole{}
	for k, v := range a.userRoles {
		if v {
			roleSlice = append(roleSlice, k)
		}
	}
	return roleSlice
}

func (a User) CanCreatePost() bool {
	return a.userRoles[UserRoleAuthor]
}

func (a User) CanComment() bool {
	return a.userRoles[UserRoleCommenter]
}

func (a User) IsAdmin() bool {
	return a.userRoles[UserRoleAdmin]
}

func (a *User) AddRole(role UserRole) {
	a.userRoles[role] = true

	event := NewUserRoleAddedEvent(a.GetID(), role)
	a.RecordEvent(event)
}

func (a *User) RemoveRole(role UserRole) {
	delete(a.userRoles, role)

	event := NewUserRoleRemovedEvent(a.GetID(), role)
	a.RecordEvent(event)
}

func (a *User) UpdateDescription(newDescription string) error {
	if len(newDescription) > 255 {
		return ErrDescriptionTooLong
	}
	a.description = newDescription

	event := NewUserDescriptionUpdatedEvent(a.GetID(), newDescription)
	a.RecordEvent(event)

	return nil
}

func (a *User) UpdatePasswordHash(passwordHash string) error {
	a.passwordHash = passwordHash

	event := NewUserPasswordUpdatedEvent(a.GetID(), passwordHash)
	a.RecordEvent(event)

	return nil
}

func RebuildUser(
	id UserID,
	email string,
	passwordHash string,
	username string,
	description string,
	userRoles []UserRole,
	joinDate time.Time,
) *User {
	return nil
}
