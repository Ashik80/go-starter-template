package entities

import (
	"errors"

	"go-starter-template/pkg/domain/valueobject"
)

var (
	ErrUserIsRequired    = errors.New("User is required")
	ErrUserAlreadyExists = errors.New("User already exists")
)

type User struct {
	ID        int
	Email     valueobject.Email
	Password  valueobject.Password
	CreatedAt valueobject.Time
	UpdatedAt valueobject.Time
	Sessions  []*Session
}

func NewUser(email, password string) (*User, error) {
	emailValid, err := valueobject.NewEmail(email)
	if err != nil {
		return nil, err
	}

	passwordValid, err := valueobject.NewPassword(password)
	if err != nil {
		return nil, err
	}

	currentTime := valueobject.NewCurrentTime()

	return &User{
		Email:     emailValid,
		Password:  passwordValid,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Sessions:  []*Session{},
	}, nil
}

func (u *User) SetEmail(email string) error {
	if u == nil {
		return ErrUserIsRequired
	}

	emailValid, err := valueobject.NewEmail(email)
	if err != nil {
		return err
	}

	u.Email = emailValid
	u.UpdatedAt = valueobject.NewCurrentTime()

	return nil
}

func (u *User) SetPassword(password string) error {
	if u == nil {
		return ErrUserIsRequired
	}

	pass, err := valueobject.NewPassword(password)
	if err != nil {
		return err
	}

	u.Password = pass
	u.UpdatedAt = valueobject.NewCurrentTime()

	return nil
}
