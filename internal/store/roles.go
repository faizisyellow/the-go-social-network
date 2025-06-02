package store

import (
	"context"
	"database/sql"
)

type Role struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Level       int    `json:"level"`
	Description string `json:"description"`
}

type RolesStore struct {
	db *sql.DB
}

func (r *RolesStore) GetByName(ctx context.Context, rolename string) (*Role, error) {
	query := `SELECT id, name, level, description FROM roles WHERE name = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	roles := &Role{}
	err := r.db.QueryRowContext(ctx, query, rolename).Scan(&roles.ID, &roles.Name, &roles.Level, &roles.Description)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return roles, nil
}
