package domain

type UserRepository interface {
	All() ([]User, error)
	FindByID(id UserID) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByUsername(username string) (*User, error)
	Exists(id UserID) (bool, error)
	UsernameExists(username string) (bool, error)
	EmailExists(email string) (bool, error)
	Create(user *User) (*User, error)
	UpdateRoles(id UserID, roles []UserRole) error
	UpdateDescription(id UserID, newDescription string) error
	UpdatePasswordHash(id UserID, passwordHash string) error
}
