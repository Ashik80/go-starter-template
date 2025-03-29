package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-starter-template/internal/domain/entities"
	"go-starter-template/internal/domain/repositories"
	"go-starter-template/internal/domain/valueobject"

	"github.com/google/uuid"
)

type SessionDTO struct {
	ID        uuid.UUID
	UserID    int
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s SessionDTO) toSession() *entities.Session {
	return &entities.Session{
		ID:        s.ID,
		ExpiresAt: valueobject.NewTime(s.ExpiresAt),
		CreatedAt: valueobject.NewTime(s.CreatedAt),
		UpdatedAt: valueobject.NewTime(s.UpdatedAt),
		User:      nil,
	}
}

type PQSessionRepository struct {
	db *sql.DB
}

func NewPQSessionRepository(db *sql.DB) repositories.SessionRepository {
	return &PQSessionRepository{db}
}

func (s *PQSessionRepository) Create(ctx context.Context, session *entities.Session) (*entities.Session, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var createdSession SessionDTO

	err = tx.QueryRowContext(ctx,
		"INSERT INTO sessions (id, user_id, expires_at) VALUES ($1, $2, $3) RETURNING *",
		session.ID,
		session.User.ID,
		session.ExpiresAt.ToTime(),
	).Scan(
		&createdSession.ID,
		&createdSession.UserID,
		&createdSession.ExpiresAt,
		&createdSession.CreatedAt,
		&createdSession.UpdatedAt,
	)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("failed to rollback transaction: %w; original error: %w", rollbackErr, err)
		}
		return nil, fmt.Errorf("failed to create session for user %s: %w", session.User.Email, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	newSession := createdSession.toSession()
	newSession.AddUser(session.User)

	return newSession, nil
}

func (s *PQSessionRepository) Get(ctx context.Context, sessionId string) (*entities.Session, error) {
	var session SessionDTO
	err := s.db.QueryRowContext(ctx, "SELECT * FROM sessions WHERE id = $1", sessionId).Scan(
		&session.ID,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	return session.toSession(), nil
}

func (s *PQSessionRepository) GetWithUser(ctx context.Context, sessionId string) (*entities.Session, error) {
	var sessionDTO SessionDTO
	var userDTO UserDTO

	query := `
		SELECT
			sessions.id, sessions.expires_at, sessions.created_at, sessions.updated_at,
			users.id as user_id, users.email, users.password, users.created_at, users.updated_at
		FROM sessions
		INNER JOIN users ON sessions.user_id = users.id
		WHERE sessions.id = $1
	`
	err := s.db.QueryRowContext(ctx, query, sessionId).Scan(
		&sessionDTO.ID,
		&sessionDTO.ExpiresAt,
		&sessionDTO.CreatedAt,
		&sessionDTO.UpdatedAt,
		&userDTO.ID,
		&userDTO.Email,
		&userDTO.Password,
		&userDTO.CreatedAt,
		&userDTO.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	session := sessionDTO.toSession()
	session.AddUser(userDTO.toUser())

	return session, nil
}

func (s *PQSessionRepository) Delete(ctx context.Context, session *entities.Session) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM sessions WHERE id = $1", session.ID)
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
