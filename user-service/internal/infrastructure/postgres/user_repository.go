// internal/infrastructure/postgres/user_repository.go
package postgres

import (
	"context"
	"database/sql"
	"user-service/internal/domain/entities"
	"user-service/internal/domain/repositories"

	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repositories.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Get(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	// Implementation
}

// Аналогично для других репозиториев...
