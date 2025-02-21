package store

import (
	"context"
	"fmt"
	"time"

	"go-starter-template/ent"
	"go-starter-template/ent/todo"
)

type (
	Todo struct {
		ID          int       `json:"id,omitempty"`
		Title       string    `json:"title,omitempty"`
		Description string    `json:"description,omitempty"`
		CreatedAt   time.Time `json:"created_at,omitempty"`
		UpdatedAt   time.Time `json:"updated_at,omitempty"`
	}

	TodoStore interface {
		List(ctx context.Context) ([]*Todo, error)
		Get(ctx context.Context, id int) (*Todo, error)
		Create(ctx context.Context, todoDto *TodoCreateDto) (*Todo, error)
		Update(ctx context.Context, todo *Todo, todoDto *TodoCreateDto) (*Todo, error)
		Delete(ctx context.Context, todo *Todo) error
	}
)

type (
	EntTodoStore struct {
		orm *ent.Client
	}

	TodoCreateDto struct {
		Title       string `json:"title"`
		Description string `json:"description,omitempty"`
	}
)

func NewEntTodoStore(orm *ent.Client) *EntTodoStore {
	return &EntTodoStore{orm}
}

func NewTodoCreateDto(title string) *TodoCreateDto {
	return &TodoCreateDto{
		Title: title,
	}
}

func (t *EntTodoStore) List(ctx context.Context) ([]*Todo, error) {
	todos, err := t.orm.Todo.Query().Order(todo.ByID()).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get todos: %w\n", err)
	}
	return mapTodos(todos), nil
}

func (t *EntTodoStore) Get(ctx context.Context, id int) (*Todo, error) {
	todo, err := t.orm.Todo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo: %w\n", err)
	}
	return mapTodo(todo), nil
}

func (t *EntTodoStore) Create(ctx context.Context, todoDto *TodoCreateDto) (*Todo, error) {
	query := t.orm.Todo.Create().SetTitle(todoDto.Title)

	if todoDto.Description != "" {
		query.SetDescription(todoDto.Description)
	}

	todo, err := query.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create todos: %w\n", err)
	}
	return mapTodo(todo), nil
}

func (t *EntTodoStore) Update(ctx context.Context, todo *Todo, todoDto *TodoCreateDto) (*Todo, error) {
	query := t.orm.Todo.UpdateOneID(todo.ID).SetTitle(todoDto.Title)

	if todoDto.Description != "" {
		query.SetDescription(todoDto.Description)
	}

	updatedTodo, err := query.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update todo: %w\n", err)
	}
	return mapTodo(updatedTodo), nil
}

func (t *EntTodoStore) Delete(ctx context.Context, todo *Todo) error {
	if err := t.orm.Todo.DeleteOneID(todo.ID).Exec(ctx); err != nil {
		return fmt.Errorf("failed to delete todo: %w\n", err)
	}
	return nil
}

func mapTodo(todo *ent.Todo) *Todo {
	return &Todo{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}
}

func mapTodos(todos []*ent.Todo) []*Todo {
	var ts []*Todo
	for _, todo := range todos {
		ts = append(ts, mapTodo(todo))
	}
	return ts
}
