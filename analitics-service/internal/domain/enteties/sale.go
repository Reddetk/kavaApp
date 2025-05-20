// internal/domain/entities/sale.go
package entities

import (
	"errors"
	"fmt"
	"time"
)

// Sale представляет отдельную продажу товара
type Sale struct {
	BaseEntity
	ProductID     string    `json:"product_id"`
	Quantity      int       `json:"quantity"`
	Price         float64   `json:"price"`
	DiscountRate  float64   `json:"discount_rate"`
	PurchaseDate  time.Time `json:"purchase_date"`
	CustomerID    string    `json:"customer_id"`
	TransactionID string    `json:"transaction_id"`
}

// Validate проверяет корректность данных в структуре Sale
func (s *Sale) Validate() error {
	if s.ProductID == "" {
		return errors.New("product ID is required")
	}

	if s.Quantity <= 0 {
		return fmt.Errorf("quantity must be positive, got %d", s.Quantity)
	}

	if s.Price < 0 {
		return fmt.Errorf("price cannot be negative, got %f", s.Price)
	}

	if s.DiscountRate < 0 || s.DiscountRate > 100 {
		return fmt.Errorf("discount rate must be between 0 and 100, got %f", s.DiscountRate)
	}

	if s.PurchaseDate.IsZero() {
		return errors.New("purchase date is required")
	}

	if s.CustomerID == "" {
		return errors.New("customer ID is required")
	}

	if s.TransactionID == "" {
		return errors.New("transaction ID is required")
	}

	return nil
}
