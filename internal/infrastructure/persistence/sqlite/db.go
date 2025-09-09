package sqlite

import (
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func NewDB() (*DB, error) {
	db, err := sqlx.Connect("sqlite3", "file:blog.db")
	if err != nil {
		return nil, err
	}

	return &DB{
		DB: db,
	}, nil
}
