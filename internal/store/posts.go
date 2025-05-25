package store

import (
	"context"
	"database/sql"
	"errors"
)

type Post struct {
	ID        int      `json:"id"`
	Content   string   `json:"content"`
	Title     string   `json:"title"`
	UserID    int      `json:"user_id"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_td"`
}

type PostStore struct {
	db *sql.DB
}

func (p *PostStore) Create(ctx context.Context, payload *Post) error {
	qry := `INSERT INTO posts(content,title,user_id,tags) VALUES(?,?,?,?)`

	p.db.ExecContext(ctx, qry, payload.Content, payload.Title, payload.UserID, payload.Tags)

	rqry := `SELECT id,created_at,updated_at FROM posts WHERE id=(SELECT LAST_INSERT_ID)`

	row := p.db.QueryRow(rqry)
	err := row.Scan(&payload.ID, &payload.CreatedAt, &payload.UpdatedAt)

	if err == sql.ErrNoRows {
		return errors.New("empty row")
	} else if err != nil {
		return err
	}

	return nil
}
