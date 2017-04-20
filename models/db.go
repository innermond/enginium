package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDatabase(ds string) error {

	if DB != nil {
		return nil
	}

	var err error

	DB, err = sql.Open("mysql", ds)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	return nil
}
