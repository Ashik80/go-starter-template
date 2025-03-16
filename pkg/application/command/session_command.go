package command

import (
	"go-starter-template/pkg/application/result"
	"go-starter-template/pkg/domain/entities"
)

type CreateSessionCommand struct {
	User         *entities.User
	ExtendByHour int
}

type CreateSessionCommandResult struct {
	Session *result.SessionResult
}

func NewCreateSessionCommandResult(session *entities.Session) *CreateSessionCommandResult {
	return &CreateSessionCommandResult{
		Session: result.NewSessionResult(session),
	}
}
