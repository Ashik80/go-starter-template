package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		Port           string
		CSRFAuthKey    string
		DatabaseConfig *DatabaseConfig
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
		return nil, fmt.Errorf("failed to load .env file %v\n", err)
	}

	config := &Config{
		Port:        os.Getenv("PORT"),
		CSRFAuthKey: os.Getenv("CSRF_AUTH_KEY"),
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
