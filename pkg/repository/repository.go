package repository

import (
	"context"
	"time"

	"go-starter-template/pkg/entity"
)

type UserRepository interface {
	Create(ctx context.Context, email string, passwordHash string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}

type SessionRepository interface {
	Create(ctx context.Context, user *entity.User, expiresAt time.Time) (*entity.Session, error)
	Get(ctx context.Context, sessionId string) (*entity.Session, error)
	GetWithUser(ctx context.Context, sessionId string) (*entity.Session, error)
	Delete(ctx context.Context, sessionId string) error
}

type TodoRepository interface {
	List(ctx context.Context) ([]*entity.Todo, error)
	Get(ctx context.Context, id int) (*entity.Todo, error)
	Create(ctx context.Context, todoDto *TodoCreateDto) (*entity.Todo, error)
	Update(ctx context.Context, id int, todoDto *TodoCreateDto) (*entity.Todo, error)
	Delete(ctx context.Context, todo *entity.Todo) error
}

type Repository struct {
	Users    UserRepository
	Sessions SessionRepository
	Todos    TodoRepository
}
