package application

import (
	"context"
	"errors"
	"log"
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

// RecalculateUserMetrics пересчитывает метрики пользователя на основе его транзакций
func (s *RetentionService) RecalculateUserMetrics(ctx context.Context, userID uuid.UUID) error {
	// Используем метод CalculateMetrics из репозитория метрик
	_, err := s.metricsRepo.CalculateMetrics(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateChurnProbability обновляет вероятность оттока пользователя
func (s *RetentionService) UpdateChurnProbability(ctx context.Context, userID uuid.UUID) error {
	// Получаем метрики пользователя
	metrics, err := s.metricsRepo.Get(ctx, userID)
	if err != nil {
		return err
	}
	if metrics == nil {
		return errors.New("user metrics not found")
	}

	// Рассчитываем вероятность оттока
	probability, err := s.survivalSvc.PredictChurnProbability(userID, *metrics)
	if err != nil {
		return err
	}

	// Обновляем состояние пользователя в модели Маркова
	if err := s.transitionSvc.UpdateUserState(userID, probability); err != nil {
		log.Printf("Warning: failed to update user state: %v", err)
		// Не возвращаем ошибку, так как это некритичная операция
	}

	return nil
}