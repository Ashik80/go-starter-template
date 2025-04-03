package factories

import (
	"database/sql"

	"go-starter-template/internal/application/services"
	"go-starter-template/internal/infrastructure/db/postgres"
)

func NewTodoServiceWithPQRepository(db *sql.DB) services.ITodoService {
	return services.NewTodoService(postgres.NewPQTodoRepository(db))
}
