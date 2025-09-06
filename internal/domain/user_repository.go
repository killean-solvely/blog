package domain

type UserRepository interface {
	All() ([]User, error)
	FindByID(id UserID) (*User, error)
	Create(user *User) (*User, error)
	UpdateRoles(roles []UserRole) error
	UpdateDescription(newDescription string) error
}
