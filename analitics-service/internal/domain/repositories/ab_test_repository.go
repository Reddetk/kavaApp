package repositories

import (
	"context"
	"time"

	"analitics-service/internal/domain/entities"
)

// ABTestRepository определяет интерфейс для работы с A/B тестами
type ABTestRepository interface {
	// GetTestResults возвращает результаты A/B тестов за указанный период
	GetTestResults(ctx context.Context, startDate, endDate time.Time) ([]entities.ABTestResult, error)

	// GetTestResultByID возвращает результат A/B теста по его ID
	GetTestResultByID(ctx context.Context, testID string) (entities.ABTestResult, error)

	// SaveTestResult сохраняет результат A/B теста
	SaveTestResult(ctx context.Context, result entities.ABTestResult) error

	// GetTestsByProduct возвращает A/B тесты для конкретного продукта
	GetTestsByProduct(ctx context.Context, productID string) ([]entities.ABTestResult, error)

	// GetTestsByCategory возвращает A/B тесты для продуктов определенной категории
	GetTestsByCategory(ctx context.Context, category string) ([]entities.ABTestResult, error)
}
