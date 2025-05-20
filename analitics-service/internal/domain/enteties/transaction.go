// internal/domain/entities/transaction.go
package entities

import (
	"errors"
	"fmt"
	"time"
)

// Transaction представляет собой транзакцию покупки
type Transaction struct {
	BaseEntity
	CustomerID   string    `json:"customer_id"`
	Date         time.Time `json:"date"`
	TotalAmount  float64   `json:"total_amount"`
	Items        []Item    `json:"items"`
	DiscountUsed bool      `json:"discount_used"`
	CouponCode   string    `json:"coupon_code,omitempty"`
}

// Validate проверяет корректность данных в структуре Transaction
func (t *Transaction) Validate() error {
	if t.CustomerID == "" {
		return errors.New("customer ID is required")
	}

	if t.Date.IsZero() {
		return errors.New("transaction date is required")
	}

	if t.TotalAmount < 0 {
		return errors.New("total amount cannot be negative")
	}

	if len(t.Items) == 0 {
		return errors.New("transaction must have at least one item")
	}

	// Проверяем каждый элемент транзакции
	for i, item := range t.Items {
		if err := item.Validate(); err != nil {
			return errors.New("invalid item at index " + fmt.Sprint(i) + ": " + err.Error())
		}
	}

	// Если указано, что использован купон, но код купона пустой
	if t.DiscountUsed && t.CouponCode == "" {
		return errors.New("coupon code is required when discount is used")
	}

	return nil
}
