package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	return sql.Open("sqlite", path)
}

func Query(db *sql.DB, qString string) (*sql.Rows, error) {
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db.Query(qString)
}
