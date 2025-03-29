package query

import (
	"go-starter-template/internal/application/result"
	"go-starter-template/internal/domain/entities"
)

type GetSessionQuery struct {
	Session *result.SessionResult
}

func NewGetSessionQuery(session *entities.Session) *GetSessionQuery {
	return &GetSessionQuery{
		Session: result.NewSessionResult(session),
	}
}
