package infrastructure

import (
	"database/sql"
	"fmt"
)

func NewDatabaseConfig(conf *Config) (*sql.DB, error) {
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
