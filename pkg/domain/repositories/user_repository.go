package repositories

import (
	"context"
	"errors"

	"go-starter-template/pkg/domain/entities"
)

var (
	ErrNoRows = errors.New("no results found")
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
}
