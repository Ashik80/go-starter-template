package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-starter-template/pkg/domain/entities"
	"go-starter-template/pkg/domain/repositories"
	"go-starter-template/pkg/domain/valueobject"
)

type UserDTO struct {
	ID        int
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u UserDTO) toUser() *entities.User {
	email, _ := valueobject.NewEmail(u.Email)
	password, _ := valueobject.NewPassword(u.Password)
	return &entities.User{
		ID:        int(u.ID),
		Email:     email,
		Password:  password,
		CreatedAt: valueobject.NewTime(u.CreatedAt),
		UpdatedAt: valueobject.NewTime(u.UpdatedAt),
	}
}

type PQUserStore struct {
	db *sql.DB
}

func NewPQUserStore(db *sql.DB) repositories.UserRepository {
	return &PQUserStore{db}
}

func (s *PQUserStore) Create(ctx context.Context, user *entities.User) (*entities.User, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var userDTO UserDTO
	err = tx.QueryRowContext(ctx,
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *",
		user.Email.ToString(), user.Password.ToString(),
	).Scan(
		&userDTO.ID,
		&userDTO.Email,
		&userDTO.Password,
		&userDTO.CreatedAt,
		&userDTO.UpdatedAt,
	)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("failed to rollback transaction: %w; original error: %w", rollbackErr, err)
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return userDTO.toUser(), nil
}

func (s *PQUserStore) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user UserDTO
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
	return user.toUser(), nil
}
