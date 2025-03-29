package valueobject

import (
	"errors"
	"net/mail"
)

var (
	ErrEmailIsRequired = errors.New("Email is required")
	ErrEmailIsInvalid  = errors.New("Email is invalid")
)

type Email struct {
	value string
}

func NewEmail(email string) (Email, error) {
	if email == "" {
		return Email{}, ErrEmailIsRequired
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return Email{}, ErrEmailIsInvalid
	}
	return Email{value: email}, nil
}

func (e Email) ToString() string {
	return e.value
}
