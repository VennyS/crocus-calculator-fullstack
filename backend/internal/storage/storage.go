package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func BootstrapDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn) // Connect сразу пингует
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	return db, nil
}
