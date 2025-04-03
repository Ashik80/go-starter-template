package services

import (
	"context"

	"go-starter-template/internal/application/command"
	"go-starter-template/internal/application/query"
	"go-starter-template/internal/domain/entities"
	"go-starter-template/internal/domain/repositories"
)

type ITodoService interface {
	ListTodos(ctx context.Context) (*query.GetTodoListQuery, error)
	GetTodo(ctx context.Context, id int) (*query.GetTodoQuery, error)
	CreateTodo(ctx context.Context, todoCommand *command.CreateTodoCommand) (*command.CreateTodoCommandResult, error)
	UpdateTodo(ctx context.Context, todoCommand *command.UpdateTodoCommand) (*command.UpdateTodoCommandResult, error)
	DeleteTodo(ctx context.Context, todoCommand *command.DeleteTodoCommand) error
}

type TodoService struct {
	todoRepository repositories.ITodoRepository
}

func NewTodoService(repo repositories.ITodoRepository) ITodoService {
	return &TodoService{
		todoRepository: repo,
	}
}

func (s *TodoService) ListTodos(ctx context.Context) (*query.GetTodoListQuery, error) {
	todos, err := s.todoRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	return query.NewGetTodoListQuery(todos), nil
}

func (s *TodoService) GetTodo(ctx context.Context, id int) (*query.GetTodoQuery, error) {
	todo, err := s.todoRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return query.NewGetTodoQuery(todo), nil
}

func (s *TodoService) CreateTodo(ctx context.Context, todoCommand *command.CreateTodoCommand) (*command.CreateTodoCommandResult, error) {
	todo, err := entities.NewTodo(todoCommand.Title, todoCommand.Description)
	if err != nil {
		return nil, err
	}

	result, err := s.todoRepository.Create(ctx, todo)
	if err != nil {
		return nil, err
	}

	return command.NewCreateTodoCommandResult(result), nil
}

func (s *TodoService) UpdateTodo(ctx context.Context, todoCommand *command.UpdateTodoCommand) (*command.UpdateTodoCommandResult, error) {
	updatedTodo, err := entities.NewTodoWithID(todoCommand.ID, todoCommand.Title, todoCommand.Description)
	if err != nil {
		return nil, err
	}

	updatedTodo, err = s.todoRepository.Update(ctx, updatedTodo)
	if err != nil {
		return nil, err
	}

	return command.NewUpdateTodoCommandResult(updatedTodo), nil
}

func (s *TodoService) DeleteTodo(ctx context.Context, todoCommand *command.DeleteTodoCommand) error {
	todo, err := s.todoRepository.Get(ctx, todoCommand.ID)
	if err != nil {
		return err
	}
	return s.todoRepository.Delete(ctx, todo)
}
