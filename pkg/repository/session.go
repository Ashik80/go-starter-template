package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-starter-template/pkg/entity"
	"time"

	"github.com/google/uuid"
)

type PQSessionStore struct {
	db *sql.DB
}

func NewPQSessionStore(db *sql.DB) SessionRepository {
	return &PQSessionStore{db}
}

func (s *PQSessionStore) Create(ctx context.Context, user *entity.User, expiresAt time.Time) (*entity.Session, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var session entity.Session
	var userId int

	err = tx.QueryRowContext(ctx,
		"INSERT INTO sessions (id, user_id, expires_at) VALUES ($1, $2, $3) RETURNING *",
		uuid.New(),
		user.ID,
		expiresAt,
	).Scan(
		&session.ID,
		&userId,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("failed to rollback transaction: %w; original error: %w", rollbackErr, err)
		}
		return nil, fmt.Errorf("failed to create session for user %s: %w", user.Email, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &session, nil
}

func (s *PQSessionStore) Get(ctx context.Context, sessionId string) (*entity.Session, error) {
	var session entity.Session
	err := s.db.QueryRowContext(ctx, "SELECT * FROM sessions WHERE id = $1", sessionId).Scan(
		&session.ID,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil && err == sql.ErrNoRows {
		return nil, newNotFoundError("session")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	return &session, nil
}

func (s *PQSessionStore) GetWithUser(ctx context.Context, sessionId string) (*entity.Session, error) {
	var session entity.Session
	var user entity.User
	query := `
		SELECT
			sessions.id, sessions.expires_at, sessions.created_at, sessions.updated_at,
			users.id as user_id, users.email, users.password, users.created_at, users.updated_at
		FROM sessions
		INNER JOIN users ON sessions.user_id = users.id
		WHERE sessions.id = $1
	`
	err := s.db.QueryRowContext(ctx, query, sessionId).Scan(
		&session.ID,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.UpdatedAt,
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil && err == sql.ErrNoRows {
		return nil, newNotFoundError("session")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	session.User = &user
	return &session, nil
}

func (s *PQSessionStore) Delete(ctx context.Context, sessionId string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM sessions WHERE id = $1", sessionId)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction: %w; original error: %w", rollbackErr, err)
		}
		return fmt.Errorf("failed to delete session: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
