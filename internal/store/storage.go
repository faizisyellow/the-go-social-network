package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("records not found")
	ErrConflict          = errors.New("resource conflict or resource already exist")
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
		GetByID(context.Context, int) (*User, error)
	}

	Comments interface {
		GetPostByID(context.Context, int) ([]Comment, error)
		Create(ctx context.Context, userId, postID int, content string) error
	}

	Followers interface {
		Follow(ctx context.Context, toFollowUserID, userID int) error
		UnFollow(ctx context.Context, toUnFollowUserID, userID int) error
	}
}

func NewStorage(db *sql.DB) Storage {

	return Storage{
		Posts:     &PostStore{db},
		Users:     &UsersStore{db},
		Comments:  &CommentsStore{db: db},
		Followers: &FollowersStore{db: db},
	}
}
