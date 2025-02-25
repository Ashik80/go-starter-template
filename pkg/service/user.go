package service

import (
	"context"
	"go-starter-template/pkg/entity"
	"go-starter-template/pkg/repository"
)

type userService struct {
	repository repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repository: repo,
	}
}

func (s *userService) CreateUser(ctx context.Context, input CreateUserInput) (*entity.User, error) {
	return s.repository.Create(ctx, input.Email, input.Password)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.repository.GetByEmail(ctx, email)
}
