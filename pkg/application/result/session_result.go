package result

import (
	"time"

	"github.com/google/uuid"

	"go-starter-template/pkg/domain/entities"
)

type SessionResult struct {
	ID        uuid.UUID
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	User      *UserResult
}

func (s *SessionResult) Expired() bool {
	return s.ExpiresAt.Before(time.Now())
}

func NewSessionResult(session *entities.Session) *SessionResult {
	return &SessionResult{
		ID:        session.ID,
		ExpiresAt: session.ExpiresAt.ToTime(),
		CreatedAt: session.CreatedAt.ToTime(),
		UpdatedAt: session.UpdatedAt.ToTime(),
		User:      NewUserResult(session.User),
	}
}
