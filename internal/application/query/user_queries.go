package query

import (
	"go-starter-template/internal/application/result"
	"go-starter-template/internal/domain/entities"
)

type GetUserQuery struct {
	User *result.UserResult
}

func NewGetUserQuery(user *entities.User) *GetUserQuery {
	return &GetUserQuery{
		User: result.NewUserResult(user),
	}
}
