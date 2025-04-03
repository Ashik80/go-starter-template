package services

import (
	"context"

	"go-starter-template/internal/application/command"
	"go-starter-template/internal/application/query"
	"go-starter-template/internal/domain/entities"
	"go-starter-template/internal/domain/repositories"
)

type ISessionService interface {
	CreateSession(ctx context.Context, sessionCommand *command.CreateSessionCommand) (*command.CreateSessionCommandResult, error)
	GetSession(ctx context.Context, sessionId string) (*query.GetSessionQuery, error)
	DeleteSession(ctx context.Context, sessionId string) error
}

type SessionService struct {
	sessionRepository repositories.ISessionRepository
}

func NewSessionService(sessionRepository repositories.ISessionRepository) ISessionService {
	return &SessionService{
		sessionRepository: sessionRepository,
	}
}

func (s *SessionService) CreateSession(ctx context.Context, sessionCommand *command.CreateSessionCommand) (*command.CreateSessionCommandResult, error) {
	newSession := entities.NewSession(sessionCommand.User)
	if sessionCommand.ExtendByHour > 0 {
		newSession.ExpiresAt = newSession.ExpiresAt.ExtendByHour(sessionCommand.ExtendByHour)
	}
	session, err := s.sessionRepository.Create(ctx, newSession)
	if err != nil {
		return nil, err
	}

	return command.NewCreateSessionCommandResult(session), nil
}

func (s *SessionService) GetSession(ctx context.Context, sessionId string) (*query.GetSessionQuery, error) {
	session, err := s.sessionRepository.GetWithUser(ctx, sessionId)
	if err != nil {
		return nil, err
	}
	return query.NewGetSessionQuery(session), nil
}

func (s *SessionService) DeleteSession(ctx context.Context, sessionId string) error {
	session, err := s.sessionRepository.Get(ctx, sessionId)
	if err != nil {
		return err
	}
	return s.sessionRepository.Delete(ctx, session)
}
