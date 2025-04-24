// internal/domain/services/segmentation_service.go
package services

import "user-service/internal/domain/entities"

type SegmentationService interface {
	PerformRFMClustering(users []entities.UserMetrics) ([]entities.Segment, error)
	PerformBehaviorClustering(transactions []entities.Transaction) ([]entities.Segment, error)
	AssignUserToSegment(userID uuid.UUID, metrics entities.UserMetrics) (entities.Segment, error)
}
