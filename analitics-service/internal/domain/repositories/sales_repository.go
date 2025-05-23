package repositories

import (
	"context"
	"time"

	"analitics-service/internal/domain/entities"
)

// SalesRepository определяет интерфейс для работы с продажами
type SalesRepository interface {
	// GetSalesByPeriod возвращает продажи за указанный период
	GetSalesByPeriod(ctx context.Context, startDate, endDate time.Time) ([]entities.Sale, error)

	// GetSalesByProductID возвращает продажи для конкретного продукта
	GetSalesByProductID(ctx context.Context, productID string, startDate, endDate time.Time) ([]entities.Sale, error)

	// GetSalesByCustomerID возвращает продажи для конкретного клиента
	GetSalesByCustomerID(ctx context.Context, customerID string, startDate, endDate time.Time) ([]entities.Sale, error)

	// CreateSale создает новую запись о продаже
	CreateSale(ctx context.Context, sale entities.Sale) error

	// GetSaleByID возвращает продажу по её ID
	GetSaleByID(ctx context.Context, saleID string) (entities.Sale, error)

	// GetDailySalesData возвращает агрегированные данные о продажах по дням
	GetDailySalesData(ctx context.Context, startDate, endDate time.Time) ([]entities.DailyTransactionData, error)
}
