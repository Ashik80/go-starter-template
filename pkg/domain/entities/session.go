package entities

import (
	"time"

	"github.com/google/uuid"

	"go-starter-template/pkg/domain/valueobject"
)

type Session struct {
	ID        uuid.UUID
	ExpiresAt valueobject.Time
	CreatedAt valueobject.Time
	UpdatedAt valueobject.Time
	User      *User
}

func NewSession(user *User) *Session {
	return &Session{
		ID:        uuid.New(),
		ExpiresAt: valueobject.Time(time.Now().Add(time.Hour * 1)),
		CreatedAt: valueobject.Time(time.Now()),
		UpdatedAt: valueobject.Time(time.Now()),
		User:      user,
	}
}

func (s *Session) Expired() bool {
	return time.Now().After(s.ExpiresAt.ToTime())
}

func (s *Session) ExtendByHour(hours int) {
	s.ExpiresAt = valueobject.Time(time.Now().Add(time.Hour * time.Duration(hours)))
}

func (s *Session) AddUser(user *User) {
	s.User = user
}
