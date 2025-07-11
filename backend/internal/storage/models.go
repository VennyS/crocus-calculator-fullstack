package storage

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64        `db:"id"`
	Login     string       `db:"login" json:"login"`
	Password  string       `db:"password" json:"-"`
	Permit    string       `db:"permit" json:"permit"`
	Fullname  string       `db:"full_name" json:"full_name"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt time.Time    `db:"updated_at" json:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at" json:"deleted_at,omitempty"`
}
