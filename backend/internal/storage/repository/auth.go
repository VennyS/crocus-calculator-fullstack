package repository

import (
	"crocus-calculator/internal/storage"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var (
	ErrLoginAlreadyExists = errors.New("user with this login already exists")
)

type authPostgresRepo struct {
	db *sqlx.DB
}

func NewAuthPostgresRepo(db *sqlx.DB) authPostgresRepo {
	return authPostgresRepo{db: db}
}

func (r authPostgresRepo) Create(user *storage.User) error {
	query := `
		INSERT INTO users (login, password, permit, full_name, created_at, deleted_at)
		VALUES (:login, :password, :permit, :full_name, :created_at, :deleted_at)
		RETURNING id
	`

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}

	err = stmt.Get(&user.ID, user)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrLoginAlreadyExists
		}
		return err
	}

	return nil
}
