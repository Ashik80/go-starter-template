package service

import (
	"context"

	"go-starter-template/pkg/entity"
	"go-starter-template/pkg/repository"
)

type todoService struct {
	repository repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{
		repository: repo,
	}
}

func (s *todoService) ListTodos(ctx context.Context) ([]*entity.Todo, error) {
	return s.repository.List(ctx)
}

func (s *todoService) GetTodo(ctx context.Context, id int) (*entity.Todo, error) {
	return s.repository.Get(ctx, id)
}

func (s *todoService) CreateTodo(ctx context.Context, input CreateTodoInput) (*entity.Todo, error) {
	dto := &repository.TodoCreateDto{
		Title:       input.Title,
		Description: input.Description,
	}
	return s.repository.Create(ctx, dto)
}

func (s *todoService) UpdateTodo(ctx context.Context, id int, input UpdateTodoInput) (*entity.Todo, error) {
	dto := &repository.TodoCreateDto{
		Title:       input.Title,
		Description: input.Description,
	}
	return s.repository.Update(ctx, id, dto)
}

func (s *todoService) DeleteTodo(ctx context.Context, todo *entity.Todo) error {
	return s.repository.Delete(ctx, todo)
}
