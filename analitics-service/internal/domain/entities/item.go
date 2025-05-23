// internal/domain/entities/item.go
package entities

import (
	"errors"
	"fmt"
)

// Item представляет собой элемент в транзакции
type Item struct {
	ProductID   string  `json:"product_id"`
	Name        string  `json:"name"`
	CategoryID  string  `json:"category_id"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	DiscountPct float64 `json:"discount_pct,omitempty"`
}

// Validate проверяет корректность данных в структуре Item
func (i *Item) Validate() error {
	if i.ProductID == "" {
		return errors.New("product ID is required")
	}

	if i.Name == "" {
		return errors.New("item name is required")
	}

	if i.Price < 0 {
		return fmt.Errorf("price cannot be negative, got %f", i.Price)
	}

	if i.Quantity <= 0 {
		return fmt.Errorf("quantity must be positive, got %d", i.Quantity)
	}

	if i.DiscountPct < 0 || i.DiscountPct > 100 {
		return fmt.Errorf("discount percentage must be between 0 and 100, got %f", i.DiscountPct)
	}

	return nil
}
