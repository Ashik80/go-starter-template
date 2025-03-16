package entities

import (
	"errors"
	"net/mail"
	"time"
	"unicode"

	"go-starter-template/pkg/domain/valueobject"
)

var (
	ErrEmailIsRequired    = errors.New("Email is required")
	ErrPasswordIsRequired = errors.New("Password is required")
	ErrEmailIsInvalid     = errors.New("Email is invalid")
	ErrUserAlreadyExists  = errors.New("User already exists")
)

type User struct {
	ID        int
	Email     string
	Password  string
	CreatedAt valueobject.Time
	UpdatedAt valueobject.Time
	Sessions  []*Session
}

func NewUser(email, password string) *User {
	return &User{
		Email:     email,
		Password:  password,
		CreatedAt: valueobject.Time(time.Now()),
		UpdatedAt: valueobject.Time(time.Now()),
	}
}

func (u *User) Validate() error {
	if u.Email == "" {
		return ErrEmailIsRequired
	}
	if u.Password == "" {
		return ErrPasswordIsRequired
	}
	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		return ErrEmailIsInvalid
	}
	if err := u.IsStrongPassword(); err != nil {
		return err
	}
	return nil
}

func (u *User) SetPassword(password string) {
	u.Password = password
}

type PasswordValidationError struct {
	Errors []string
}

func (p *PasswordValidationError) Error() string {
	return "Password must have minimum length of 8 characters and must contain at least 1 uppercase letter, 1 lowercase letter, 1 digit and 1 symbol"
}

func (u *User) IsStrongPassword() *PasswordValidationError {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	pve := &PasswordValidationError{
		Errors: []string{},
	}
	if len(u.Password) > 7 {
		hasMinLen = true
	}
	for _, char := range u.Password {
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsLower(char) {
			hasLower = true
		}
		if unicode.IsDigit(char) {
			hasNumber = true
		}
		if unicode.IsSymbol(char) || unicode.IsPunct(char) {
			hasSpecial = true
		}
	}
	if !hasMinLen {
		pve.Errors = append(pve.Errors, "Password must be at least 8 characters long")
	}
	if !hasUpper {
		pve.Errors = append(pve.Errors, "Password must have at least 1 uppercase character")
	}
	if !hasLower {
		pve.Errors = append(pve.Errors, "Password must have at least 1 lowercase character")
	}
	if !hasNumber {
		pve.Errors = append(pve.Errors, "Password must have at least 1 digit")
	}
	if !hasSpecial {
		pve.Errors = append(pve.Errors, "Password must have at least 1 special character or symbol")
	}
	if len(pve.Errors) > 0 {
		return pve
	}
	return pve
}
