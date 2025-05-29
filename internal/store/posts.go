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
	Version   int       `json:"version"`
	User      User      `json:"user"`
}

type PostWithMetaData struct {
	Post
	CommentCount int `json:"comment_count"`
}

type PostStore struct {
	db *sql.DB
}

func (p *PostStore) Create(ctx context.Context, payload *Post) error {
	qry := `INSERT INTO posts (content, title, user_id) VALUES(?,?,?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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

	qry := `SELECT id, title, content, user_id, version, created_at, updated_at  FROM posts WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	row := p.db.QueryRowContext(ctx, qry, id)

	var post Post

	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.Version, &post.CreatedAt, &post.UpdatedAt)

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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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
	query := `UPDATE posts SET title = ?, content = ?, version = version + 1 WHERE id = ? AND version = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := p.db.ExecContext(ctx, query, &payload.Title, &payload.Content, &payload.ID, &payload.Version)
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

func (p *PostStore) GetUserFeed(ctx context.Context, userId int, fp PaginatedFeedQuery) ([]PostWithMetaData, error) {
	query := `
	SELECT 
		p.id,
		p.title,
		p.content,
		p.created_at,
		u.username,
		COUNT(c.id) AS comment_count
	FROM
		posts p
			LEFT JOIN
		comments c ON c.post_id = p.id
			LEFT JOIN
		users u ON u.id = p.user_id
			JOIN
		followers f ON f.follower_id = p.user_id
			OR p.user_id = ?
	WHERE
        f.user_id = ? OR p.user_id = ?
	GROUP BY (p.id)
	ORDER BY p.created_at ` + fp.Sort + `
	LIMIT ?
	OFFSET ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := p.db.QueryContext(ctx, query, userId, userId, userId, fp.Limit, fp.Offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feeds []PostWithMetaData

	for rows.Next() {
		var post PostWithMetaData
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.User.Username, &post.CommentCount)
		if err != nil {
			return nil, err
		}

		feeds = append(feeds, post)
	}

	return feeds, nil
}
