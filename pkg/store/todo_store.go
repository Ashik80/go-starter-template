package store

import (
	"context"
	"fmt"
	"gohtmx/ent"
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

func (t *TodoStore) List(ctx context.Context) ([]*ent.Todo, error) {
	todos, err := t.orm.Todo.Query().All(ctx)
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

func (t *TodoStore) Create(ctx context.Context, todoDto TodoCreateDto) (*ent.Todo, error) {
	todo, err := t.orm.Todo.Create().SetTitle(todoDto.Title).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create todos: %v", err)
	}
	return todo, nil
}
