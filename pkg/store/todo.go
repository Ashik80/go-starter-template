package store

import (
	"context"
	"fmt"

	"go-starter-template/ent"
	"go-starter-template/ent/todo"
)

type (
	TodoStore struct {
		orm *ent.Client
	}

	TodoCreateDto struct {
		Title string `json:"title"`
	}
)

func NewTodoStore(orm *ent.Client) *TodoStore {
	return &TodoStore{orm}
}

func NewTodoCreateDto(title string) *TodoCreateDto {
	return &TodoCreateDto{
		Title: title,
	}
}

func (t *TodoStore) List(ctx context.Context) ([]*ent.Todo, error) {
	todos, err := t.orm.Todo.Query().Order(todo.ByID()).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get todos: %v", err)
	}
	return todos, nil
}

func (t *TodoStore) Get(ctx context.Context, id int) (*ent.Todo, error) {
	todo, err := t.orm.Todo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo: %v", err)
	}
	return todo, nil
}

func (t *TodoStore) Create(ctx context.Context, todoDto *TodoCreateDto) (*ent.Todo, error) {
	todo, err := t.orm.Todo.Create().SetTitle(todoDto.Title).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create todos: %v", err)
	}
	return todo, nil
}

func (t *TodoStore) Update(ctx context.Context, todo *ent.Todo, todoDto TodoCreateDto) (*ent.Todo, error) {
	updatedTodo, err := todo.Update().SetTitle(todoDto.Title).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update todo: %v", err)
	}
	return updatedTodo, nil
}

func (t *TodoStore) Delete(ctx context.Context, todo *ent.Todo) error {
	if err := t.orm.Todo.DeleteOne(todo).Exec(ctx); err != nil {
		return fmt.Errorf("failed to delete todo: %v", err)
	}
	return nil
}
