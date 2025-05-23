package repositories

import (
	"context"
	"time"

	"analitics-service/internal/domain/entities"
)

// ABCSegmentRepository определяет интерфейс для работы с сегментацией ABC-анализа
type ABCSegmentRepository interface {
	// SaveSegmentation сохраняет результаты сегментации продуктов
	SaveSegmentation(ctx context.Context, segmentation map[string]entities.ProductFullSegmentation) error

	// GetProductSegmentation возвращает информацию о сегментации конкретного продукта
	GetProductSegmentation(ctx context.Context, productID string) (*entities.ProductSegmentation, error)

	// GetFullSegmentation возвращает полную информацию о сегментации всех продуктов
	GetFullSegmentation(ctx context.Context) (map[string]entities.ProductFullSegmentation, error)

	// GetSegmentationByCategory возвращает сегментацию продуктов по категории
	GetSegmentationByCategory(ctx context.Context, category string) ([]entities.ProductSegmentation, error)

	// GetLatestAnalysisDate возвращает дату последнего проведенного ABC-анализа
	GetLatestAnalysisDate(ctx context.Context) (time.Time, error)
}
