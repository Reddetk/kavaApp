// internal/domain/repositories/clv_repository.go
package repositories

import (
	"context"
	"user-service/internal/domain/entities"

	"github.com/google/uuid"
)

type CLVRepository interface {
	StoreDataPoint(ctx context.Context, point entities.CLVDataPoint) error
	GetHistory(ctx context.Context, userID uuid.UUID, periodMonths int) ([]entities.CLVDataPoint, error)
	PurgeOldData(ctx context.Context, olderThanDays int) error
}
