package interfaces

import (
	"context"
	"go-starter-template/pkg/application/command"
	"go-starter-template/pkg/application/query"
)

type SessionService interface {
	CreateSession(ctx context.Context, sessionCommand *command.CreateSessionCommand) (*command.CreateSessionCommandResult, error)
	GetSession(ctx context.Context, sessionId string) (*query.SessionQueryResult, error)
	DeleteSession(ctx context.Context, sessionId string) error
}
