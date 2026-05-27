package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound            = errors.New("resource not found")
	ErrInternalServerError = errors.New("something went wrong, please try again later")
)

type Storage struct {
	Post interface {
		Create(ctx context.Context, post *Post) error
		GetById(ctx context.Context, postId int64) (*Post, error)
	}
	User interface {
		Create(ctx context.Context, user *User) error
	}
	Comment interface {
		GetByPostId(ctx context.Context, postId int64) ([]Comment, error)
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Post:    &PostStore{db: db},
		User:    &UserStore{db: db},
		Comment: &CommentStore{db: db},
	}
}
