package interfaces

import (
	"context"
	"go-starter-template/pkg/application/command"
	"go-starter-template/pkg/application/query"
)

type UserService interface {
	Login(ctx context.Context, loginCommand *command.CreateLoginCommand) (*command.CreateLoginCommandResult, error)
	Signup(ctx context.Context, signupCommand *command.CreateSignupCommand) (*command.CreateSignupCommandResult, error)
	GetUserByEmail(ctx context.Context, email string) (*query.GetUserQuery, error)
}
