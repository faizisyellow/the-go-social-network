package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("a user with this email already exists")
)

type User struct {
	ID        int          `json:"id"`
	Username  string       `json:"username"`
	Email     string       `json:"email"`
	Password  HashPassword `json:"-"`
	CreatedAt string       `json:"created_at"`
}

type HashPassword struct {
	Text *string
	Hash []byte
}

func (h *HashPassword) Set(text string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	h.Text = &text
	h.Hash = hashed

	return nil
}

type UsersStore struct {
	db *sql.DB
}

func (u *UsersStore) Create(ctx context.Context, tx *sql.Tx, payload *User) error {
	qry := `INSERT INTO users (username,password,email) VALUES(?,?,?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := tx.ExecContext(ctx, qry, &payload.Username, &payload.Password.Hash, &payload.Email)

	//TODO: fix the handling error duplicate key
	if err != nil {
		duplicateKey := "Error 1062"
		switch {
		case strings.Contains(err.Error(), duplicateKey):
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	rqry := `SELECT id,created_at FROM users WHERE id = ?`
	row := tx.QueryRow(rqry, id)

	err = row.Scan(&payload.ID, &payload.CreatedAt)
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

func (u *UsersStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userID int) error {
	query := `INSERT INTO user_invitations(token,user_id,expire) VALUES(?,?,?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(exp))
	if err != nil {
		return err
	}

	return nil
}

func (u *UsersStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {

	return withTx(u.db, ctx, func(tx *sql.Tx) error {
		// transcations

		if err := u.Create(ctx, tx, user); err != nil {
			return err
		}

		if err := u.createUserInvitation(ctx, tx, token, invitationExp, user.ID); err != nil {
			return err
		}

		return nil
	})
}
