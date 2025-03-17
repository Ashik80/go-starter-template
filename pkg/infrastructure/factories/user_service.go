package factories

import (
	"database/sql"

	"go-starter-template/pkg/application/interfaces"
	"go-starter-template/pkg/application/services"
	"go-starter-template/pkg/infrastructure/db/postgres"
	"go-starter-template/pkg/infrastructure/security"
)

func NewUserServiceWithPQRepository(db *sql.DB) interfaces.UserService {
	return services.NewUserService(
		postgres.NewPQUserRepository(db),
		security.NewBcryptPasswordHasher(),
		services.NewSessionService(postgres.NewPQSessionRepository(db)),
	)
}
