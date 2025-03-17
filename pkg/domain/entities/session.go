package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"go-starter-template/pkg/domain/valueobject"
)

var (
	ErrSessionIsRequired = errors.New("Session is required")
)

type Session struct {
	ID        uuid.UUID
	ExpiresAt valueobject.Time
	CreatedAt valueobject.Time
	UpdatedAt valueobject.Time
	User      *User
}

func NewSession(user *User) *Session {
	currentTime := valueobject.NewCurrentTime()
	return &Session{
		ID:        uuid.New(),
		ExpiresAt: currentTime.ExtendByHour(1),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		User:      user,
	}
}

func (s *Session) Expired() bool {
	return time.Now().After(s.ExpiresAt.ToTime())
}

func (s *Session) SetExpiresAt(expiresAt valueobject.Time) error {
	if s == nil {
		return ErrSessionIsRequired
	}
	s.ExpiresAt = expiresAt
	s.UpdatedAt = valueobject.NewCurrentTime()
	return nil
}

func (s *Session) AddUser(user *User) error {
	if s.User != nil {
		return ErrSessionIsRequired
	}

	s.User = user
	s.UpdatedAt = valueobject.NewCurrentTime()

	return nil
}
