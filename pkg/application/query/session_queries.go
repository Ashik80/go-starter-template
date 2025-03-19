package query

import (
	"go-starter-template/pkg/application/result"
	"go-starter-template/pkg/domain/entities"
)

type GetSessionQuery struct {
	Session *result.SessionResult
}

func NewGetSessionQuery(session *entities.Session) *GetSessionQuery {
	return &GetSessionQuery{
		Session: result.NewSessionResult(session),
	}
}
