package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func NewDB() (*sql.DB, error) {
	connStr := "user=postgres dbname=bot password=1 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, errors.Wrap(err, "error opening DB connection")
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "error pinging DB")
	}

	return db, nil
}
