package service

import (
	"go-starter-template/pkg/repository"
)

func NewServices(repo *repository.Repository) *Services {
	return &Services{
		Todo:    NewTodoService(repo.Todos),
		User:    NewUserService(repo.Users),
		Session: NewSessionService(repo.Sessions),
	}
}
