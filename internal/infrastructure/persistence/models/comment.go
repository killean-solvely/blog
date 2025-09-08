package models

import "time"

type Comment struct {
	ID            string     `db:"id"`
	PostID        string     `db:"post_id"`
	CommenterID   string     `db:"commenter_id"`
	Content       string     `db:"content"`
	CreatedAt     time.Time  `db:"created_at"`
	LastUpdatedAt *time.Time `db:"last_updated_at"`
	ArchivedAt    *time.Time `db:"archived_at"`
}
