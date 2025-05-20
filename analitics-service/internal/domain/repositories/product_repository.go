package repositories

import (
	"context"

	"analitics-service/internal/domain/enteties"
)

// ProductRepository определяет интерфейс для работы с продуктами
type ProductRepository interface {
	// GetAllProducts возвращает все продукты
	GetAllProducts(ctx context.Context) ([]entities.Product, error)

	// GetProductByID возвращает продукт по его ID
	GetProductByID(ctx context.Context, productID string) (entities.Product, error)

	// CreateProduct создает новый продукт
	CreateProduct(ctx context.Context, product entities.Product) error

	// UpdateProduct обновляет существующий продукт
	UpdateProduct(ctx context.Context, product entities.Product) error

	// DeleteProduct удаляет продукт по его ID
	DeleteProduct(ctx context.Context, productID string) error
}
