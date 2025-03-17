package command

import (
	"go-starter-template/pkg/application/result"
	"go-starter-template/pkg/domain/entities"
)

type CreateSignupCommand struct {
	Email    string
	Password string
}

type CreateSignupCommandResult struct {
	User *result.UserResult
}

func NewSignupUserCommandResult(user *entities.User) *CreateSignupCommandResult {
	return &CreateSignupCommandResult{
		User: result.NewUserResult(user),
	}
}

type CreateLoginCommand struct {
	Email    string
	Password string
	Remember bool
}

type CreateLoginCommandResult struct {
	Session *result.SessionResult
	User    *result.UserResult
}

func NewLoginUserCommandResult(user *entities.User) *CreateLoginCommandResult {
	return &CreateLoginCommandResult{
		User: result.NewUserResult(user),
	}
}
