package command

import (
	"go-starter-template/internal/application/result"
	"go-starter-template/internal/domain/entities"
)

type CreateTodoCommand struct {
	Title       string
	Description string
}

type CreateTodoCommandResult struct {
	Todo *result.TodoResult
}

func NewCreateTodoCommandResult(todo *entities.Todo) *CreateTodoCommandResult {
	return &CreateTodoCommandResult{
		Todo: result.NewTodoResult(todo),
	}
}

type UpdateTodoCommand struct {
	ID          int
	Title       string
	Description string
}

type UpdateTodoCommandResult struct {
	Todo *result.TodoResult
}

func NewUpdateTodoCommandResult(todo *entities.Todo) *UpdateTodoCommandResult {
	return &UpdateTodoCommandResult{
		Todo: result.NewTodoResult(todo),
	}
}

type DeleteTodoCommand struct {
	ID int
}
