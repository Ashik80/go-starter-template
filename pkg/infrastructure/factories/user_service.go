package factories

import (
	"database/sql"

	"go-starter-template/pkg/application/interfaces"
	"go-starter-template/pkg/application/services"
	"go-starter-template/pkg/infrastructure/db/postgres"
	"go-starter-template/pkg/infrastructure/logger"
	"go-starter-template/pkg/infrastructure/security"
)

func NewUserServiceWithPQRepository(db *sql.DB, log *logger.Logger) interfaces.UserService {
	return services.NewUserService(
		postgres.NewPQUserRepository(db),
		security.NewBcryptPasswordHasher(log),
		services.NewSessionService(postgres.NewPQSessionRepository(db)),
	)
}
