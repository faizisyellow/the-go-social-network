package store

import (
	"context"
	"database/sql"
	"errors"
)

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

type UsersStore struct {
	db *sql.DB
}

func (u *UsersStore) Create(ctx context.Context, payload *User) error {
	qry := `INSERT INTO users (username,password,email) VALUES(?,?,?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	u.db.ExecContext(ctx, qry, payload.Username, payload.Password, payload.Email)

	rqry := `SELECT id,created_at FROM users WHERE id=(SELECT LAST_INSERT_ID)`
	row := u.db.QueryRow(rqry)

	err := row.Scan(&payload.ID, &payload.CreatedAt)
	if err == sql.ErrNoRows {
		return errors.New("empty row")
	} else if err != nil {
		return err
	}

	return nil
}
