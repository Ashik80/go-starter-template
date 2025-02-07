package store

import (
	"context"
	"fmt"
	"time"

	"go-starter-template/ent"
	"go-starter-template/ent/user"
)

type (
	User struct {
		ID        int        `json:"id,omitempty"`
		Email     string     `json:"email,omitempty"`
		Password  string     `json:"password,omitempty"`
		CreatedAt time.Time  `json:"created_at,omitempty"`
		UpdatedAt time.Time  `json:"updated_at,omitempty"`
		Sessions  []*Session `json:"sessions,omitempty"`
	}

	UserStore interface {
		Create(ctx context.Context, email string, passwordHash string) (*User, error)
		GetByEmail(ctx context.Context, email string) (*User, error)
	}
)

type EntUserStore struct {
	orm *ent.Client
}

func NewEntUserStore(orm *ent.Client) *EntUserStore {
	return &EntUserStore{orm}
}

func (s *EntUserStore) Create(ctx context.Context, email string, passwordHash string) (*User, error) {
	query := s.orm.User.Create().SetEmail(email).SetPassword(passwordHash)
	user, err := query.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v\n", err)
	}
	return mapUser(user), nil
}

func (s *EntUserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := s.orm.User.Query().Where(user.Email(email))
	user, err := query.First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find user with email %s: %v\n", email, err)
	}
	return mapUser(user), nil
}

func mapUser(user *ent.User) *User {
	return &User{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
