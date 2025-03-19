package query

import (
	"go-starter-template/pkg/application/result"
	"go-starter-template/pkg/domain/entities"
)

type GetUserQuery struct {
	User *result.UserResult
}

func NewGetUserQuery(user *entities.User) *GetUserQuery {
	return &GetUserQuery{
		User: result.NewUserResult(user),
	}
}
