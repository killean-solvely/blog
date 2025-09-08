package sqlite

import (
	"github.com/jmoiron/sqlx"
)

type DB struct {
	db *sqlx.DB
}

func NewDB() (*DB, error) {
	db, err := sqlx.Connect("sqlite3", "file:data.sqlite")
	if err != nil {
		return nil, err
	}

	return &DB{
		db: db,
	}, nil
}
