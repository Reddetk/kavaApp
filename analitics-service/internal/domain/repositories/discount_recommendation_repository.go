package repositories

import (
	"context"

	"analitics-service/internal/domain/entities"
)

// DiscountRecommendationRepository определяет интерфейс для работы с рекомендациями по скидкам
type DiscountRecommendationRepository interface {
	// SaveRecommendation сохраняет рекомендацию по скидке
	SaveRecommendation(ctx context.Context, recommendation entities.DiscountRecommendation) error

	// GetRecommendationByProductID возвращает рекомендацию по скидке для конкретного продукта
	GetRecommendationByProductID(ctx context.Context, productID string) (entities.DiscountRecommendation, error)

	// GetRecommendationsByCategory возвращает рекомендации по скидкам для продуктов определенной категории
	GetRecommendationsByCategory(ctx context.Context, category string) ([]entities.DiscountRecommendation, error)

	// GetRecommendationsBySegment возвращает рекомендации по скидкам для продуктов определенного сегмента
	GetRecommendationsBySegment(ctx context.Context, segment entities.Segment) ([]entities.DiscountRecommendation, error)

	// GetLatestRecommendations возвращает последние рекомендации по скидкам
	GetLatestRecommendations(ctx context.Context, limit int) ([]entities.DiscountRecommendation, error)
}
