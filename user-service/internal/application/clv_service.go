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
	sr repositories.SegmentRepository,
	cs services.CLVService,
	rs *RetentionService,
	logger *zap.Logger,
) *CLVService {
	return &CLVService{
		userRepo:      ur,
		metricsRepo:   mr,
		segmentRepo:   sr,
		clvCalculator: cs,
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

	// Получение информации о сегменте пользователя
	segment, err := s.segmentRepo.Get(ctx, metrics.LastSegmentID)
	if err != nil {
		s.logger.Warn("failed to get user segment, using default values",
			zap.String("userID", userID.String()),
			zap.Error(err))
	}

	// Корректировка среднего чека в зависимости от сегмента
	adjustedAvgCheck := s.adjustAvgCheckBySegment(metrics.AvgCheck, segment)

	// Вызов доменного сервиса для расчета CLV
	clv, err := s.clvCalculator.CalculateCLV(
		userID,
		retentionRate,
		adjustedAvgCheck,
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

	if err := s.metricsRepo.Create(ctx, metrics); err != nil {
		s.logger.Warn("failed to save calculated metrics",
			zap.String("userID", userID.String()),
			zap.Error(err))
	}

	return metrics, nil
}

// adjustAvgCheckBySegment корректирует средний чек по сегменту
func (s *CLVService) adjustAvgCheckBySegment(avgCheck float64, segment *entities.Segment) float64 {
	if segment == nil {
		return avgCheck
	}

	// Логика корректировки на основе данных сегмента
	if multiplier, ok := segment.CentroidData["avg_check_multiplier"].(float64); ok {
		return avgCheck * multiplier
	}
	return avgCheck
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
