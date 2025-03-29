package interfaces

import (
	"context"

	"go-starter-template/internal/application/command"
	"go-starter-template/internal/application/query"
)

type SessionService interface {
	CreateSession(ctx context.Context, sessionCommand *command.CreateSessionCommand) (*command.CreateSessionCommandResult, error)
	GetSession(ctx context.Context, sessionId string) (*query.GetSessionQuery, error)
	DeleteSession(ctx context.Context, sessionId string) error
}
