package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type User struct {
	Id        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserStore struct {
	db *sql.DB
}

func (u *UserStore) Create(ctx context.Context, user *User) (*User, error) {
	query := `
		INSERT INTO users (first_name, last_name, username, email, password)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id,first_name, last_name, username, email, created_at, updated_at
	`
	err := u.db.
		QueryRowContext(ctx, query, user.FirstName, user.LastName, user.Username, user.Email, user.Password).
		Scan(&user.Id, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserStore) GetById(ctx context.Context, userId int64) (*User, error) {
	query := `
		SELECT id, first_name, last_name, username, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`
	var user User
	err := u.db.
		QueryRowContext(ctx, query, userId).
		Scan(&user.Id, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}
