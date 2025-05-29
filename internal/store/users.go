package store

import (
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int      `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hashed

	return nil
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

func (u *UsersStore) GetByID(ctx context.Context, userId int) (*User, error) {
	query := `SELECT id, username, email, password, created_at FROM users WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res := u.db.QueryRowContext(ctx, query, userId)

	user := User{}
	err := res.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
