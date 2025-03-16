package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"go-starter-template/pkg/domain/entities"
	"go-starter-template/pkg/domain/repositories"
)

type PQUserStore struct {
	db *sql.DB
}

func NewPQUserStore(db *sql.DB) repositories.UserRepository {
	return &PQUserStore{db}
}

func (s *PQUserStore) Create(ctx context.Context, user *entities.User) (*entities.User, error) {
	var createdUser entities.User
	err := s.db.QueryRowContext(ctx,
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *",
		user.Email, user.Password,
	).Scan(
		&createdUser.ID,
		&createdUser.Email,
		&createdUser.Password,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &createdUser, nil
}

func (s *PQUserStore) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	err := s.db.QueryRowContext(ctx, "SELECT * FROM users WHERE email = $1", email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repositories.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}
