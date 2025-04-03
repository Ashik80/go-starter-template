package services

import (
	"context"

	"go-starter-template/internal/application/command"
	"go-starter-template/internal/application/query"
	"go-starter-template/internal/domain/entities"
	"go-starter-template/internal/domain/repositories"
	"go-starter-template/pkg/security"
)

type IUserService interface {
	Login(ctx context.Context, loginCommand *command.CreateLoginCommand) (*command.CreateLoginCommandResult, error)
	Signup(ctx context.Context, signupCommand *command.CreateSignupCommand) (*command.CreateSignupCommandResult, error)
	GetUserByEmail(ctx context.Context, email string) (*query.GetUserQuery, error)
}

type UserService struct {
	userRepository repositories.IUserRepository
	sessionService ISessionService
	passwordHasher security.IPasswordHasher
}

func NewUserService(
	userRepository repositories.IUserRepository,
	passwordHasher security.IPasswordHasher,
	sessionService ISessionService,
) IUserService {
	return &UserService{
		userRepository: userRepository,
		passwordHasher: passwordHasher,
		sessionService: sessionService,
	}
}

func (s *UserService) Login(ctx context.Context, loginCommand *command.CreateLoginCommand) (*command.CreateLoginCommandResult, error) {
	user, err := s.userRepository.GetByEmail(ctx, loginCommand.Email)
	if err != nil {
		return nil, err
	}

	if err = s.passwordHasher.CompareHashAndPassword(user.Password.ToString(), loginCommand.Password); err != nil {
		return nil, err
	}

	var extendSessionExpiryByHour int
	if loginCommand.Remember {
		extendSessionExpiryByHour = 30 * 24
	}

	sessionResult, err := s.sessionService.CreateSession(ctx, &command.CreateSessionCommand{
		User:         user,
		ExtendByHour: extendSessionExpiryByHour,
	})

	if err != nil {
		return nil, err
	}

	result := command.NewLoginUserCommandResult(user)
	result.Session = sessionResult.Session

	return result, nil
}

func (s *UserService) Signup(ctx context.Context, signupCommand *command.CreateSignupCommand) (*command.CreateSignupCommandResult, error) {
	existingUser, err := s.userRepository.GetByEmail(ctx, signupCommand.Email)
	if err != nil && err != repositories.ErrNoRows {
		return nil, err
	}
	if existingUser != nil {
		return nil, entities.ErrUserAlreadyExists
	}

	user, err := entities.NewUser(signupCommand.Email, signupCommand.Password)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := s.passwordHasher.GenerateFromPassword(signupCommand.Password, 10)
	if err != nil {
		return nil, err
	}
	user.SetPassword(hashedPassword)

	createdUser, err := s.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return command.NewSignupUserCommandResult(createdUser), nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*query.GetUserQuery, error) {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return query.NewGetUserQuery(user), nil
}
