package store

import (
	"context"
	"database/sql"
	"errors"
)

type Post struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int       `json:"user_id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comments"`
}

type PostStore struct {
	db *sql.DB
}

func (p *PostStore) Create(ctx context.Context, payload *Post) error {
	qry := `INSERT INTO posts (content, title, user_id) VALUES(?,?,?)`

	result, err := p.db.ExecContext(ctx, qry, payload.Content, payload.Title, payload.UserID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	rqry := `SELECT id, created_at, updated_at FROM posts WHERE id = ?`

	row := p.db.QueryRow(rqry, id)

	err = row.Scan(&payload.ID, &payload.CreatedAt, &payload.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostStore) GetPostByID(ctx context.Context, id int) (*Post, error) {

	qry := `SELECT id, title, content, user_id, created_at, updated_at  FROM posts WHERE id = ?`

	row := p.db.QueryRowContext(ctx, qry, id)

	var post Post

	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt, &post.UpdatedAt)

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

func (p *PostStore) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM posts WHERE id = ?`

	res, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (p *PostStore) Update(ctx context.Context, payload *Post) error {
	query := `UPDATE posts SET title = ?, content = ? WHERE id = ?`

	res, err := p.db.ExecContext(ctx, query, &payload.Title, &payload.Content, &payload.ID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
