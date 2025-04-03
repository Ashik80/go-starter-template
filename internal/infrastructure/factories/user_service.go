package factories

import (
	"database/sql"

	"go-starter-template/internal/application/services"
	"go-starter-template/internal/infrastructure/db/postgres"
	"go-starter-template/pkg/logger"
	"go-starter-template/pkg/security"
)

func NewUserServiceWithPQRepository(db *sql.DB, log *logger.Logger) services.IUserService {
	return services.NewUserService(
		postgres.NewPQUserRepository(db),
		security.NewBcryptPasswordHasher(log),
		services.NewSessionService(postgres.NewPQSessionRepository(db)),
	)
}
