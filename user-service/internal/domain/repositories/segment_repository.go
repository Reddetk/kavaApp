// internal/domain/repositories/segment_repository.go
package repositories

import (
	"context"
	"user-service/internal/domain/entities"

	"github.com/google/uuid"
)

type SegmentRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*entities.Segment, error)
	Create(ctx context.Context, segment *entities.Segment) error
	Update(ctx context.Context, segment *entities.Segment) error
	GetByType(ctx context.Context, segmentType string) ([]*entities.Segment, error)
}
