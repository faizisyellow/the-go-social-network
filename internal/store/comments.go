package store

import (
	"context"
	"database/sql"
)

type Comment struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	PostID    int    `json:"post_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}

type CommentsStore struct {
	db *sql.DB
}

func (c *CommentsStore) GetPostByID(ctx context.Context, postID int) ([]Comment, error) {
	query := `SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, users.username, users.id
	 FROM comments c JOIN users ON c.user_id=users.id WHERE c.post_id = ? ORDER BY c.created_at DESC`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := c.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := []Comment{}

	for rows.Next() {
		var c Comment
		c.User = User{}

		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.User.Username, &c.User.ID)
		if err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}

	return comments, nil
}

func (c *CommentsStore) Create(ctx context.Context, userID, postID int, content string) error {
	query := `INSERT INTO comments(user_id,post_id,content) VALUES(?,?,?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := c.db.ExecContext(ctx, query, &userID, &postID, &content)
	if err != nil {
		return err
	}

	return nil
}
