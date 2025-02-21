package service

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	GenerateFromPassword(password string, cost int) (string, error)
	CompareHashAndPassword(hashedPassword, password string) error
}

type BcryptPasswordHasher struct{}

func NewBcryptPasswordHasher() PasswordHasher {
	return &BcryptPasswordHasher{}
}

func (b *BcryptPasswordHasher) GenerateFromPassword(password string, cost int) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		log.Printf("ERROR: failed to generate password hash: %v\n", err)
		return "", err
	}
	return string(hash), nil
}

func (b *BcryptPasswordHasher) CompareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
