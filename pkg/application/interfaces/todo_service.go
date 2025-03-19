package interfaces

import (
	"context"

	"go-starter-template/pkg/application/command"
	"go-starter-template/pkg/application/query"
)

type TodoService interface {
	ListTodos(ctx context.Context) (*query.GetTodoListQuery, error)
	GetTodo(ctx context.Context, id int) (*query.GetTodoQuery, error)
	CreateTodo(ctx context.Context, todoCommand *command.CreateTodoCommand) (*command.CreateTodoCommandResult, error)
	UpdateTodo(ctx context.Context, todoCommand *command.UpdateTodoCommand) (*command.UpdateTodoCommandResult, error)
	DeleteTodo(ctx context.Context, todoCommand *command.DeleteTodoCommand) error
}
