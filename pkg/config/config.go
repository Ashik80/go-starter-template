package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		Env            string
		Port           string
		CSRFAuthKey    string
		DatabaseConfig *DatabaseConfig
		AllowedOrigins string
	}

	DatabaseConfig struct {
		User     string
		Password string
		Name     string
		Host     string
		Port     string
		SslMode  string
	}
)

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env file %w\n", err)
	}

	config := &Config{
		Env:            os.Getenv("ENV"),
		Port:           os.Getenv("PORT"),
		CSRFAuthKey:    os.Getenv("CSRF_AUTH_KEY"),
		AllowedOrigins: os.Getenv("ALLOWED_ORIGINS"),
		DatabaseConfig: &DatabaseConfig{
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			SslMode:  os.Getenv("DB_SSL_MODE"),
		},
	}

	return config, nil
}
