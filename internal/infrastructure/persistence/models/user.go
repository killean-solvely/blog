package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Username     string    `json:"username"`
	Description  string    `json:"description"`
	UserRoles    string    `json:"user_roles"`
	JoinDate     time.Time `json:"join_date"`
}
