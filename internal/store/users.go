package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"faizisyellow.github.com/thegosocialnetwork/internal/helpers"
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
	IsActive  bool         `json:"is_active"`
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

func (u *UsersStore) GetUserInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.created_at, u.is_active FROM users u
		JOIN user_invitations ui ON u.id = ui.user_id
		WHERE ui.token = ? AND ui.expire > ? 
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	hashToken, err := helpers.HashToken(token)
	if err != nil {
		return nil, err
	}

	user := &User{}
	err = tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.IsActive)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (u *UsersStore) update(ctx context.Context, tx *sql.Tx, user *User) error {

	query := `UPDATE users SET username = ?, email = ?, is_active = ? WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, &user.Username, &user.Email, &user.IsActive, &user.ID)
	if err != nil {
		return err
	}

	return nil

}

func (u *UsersStore) deleteUserInvitations(ctx context.Context, tx *sql.Tx, userID int) error {
	query := `DELETE FROM user_invitations WHERE user_id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (u *UsersStore) Activate(ctx context.Context, token string) error {

	return withTx(u.db, ctx, func(tx *sql.Tx) error {
		// 1. Find the user that this token invitation belongs to
		user, err := u.GetUserInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		// 2. Update the user
		user.IsActive = true
		if err := u.update(ctx, tx, user); err != nil {
			return err
		}

		// 3. Clean the invitations
		if err := u.deleteUserInvitations(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil
	})
}

// delete user
func (u *UsersStore) delete(ctx context.Context, tx *sql.Tx, userID int) error {
	query := `DELETE FROM users WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	isSuccess, _ := res.RowsAffected()

	if isSuccess == 0 {
		return ErrNotFound
	}

	return nil
}

func (u *UsersStore) Delete(ctx context.Context, userID int) error {
	return withTx(u.db, ctx, func(tx *sql.Tx) error {
		if err := u.delete(ctx, tx, userID); err != nil {
			return err
		}

		if err := u.deleteUserInvitations(ctx, tx, userID); err != nil {
			return err
		}

		return nil
	})
}
