// internal/domain/repositories/transaction_repository.go
package repositories

import (
	"context"
	"time"
	"user-service/internal/domain/entities"

	"github.com/google/uuid"
)

type TransactionRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Transaction, error)
	GetByPeriod(ctx context.Context, start, end time.Time) ([]*entities.Transaction, error)
	Create(ctx context.Context, transaction *entities.Transaction) error
	Ping(ctx context.Context) error
}
