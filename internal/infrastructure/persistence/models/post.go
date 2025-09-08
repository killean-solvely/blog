package models

import "time"

type Post struct {
	ID           string     `db:"id"`
	AuthorID     string     `db:"author_id"`
	Title        string     `db:"title"`
	Content      string     `db:"content"`
	CreatedAt    time.Time  `db:"created_at"`
	LastEditedAt *time.Time `db:"last_edited_at"`
	ArchivedAt   *time.Time `db:"archived_at"`
}
