package query

import (
	"go-starter-template/pkg/application/result"
	"go-starter-template/pkg/domain/entities"
)

type SessionQueryResult struct {
	Session *result.SessionResult
}

func NewSessionQueryResult(session *entities.Session) *SessionQueryResult {
	return &SessionQueryResult{
		Session: result.NewSessionResult(session),
	}
}
