// internal/domain/entities/product.go
package entities

import (
	"errors"
)

// Product представляет товар в системе
type Product struct {
	BaseEntity
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	CategoryID  string  `json:"category_id"`
	SubCategory string  `json:"sub_category"`
	Price       float64 `json:"price"`
	Cost        float64 `json:"cost"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	IsActive    bool    `json:"is_active"`
}

// Validate проверяет корректность данных в структуре Product
func (p *Product) Validate() error {
	if p.Name == "" {
		return errors.New("product name is required")
	}

	if p.Price < 0 {
		return errors.New("product price cannot be negative")
	}

	if p.Cost < 0 {
		return errors.New("product cost cannot be negative")
	}

	if p.CategoryID == "" {
		return errors.New("category ID is required")
	}

	return nil
}
