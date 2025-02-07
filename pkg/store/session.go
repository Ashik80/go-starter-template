package store

import (
	"context"
	"fmt"
	"go-starter-template/ent"
	"go-starter-template/ent/session"
	"time"

	"github.com/google/uuid"
)

type (
	Session struct {
		ID        uuid.UUID `json:"id,omitempty"`
		ExpiresAt time.Time `json:"expires_at,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
		User      *User     `json:"user,omitempty"`
	}

	SessionStore interface {
		Create(ctx context.Context, user *User, expiresAt time.Time) (*Session, error)
		Get(ctx context.Context, sessionId string) (*Session, error)
		GetWithUser(ctx context.Context, sessionId string) (*Session, error)
	}
)

type EntSessionStore struct {
	orm *ent.Client
}

func NewEntSessionStore(orm *ent.Client) *EntSessionStore {
	return &EntSessionStore{orm}
}

func (s *EntSessionStore) Create(ctx context.Context, user *User, expiresAt time.Time) (*Session, error) {
	sess, err := s.orm.Session.Create().SetExpiresAt(expiresAt).SetUserID(user.ID).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create session for user %s: %v\n", user.Email, err)
	}
	return mapSession(sess), err
}

func (s *EntSessionStore) Get(ctx context.Context, sessionId string) (*Session, error) {
	sess, err := s.orm.Session.Get(ctx, uuid.MustParse(sessionId))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve session: %v\n", err)
	}
	return mapSession(sess), err
}

func (s *EntSessionStore) GetWithUser(ctx context.Context, sessionId string) (*Session, error) {
	sess, err := s.orm.Session.Query().Where(session.ID(uuid.MustParse(sessionId))).WithUser().Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve session: %v\n", err)
	}
	mappedSession := mapSession(sess)
	mappedSession.User = mapUser(sess.Edges.User)
	return mappedSession, err
}

func mapSession(sess *ent.Session) *Session {
	return &Session{
		ID:        sess.ID,
		ExpiresAt: sess.ExpiresAt,
		CreatedAt: sess.CreatedAt,
		UpdatedAt: sess.UpdatedAt,
	}
}
