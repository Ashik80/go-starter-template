package repositories

import (
	"context"

	"go-starter-template/internal/domain/entities"
)

type ITodoRepository interface {
	List(ctx context.Context) ([]*entities.Todo, error)
	Get(ctx context.Context, id int) (*entities.Todo, error)
	Create(ctx context.Context, todoDto *entities.Todo) (*entities.Todo, error)
	Update(ctx context.Context, todoDto *entities.Todo) (*entities.Todo, error)
	Delete(ctx context.Context, todo *entities.Todo) error
}
