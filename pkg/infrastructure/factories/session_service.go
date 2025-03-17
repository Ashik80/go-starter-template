package factories

import (
	"database/sql"

	"go-starter-template/pkg/application/interfaces"
	"go-starter-template/pkg/application/services"
	"go-starter-template/pkg/infrastructure/db/postgres"
)

func NewSessionServiceWithPQRepository(db *sql.DB) interfaces.SessionService {
	return services.NewSessionService(postgres.NewPQSessionRepository(db))
}
