package store

import (
	"context"
	"database/sql"
	"time"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct{}

func (m *MockUserStore) Create(context.Context, *sql.Tx, *User) error {
	return nil
}

func (m *MockUserStore) GetByID(context.Context, int) (*User, error) {
	return nil, nil
}

func (m *MockUserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return nil
}

func (m *MockUserStore) Activate(context.Context, string) error {
	return nil
}

func (m *MockUserStore) Delete(context.Context, int) error {
	return nil
}

func (m *MockUserStore) GetByEmail(context.Context, string) (*User, error) {
	return nil, nil
}
