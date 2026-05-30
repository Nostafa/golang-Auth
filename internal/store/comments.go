package store

import (
	"context"
	"database/sql"
	"time"
)

type Comment struct {
	Id        int64     `json:"id"`
	UserId    int64     `json:"user_id"`
	PostId    int64     `json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `json:"user"`
}

type CommentStore struct {
	db *sql.DB
}

func (c *CommentStore) Create(ctx context.Context, comment *Comment) error {
	query := `
		INSERT INTO comments (user_id, post_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	err := c.db.QueryRowContext(ctx, query, comment.UserId, comment.PostId, comment.Content).Scan(&comment.Id, &comment.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommentStore) GetByPostId(ctx context.Context, postId int64) ([]Comment, error) {
	query := `
		SELECT
			c.id,
			c.user_id,
			c.post_id,
			c.content,
			c.created_at,
			u.id,
			u.first_name,
			u.last_name,
			u.username,
			u.email
		FROM comments c
		INNER JOIN users u ON u.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at ASC
	`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := c.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.Id,
			&comment.UserId,
			&comment.PostId,
			&comment.Content,
			&comment.CreatedAt,
			&comment.User.Id,
			&comment.User.FirstName,
			&comment.User.LastName,
			&comment.User.Username,
			&comment.User.Email,
		)

		if err != nil {
			return nil, err
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
