// internal/domain/repositories/user_repository.go
package repositories

import (
	"context"
	"user-service/internal/domain/entities"

	"github.com/google/uuid"
)

type UserRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*entities.User, error)
	Create(ctx context.Context, user *entities.User) error
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entities.User, error)
	Ping(ctx context.Context) error
}
