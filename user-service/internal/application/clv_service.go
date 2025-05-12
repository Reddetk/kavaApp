package application

import (
	"context"
	"fmt"
	"math"
	"time"

	"user-service/internal/domain/entities"
	"user-service/internal/domain/repositories"
	"user-service/internal/domain/services"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CLVService реализует логику расчета Customer Lifetime Value
type CLVService struct {
	userRepo      repositories.UserRepository
	metricsRepo   repositories.UserMetricsRepository
	segmentRepo   repositories.SegmentRepository
	clvCalculator services.CLVService
	retentionSvc  *RetentionService
	logger        *zap.Logger
}

// NewCLVService создает новый экземпляр сервиса CLV
func NewCLVService(
	ur repositories.UserRepository,
	mr repositories.UserMetricsRepository,
	cc services.CLVService,
	rs *RetentionService,
) *CLVService {
	logger, _ := zap.NewProduction()
	return &CLVService{
		userRepo:      ur,
		metricsRepo:   mr,
		clvCalculator: cc,
		retentionSvc:  rs,
		logger:        logger.With(zap.String("service", "CLVService")),
	}
}

// CalculateUserCLV рассчитывает и обновляет CLV для конкретного пользователя
func (s *CLVService) CalculateUserCLV(ctx context.Context, userID uuid.UUID) (float64, error) {
	// Получение метрик пользователя
	metrics, err := s.getOrCalculateMetrics(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user metrics: %w", err)
	}

	// Прогнозирование вероятности оттока
	churnProb, err := s.retentionSvc.PredictChurnProbability(ctx, userID)
	if err != nil {
		s.logger.Warn("failed to predict churn probability, using default",
			zap.String("userID", userID.String()),
			zap.Error(err))
		churnProb = 0.2 // Значение по умолчанию
	}

	// Расчет коэффициента удержания
	retentionRate := math.Max(0, 1-math.Min(churnProb, 1))

	// Вызов доменного сервиса для расчета CLV
	clv, err := s.clvCalculator.CalculateCLV(
		userID,
		retentionRate,
		metrics.AvgCheck,
	)
	if err != nil {
		return 0, fmt.Errorf("CLV calculation failed: %w", err)
	}

	// Обновление метрик пользователя
	if err := s.updateUserCLV(ctx, userID, clv); err != nil {
		s.logger.Error("failed to update user CLV",
			zap.String("userID", userID.String()),
			zap.Error(err))
	}

	return clv, nil
}

// BatchUpdateCLV выполняет пакетное обновление CLV для всех пользователей
func (s *CLVService) BatchUpdateCLV(ctx context.Context, batchSize int) error {
	offset := 0
	for {
		users, err := s.userRepo.List(ctx, batchSize, offset)
		if err != nil {
			return fmt.Errorf("failed to list users: %w", err)
		}

		if len(users) == 0 {
			break
		}

		for _, user := range users {
			if _, err := s.CalculateUserCLV(ctx, user.ID); err != nil {
				s.logger.Error("batch CLV update failed",
					zap.String("userID", user.ID.String()),
					zap.Error(err))
			}
		}

		offset += len(users)
		if len(users) < batchSize {
			break
		}
	}
	return nil
}

// EstimateCLV оценивает CLV для конкретного сценария
func (s *CLVService) EstimateCLV(ctx context.Context, userID uuid.UUID, scenario string) (float64, error) {
	// Получение метрик пользователя
	metrics, err := s.getOrCalculateMetrics(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user metrics: %w", err)
	}

	// Прогнозирование вероятности оттока
	churnProb, err := s.retentionSvc.PredictChurnProbability(ctx, userID)
	if err != nil {
		s.logger.Warn("failed to predict churn probability, using default",
			zap.String("userID", userID.String()),
			zap.Error(err))
		churnProb = 0.2 // Значение по умолчанию
	}

	// Расчет коэффициента удержания с учетом сценария
	var retentionRate float64
	var avgCheck float64

	switch scenario {
	case "optimistic":
		retentionRate = math.Max(0, 1-math.Min(churnProb*0.8, 1)) // Улучшенное удержание
		avgCheck = metrics.AvgCheck * 1.2                         // Увеличенный средний чек
	case "pessimistic":
		retentionRate = math.Max(0, 1-math.Min(churnProb*1.2, 1)) // Ухудшенное удержание
		avgCheck = metrics.AvgCheck * 0.8                         // Уменьшенный средний чек
	default: // "default" или любой другой сценарий
		retentionRate = math.Max(0, 1-math.Min(churnProb, 1))
		avgCheck = metrics.AvgCheck
	}

	// Вызов доменного сервиса для расчета CLV
	clv, err := s.clvCalculator.EstimateCLV(userID, scenario)
	if err != nil {
		// Если метод EstimateCLV не реализован, используем CalculateCLV
		clv, err = s.clvCalculator.CalculateCLV(userID, retentionRate, avgCheck)
		if err != nil {
			return 0, fmt.Errorf("CLV estimation failed: %w", err)
		}
	}

	return clv, nil
}

// GetHistoricalCLV возвращает исторические данные CLV для пользователя
func (s *CLVService) GetHistoricalCLV(ctx context.Context, userID uuid.UUID, periodMonths int) ([]entities.CLVDataPoint, error) {
	// Проверяем существование пользователя
	user, err := s.userRepo.Get(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found: %s", userID)
	}

	// Вызов доменного сервиса для получения исторических данных
	history, err := s.clvCalculator.GetHistoricalCLV(userID, periodMonths)
	if err != nil {
		return nil, fmt.Errorf("failed to get historical CLV: %w", err)
	}

	// Если история пуста, создаем одну точку с текущим значением CLV
	if len(history) == 0 {
		metrics, err := s.metricsRepo.Get(ctx, userID)
		if err != nil || metrics == nil {
			// Если метрики недоступны, рассчитываем CLV
			clv, err := s.CalculateUserCLV(ctx, userID)
			if err != nil {
				return nil, fmt.Errorf("failed to calculate current CLV: %w", err)
			}

			// Создаем одну историческую точку
			history = append(history, entities.CLVDataPoint{
				UserID:    userID,
				Date:      time.Now(),
				Value:     clv,
				Scenario:  "default",
				Algorithm: "DiscountedCashFlow",
			})
		} else {
			// Используем значение CLV из метрик
			history = append(history, entities.CLVDataPoint{
				UserID:    userID,
				Date:      metrics.LastCLVUpdate,
				Value:     metrics.CLV,
				Scenario:  "default",
				Algorithm: "DiscountedCashFlow",
			})
		}
	}

	return history, nil
}

// getOrCalculateMetrics получает или рассчитывает метрики пользователя
func (s *CLVService) getOrCalculateMetrics(ctx context.Context, userID uuid.UUID) (*entities.UserMetrics, error) {
	metrics, err := s.metricsRepo.Get(ctx, userID)
	if err == nil && metrics != nil && metrics.IsValid() {
		return metrics, nil
	}

	s.logger.Info("calculating metrics for user",
		zap.String("userID", userID.String()))

	metrics, err = s.metricsRepo.CalculateMetrics(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("metrics calculation failed: %w", err)
	}

	if metrics.CLV == 0 {
		// Если CLV не рассчитан, устанавливаем начальное значение
		metrics.CLV = metrics.AvgCheck * 12 // Простая оценка: средний чек * 12 месяцев
	}

	if err := s.metricsRepo.Update(ctx, metrics); err != nil {
		s.logger.Warn("failed to save calculated metrics",
			zap.String("userID", userID.String()),
			zap.Error(err))
	}

	return metrics, nil
}

// updateUserCLV обновляет значение CLV пользователя
func (s *CLVService) updateUserCLV(ctx context.Context, userID uuid.UUID, clv float64) error {
	metrics, err := s.metricsRepo.Get(ctx, userID)
	if err != nil {
		return err
	}

	metrics.CLV = clv
	metrics.LastCLVUpdate = time.Now()

	return s.metricsRepo.Update(ctx, metrics)
}
