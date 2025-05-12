// internal/domain/services/segmentation_service.go
package services

import (
	"user-service/internal/domain/entities"

	"github.com/google/uuid"
)

type SegmentationService interface {
	PerformRFMClustering(users []entities.UserMetrics) ([]entities.Segment, error)
	PerformBehaviorClustering(transactions []entities.Transaction) ([]entities.Segment, error)
	AssignUserToSegment(userID uuid.UUID, metrics entities.UserMetrics, segments []entities.Segment) (entities.Segment, error)
}