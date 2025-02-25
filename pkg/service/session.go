package service

import (
	"context"
	"time"

	"go-starter-template/pkg/entity"
	"go-starter-template/pkg/repository"
)

type sessionService struct {
	repository repository.SessionRepository
}

func NewSessionService(repo repository.SessionRepository) SessionService {
	return &sessionService{
		repository: repo,
	}
}

func (s *sessionService) CreateSession(ctx context.Context, userID int, expiresAt time.Time) (*entity.Session, error) {
	user := &entity.User{ID: userID}
	return s.repository.Create(ctx, user, expiresAt)
}

func (s *sessionService) GetSession(ctx context.Context, token string) (*entity.Session, error) {
	return s.repository.GetWithUser(ctx, token)
}

func (s *sessionService) DeleteSession(ctx context.Context, token string) error {
	return s.repository.Delete(ctx, token)
}
