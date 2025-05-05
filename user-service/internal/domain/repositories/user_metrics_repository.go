// internal/domain/repositories/user_metrics_repository.go
package repositories

import (
	"context"
	"user-service/internal/domain/entities"

	"github.com/google/uuid"
)

type UserMetricsRepository interface {
	Get(ctx context.Context, userID uuid.UUID) (*entities.UserMetrics, error)
	Create(ctx context.Context, metrics *entities.UserMetrics) error
	Update(ctx context.Context, metrics *entities.UserMetrics) error
	CalculateMetrics(ctx context.Context, userID uuid.UUID) (*entities.UserMetrics, error)
}
