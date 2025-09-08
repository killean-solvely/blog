package models

import "time"

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Username     string
	Description  string
	UserRoles    string
	JoinDate     time.Time
}
