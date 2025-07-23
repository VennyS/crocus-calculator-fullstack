package repository

import (
	"crocus-calculator/internal/dto"
	"crocus-calculator/internal/storage"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var (
	ErrLoginAlreadyExists = errors.New("user with this login already exists")
	ErrInvalidCredentials = errors.New("invalid login or password")
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
    `

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}

	_, err = stmt.Exec(user)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrLoginAlreadyExists
		}
		return err
	}

	return nil
}

func (r authPostgresRepo) Login(credentials *dto.LoginRequest) (int64, string, error) {
	query := `
        SELECT id, password 
        FROM users 
        WHERE login = :login AND deleted_at IS NULL
    `

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return 0, "", err
	}
	defer stmt.Close()

	var result struct {
		ID       int64
		Password string
	}

	err = stmt.Get(&result, credentials)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", ErrInvalidCredentials
		}
		return 0, "", err
	}

	return result.ID, result.Password, nil
}
