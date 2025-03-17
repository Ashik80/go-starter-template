package services

import (
	"context"

	"go-starter-template/pkg/application/command"
	"go-starter-template/pkg/application/interfaces"
	"go-starter-template/pkg/application/query"
	"go-starter-template/pkg/domain/entities"
	"go-starter-template/pkg/domain/repositories"
	"go-starter-template/pkg/infrastructure/security"
)

type UserService struct {
	userRepository repositories.UserRepository
	sessionService interfaces.SessionService
	passwordHasher security.PasswordHasher
}

func NewUserService(
	userRepository repositories.UserRepository,
	passwordHasher security.PasswordHasher,
	sessionService interfaces.SessionService,
) interfaces.UserService {
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

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*query.UserQueryResult, error) {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return query.NewUserQueryResult(user), nil
}
