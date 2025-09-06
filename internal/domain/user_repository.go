package domain

type UserRepository interface {
	All() ([]User, error)
	FindByID(id UserID) (*User, error)
	Exists(id UserID) (bool, error)
	Create(user *User) (*User, error)
	UpdateRoles(id UserID, roles []UserRole) error
	UpdateDescription(id UserID, newDescription string) error
}
