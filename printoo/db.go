package printoo

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct{ *sql.DB }

type Tx struct{ *sql.Tx }

func Open(dns string) (*DB, error) {
	db, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) Begin() (*Tx, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}

	return &Tx{tx}, nil
}
