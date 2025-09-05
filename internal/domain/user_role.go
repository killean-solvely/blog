package domain

type UserRole string

const (
	UserRoleAuthor    UserRole = "AUTHOR"
	UserRoleCommenter UserRole = "COMMENTER"
	UserRoleAdmin     UserRole = "ADMIN"
)

func (ur UserRole) String() string {
	return string(ur)
}
