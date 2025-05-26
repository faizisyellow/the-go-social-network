package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("records not found")
)

type Storage struct {
	Posts interface {
		GetPostByID(context.Context, int) (*Post, error)
		Create(context.Context, *Post) error
	}

	Users interface {
		Create(context.Context, *User) error
	}
}

func NewStorage(db *sql.DB) Storage {

	return Storage{
		Posts: &PostStore{db},
		Users: &UsersStore{db},
	}
}
