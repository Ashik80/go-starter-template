package factories

import (
	"database/sql"

	"go-starter-template/internal/application/interfaces"
	"go-starter-template/internal/application/services"
	"go-starter-template/internal/infrastructure/db/postgres"
)

func NewSessionServiceWithPQRepository(db *sql.DB) interfaces.SessionService {
	return services.NewSessionService(postgres.NewPQSessionRepository(db))
}
