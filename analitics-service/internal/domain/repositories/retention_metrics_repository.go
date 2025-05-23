package repositories

import (
	"context"
	"time"

	"analitics-service/internal/domain/entities"
)

// RetentionMetricsRepository определяет интерфейс для работы с метриками удержания клиентов
type RetentionMetricsRepository interface {
	// SaveMetrics сохраняет метрики удержания
	SaveMetrics(ctx context.Context, metrics entities.RetentionMetrics) error

	// GetMetricsByPeriod возвращает метрики удержания за указанный период
	GetMetricsByPeriod(ctx context.Context, period entities.TimeRange, date time.Time) (entities.RetentionMetrics, error)

	// GetMetricsHistory возвращает историю метрик удержания
	GetMetricsHistory(ctx context.Context, period entities.TimeRange, startDate, endDate time.Time) ([]entities.RetentionMetrics, error)

	// GetLatestMetrics возвращает последние метрики удержания
	GetLatestMetrics(ctx context.Context, period entities.TimeRange) (entities.RetentionMetrics, error)
}
