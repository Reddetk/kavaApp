package repositories

import (
	"context"
)

// ProfitMarginRepository определяет интерфейс для работы с данными о прибыльности продуктов
type ProfitMarginRepository interface {
	// GetProfitMargins возвращает маржу прибыли для всех продуктов
	GetProfitMargins(ctx context.Context) (map[string]float64, error)
	
	// GetProfitMarginByProductID возвращает маржу прибыли для конкретного продукта
	GetProfitMarginByProductID(ctx context.Context, productID string) (float64, error)
	
	// UpdateProfitMargin обновляет маржу прибыли для продукта
	UpdateProfitMargin(ctx context.Context, productID string, margin float64) error
	
	// UpdateProfitMargins обновляет маржу прибыли для нескольких продуктов
	UpdateProfitMargins(ctx context.Context, margins map[string]float64) error
}