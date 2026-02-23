package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewPostgresConnection(databaseUrl string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
