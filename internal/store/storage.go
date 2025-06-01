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
		GetUserFeed(context.Context, int, PaginatedFeedQuery) ([]PostWithMetaData, error)
	}

	Users interface {
		Create(context.Context, *sql.Tx, *User) error
		GetByID(context.Context, int) (*User, error)
		CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error
		Activate(context.Context, string) error
		Delete(context.Context, int) error
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

func withTx(db *sql.DB, ctx context.Context, fnc func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fnc(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
