package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type Post struct {
	Id        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	Tags      []string  `json:"tags"`
	UserId    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Comments []Comment `json:"comments"`
}
type PostStore struct {
	db *sql.DB
}

func (p *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, tags, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	err := p.db.
		QueryRowContext(ctx, query, post.Content, post.Title, pq.Array(post.Tags), post.UserId).
		Scan(&post.Id, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostStore) GetById(ctx context.Context, postId int64) (*Post, error) {
	query := `
		SELECT id, content, title, tags, user_id, created_at, updated_at
		FROM posts
		WHERE id = $1
		LIMIT 1
	`
	var post Post
	err := p.db.
		QueryRowContext(ctx, query, postId).
		Scan(&post.Id, &post.Content, &post.Title, pq.Array(&post.Tags), &post.UserId, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &post, nil
}
