package services

import (
	"errors"
	"user-service/internal/domain/entities"
	"user-service/internal/domain/services"

	"github.com/google/uuid"
)

// Реализация сервиса анализа выживания с моделью Кокса (заглушка)
type CoxSurvivalAnalysis struct {
	// survivalCurves хранит псевдокривые выживания: userID → []S(t)
	survivalCurves map[uuid.UUID][]float64
}

// Конструктор
func NewCoxSurvivalAnalysis() services.SurvivalAnalysisService {
	return &CoxSurvivalAnalysis{
		survivalCurves: make(map[uuid.UUID][]float64),
	}
}

// Построение модели (можно заменить вызовом внешнего пайплайна или ML-кода)
func (s *CoxSurvivalAnalysis) BuildCoxModel(segment entities.Segment, users []entities.UserMetrics) error {
	for _, user := range users {
		// Пример псевдокривой выживания
		survival := []float64{1.0, 0.95, 0.90, 0.80, 0.65, 0.50, 0.30, 0.15, 0.05}
		s.survivalCurves[user.UserID] = survival
	}
	return nil
}

// Вероятность оттока (churn) на предпоследнем шаге кривой: 1 - S(t)
func (s *CoxSurvivalAnalysis) PredictChurnProbability(userID uuid.UUID, metrics entities.UserMetrics) (float64, error) {
	survival, ok := s.survivalCurves[userID]
	if !ok {
		return 0, errors.New("no survival curve for user: " + userID.String())
	}
	last := survival[len(survival)-1]
	return 1 - last, nil
}

// Ожидаемое время до события: E[T] ≈ sum(S(t))
func (s *CoxSurvivalAnalysis) PredictTimeToEvent(userID uuid.UUID, metrics entities.UserMetrics) (float64, error) {
	survival, ok := s.survivalCurves[userID]
	if !ok {
		return 0, errors.New("no survival curve for user: " + userID.String())
	}
	var expected float64
	for _, s_t := range survival {
		expected += s_t
	}
	return expected, nil
}
