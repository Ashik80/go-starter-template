package service

import (
	"context"
	"go-starter-template/pkg/entity"
	"time"
)

type TodoService interface {
	ListTodos(ctx context.Context) ([]*entity.Todo, error)
	GetTodo(ctx context.Context, id int) (*entity.Todo, error)
	CreateTodo(ctx context.Context, input CreateTodoInput) (*entity.Todo, error)
	UpdateTodo(ctx context.Context, id int, input UpdateTodoInput) (*entity.Todo, error)
	DeleteTodo(ctx context.Context, todo *entity.Todo) error
}

type UserService interface {
	CreateUser(ctx context.Context, input CreateUserInput) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
}

type SessionService interface {
	CreateSession(ctx context.Context, userID int, expiresAt time.Time) (*entity.Session, error)
	GetSession(ctx context.Context, token string) (*entity.Session, error)
	DeleteSession(ctx context.Context, token string) error
}

type Services struct {
	Todo    TodoService
	User    UserService
	Session SessionService
}

type CreateTodoInput struct {
	Title       string
	Description string
}

type UpdateTodoInput struct {
	Title       string
	Description string
}

type CreateUserInput struct {
	Email    string
	Password string
}
