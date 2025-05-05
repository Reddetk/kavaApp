package application

import (
	"context"
	"errors"
	"user-service/internal/domain/repositories"
	"user-service/internal/domain/services"

	"github.com/google/uuid"
)

type RetentionService struct {
	userRepo      repositories.UserRepository
	metricsRepo   repositories.UserMetricsRepository
	survivalSvc   services.SurvivalAnalysisService
	transitionSvc services.StateTransitionService
}

func NewRetentionService(
	ur repositories.UserRepository,
	mr repositories.UserMetricsRepository,
	ss services.SurvivalAnalysisService,
	ts services.StateTransitionService,
) *RetentionService {
	return &RetentionService{
		userRepo:      ur,
		metricsRepo:   mr,
		survivalSvc:   ss,
		transitionSvc: ts,
	}
}

// PredictChurnProbability predicts the probability that a user will churn.
func (s *RetentionService) PredictChurnProbability(ctx context.Context, userID uuid.UUID) (float64, error) {
	metrics, err := s.metricsRepo.Get(ctx, userID)
	if err != nil {
		return 0, err
	}
	if metrics == nil {
		return 0, errors.New("user metrics not found")
	}

	probability, err := s.survivalSvc.PredictChurnProbability(userID, *metrics)
	if err != nil {
		return 0, err
	}

	return probability, nil
}

// PredictTimeToEvent estimates expected time until the next event (e.g., churn or next purchase)
// using Cox proportional hazards model instead of a probability proxy.
func (s *RetentionService) PredictTimeToEvent(ctx context.Context, userID uuid.UUID) (float64, error) {
	metrics, err := s.metricsRepo.Get(ctx, userID)
	if err != nil {
		return 0, err
	}
	if metrics == nil {
		return 0, errors.New("user metrics not found")
	}

	// Используем survival analysis напрямую для оценки времени
	estimatedTime, err := s.survivalSvc.PredictTimeToEvent(userID, *metrics)
	if err != nil {
		return 0, err
	}

	return estimatedTime, nil
}
