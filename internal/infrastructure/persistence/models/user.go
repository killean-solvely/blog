package models

import "time"

type User struct {
	ID           string    `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Username     string    `db:"username"`
	Description  string    `db:"description"`
	UserRoles    string    `db:"user_roles"`
	JoinDate     time.Time `db:"join_date"`
}
