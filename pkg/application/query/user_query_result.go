package query

import (
	"go-starter-template/pkg/application/result"
	"go-starter-template/pkg/domain/entities"
)

type UserQueryResult struct {
	User *result.UserResult
}

func NewUserQueryResult(user *entities.User) *UserQueryResult {
	return &UserQueryResult{
		User: result.NewUserResult(user),
	}
}
