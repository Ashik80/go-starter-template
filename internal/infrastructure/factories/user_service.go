package factories

import (
	"database/sql"

	"go-starter-template/internal/application/interfaces"
	"go-starter-template/internal/application/services"
	"go-starter-template/internal/infrastructure/db/postgres"
	"go-starter-template/pkg/logger"
	"go-starter-template/pkg/security"
)

func NewUserServiceWithPQRepository(db *sql.DB, log *logger.Logger) interfaces.UserService {
	return services.NewUserService(
		postgres.NewPQUserRepository(db),
		security.NewBcryptPasswordHasher(log),
		services.NewSessionService(postgres.NewPQSessionRepository(db)),
	)
}
