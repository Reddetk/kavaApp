// internal/domain/services/survival_analysis_service.go
package services

import "user-service/internal/domain/entities"

type SurvivalAnalysisService interface {
	BuildCoxModel(segment entities.Segment, users []entities.UserMetrics) error
	PredictChurnProbability(userID uuid.UUID, metrics entities.UserMetrics) (float64, error)
	PredictTimeToEvent(userID uuid.UUID, metrics entities.UserMetrics) (float64, error)
}
