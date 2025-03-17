package factories

import (
	"database/sql"

	"go-starter-template/pkg/application/interfaces"
	"go-starter-template/pkg/application/services"
	"go-starter-template/pkg/infrastructure/db/postgres"
)

func NewTodoServiceWithPQRepository(db *sql.DB) interfaces.TodoService {
	return services.NewTodoService(postgres.NewPQTodoRepository(db))
}
