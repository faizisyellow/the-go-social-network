package store

import (
	"context"
	"database/sql"
	"strings"
)

type Follower struct {
	FollowedID int    `json:"followed_id"`
	FollowerID int    `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type UserFollows struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type FollowersAndFollowingCount struct {
	Following int `json:"following"`
	Followers int `json:"followers"`
}

type FollowersStore struct {
	db *sql.DB
}

func (f *FollowersStore) Follow(ctx context.Context, toFollowUserID, userID int) error {
	query := `INSERT INTO followers(followed_id,follower_id) VALUES(?,?)`

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
	query := `DELETE FROM followers WHERE followed_id = ? AND follower_id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := f.db.ExecContext(ctx, query, &toUnFollowUserID, &userID)
	if err != nil {
		return err
	}

	return nil

}

func (f *FollowersStore) TotalFollowersAndFollowing(ctx context.Context, userID int) (*FollowersAndFollowingCount, error) {
	query := `
	SELECT following, follower FROM users
	JOIN (SELECT followed_id, COUNT(*) AS follower FROM followers GROUP BY followed_id) f1 ON users.id = f1.followed_id
	JOIN (SELECT follower_id, COUNT(*) AS following FROM followers GROUP BY follower_id) f2 ON users.id = f2.follower_id
	WHERE users.id = ?
	`

	cff := &FollowersAndFollowingCount{}
	err := f.db.QueryRowContext(ctx, query, userID).Scan(&cff.Following, &cff.Followers)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return cff, nil
}

func (f *FollowersStore) GetUserFollowers(ctx context.Context, userID int, followers *[]*UserFollows) error {
	query := `SELECT id, username FROM users JOIN followers ON users.id = followers.follower_id AND followers.followed_id = ? `

	rows, err := f.db.QueryContext(ctx, query, userID)
	if err != nil {
		return err
	}

	for rows.Next() {
		var user UserFollows

		err := rows.Scan(&user.ID, &user.Username)
		if err != nil {
			return err
		}

		*followers = append(*followers, &user)
	}

	return nil
}

func (f *FollowersStore) GetUserFollowing(ctx context.Context, userID int, following *[]*UserFollows) error {
	query := `SELECT id, username FROM users JOIN followers ON users.id = followers.followed_id AND followers.follower_id = ? `

	rows, err := f.db.QueryContext(ctx, query, userID)
	if err != nil {
		return err
	}

	for rows.Next() {
		var user UserFollows

		err := rows.Scan(&user.ID, &user.Username)
		if err != nil {
			return err
		}

		*following = append(*following, &user)
	}

	return nil
}
