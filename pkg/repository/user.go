package repository

import (
	"context"
	"database/sql"
	"fmt"

	"go-starter-template/pkg/entity"
)

type PQUserStore struct {
	db *sql.DB
}

func NewPQUserStore(db *sql.DB) UserRepository {
	return &PQUserStore{db}
}

func (s *PQUserStore) Create(ctx context.Context, email string, passwordHash string) (*entity.User, error) {
	var user entity.User
	err := s.db.QueryRowContext(ctx,
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *",
		email, passwordHash,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &user, nil
}

func (s *PQUserStore) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := s.db.QueryRowContext(ctx, "SELECT * FROM users WHERE email = $1", email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil && err == sql.ErrNoRows {
		return nil, newNotFoundError("user")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}
