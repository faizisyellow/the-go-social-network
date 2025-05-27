package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("records not found")
	QueryTimeoutDuration = 5 * time.Second
)

type Storage struct {
	Posts interface {
		GetPostByID(context.Context, int) (*Post, error)
		Create(context.Context, *Post) error
		Delete(context.Context, int) error
		Update(context.Context, *Post) error
	}

	Users interface {
		Create(context.Context, *User) error
	}

	Comments interface {
		GetPostByID(context.Context, int) ([]Comment, error)
	}
}

func NewStorage(db *sql.DB) Storage {

	return Storage{
		Posts:    &PostStore{db},
		Users:    &UsersStore{db},
		Comments: &CommentsStore{db: db},
	}
}
