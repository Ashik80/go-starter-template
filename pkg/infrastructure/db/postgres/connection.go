package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"go-starter-template/pkg/infrastructure/config"
)

func NewDatabaseConfig(conf *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.DatabaseConfig.Host,
		conf.DatabaseConfig.Port,
		conf.DatabaseConfig.User,
		conf.DatabaseConfig.Password,
		conf.DatabaseConfig.Name,
	)
	return sql.Open("postgres", dsn)
}
