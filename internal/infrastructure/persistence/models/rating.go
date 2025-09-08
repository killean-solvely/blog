package models

import "time"

type Rating struct {
	ID         string     `db:"id"`
	PostID     string     `db:"post_id"`
	UserID     string     `db:"user_id"`
	RatingType string     `db:"rating_type"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
}
