package repositories

import (
	"context"

	"go-starter-template/internal/domain/entities"
)

type ISessionRepository interface {
	Create(ctx context.Context, session *entities.Session) (*entities.Session, error)
	Get(ctx context.Context, sessionId string) (*entities.Session, error)
	GetWithUser(ctx context.Context, sessionId string) (*entities.Session, error)
	Delete(ctx context.Context, session *entities.Session) error
}
