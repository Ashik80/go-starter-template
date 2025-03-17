package services

import (
	"context"

	"go-starter-template/pkg/application/command"
	"go-starter-template/pkg/application/interfaces"
	"go-starter-template/pkg/application/query"
	"go-starter-template/pkg/domain/entities"
	"go-starter-template/pkg/domain/repositories"
)

type SessionService struct {
	sessionRepository repositories.SessionRepository
}

func NewSessionService(sessionRepository repositories.SessionRepository) interfaces.SessionService {
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

func (s *SessionService) GetSession(ctx context.Context, sessionId string) (*query.SessionQueryResult, error) {
	session, err := s.sessionRepository.GetWithUser(ctx, sessionId)
	if err != nil {
		return nil, err
	}
	return query.NewSessionQueryResult(session), nil
}

func (s *SessionService) DeleteSession(ctx context.Context, sessionId string) error {
	session, err := s.sessionRepository.Get(ctx, sessionId)
	if err != nil {
		return err
	}
	return s.sessionRepository.Delete(ctx, session)
}
