package store

import (
	"context"
	"database/sql"
	"strings"
)

type Follower struct {
	UserID     int    `json:"user_id"`
	FollowerID int    `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowersStore struct {
	db *sql.DB
}

func (f *FollowersStore) Follow(ctx context.Context, toFollowUserID, userID int) error {
	query := `INSERT INTO followers(user_id,follower_id) VALUES(?,?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := f.db.ExecContext(ctx, query, &toFollowUserID, &userID)
	if err != nil {
		duplicateKey := "Error 1062"

		if strings.Contains(err.Error(), duplicateKey) {
			return ErrConflict
		} else {
			return err
		}
	}

	return nil
}

// TODO: for unfollow user if possible just create a followed column and then just toggle it (update true or false)
func (f *FollowersStore) UnFollow(ctx context.Context, toUnFollowUserID, userID int) error {
	query := `DELETE FROM followers WHERE user_id = ? AND follower_id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := f.db.ExecContext(ctx, query, &toUnFollowUserID, &userID)
	if err != nil {
		return err
	}

	return nil

}
