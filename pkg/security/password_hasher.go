package security

import (
	"fmt"

	"go-starter-template/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

type IPasswordHasher interface {
	GenerateFromPassword(password string, cost int) (string, error)
	CompareHashAndPassword(hashedPassword, password string) error
}

type BcryptPasswordHasher struct {
	log *logger.Logger
}

func NewBcryptPasswordHasher(log *logger.Logger) IPasswordHasher {
	return &BcryptPasswordHasher{log}
}

func (b *BcryptPasswordHasher) GenerateFromPassword(password string, cost int) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		b.log.Error("failed to generate password hash: %v", err)
		return "", fmt.Errorf("failed to generate password hash: %w", err)
	}
	return string(hash), nil
}

func (b *BcryptPasswordHasher) CompareHashAndPassword(hashedPassword, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return fmt.Errorf("failed to compare password: %w", err)
	}
	return nil
}
