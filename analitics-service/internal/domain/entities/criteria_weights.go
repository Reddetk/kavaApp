// internal/domain/entities/criteria_weights.go
package entities

import (
	"errors"
	"fmt"
	"math"
)

// CriteriaWeights содержит веса для разных критериев ABC-анализа
type CriteriaWeights struct {
	RevenueWeight  float64 `json:"revenue_weight"`
	QuantityWeight float64 `json:"quantity_weight"`
	ProfitWeight   float64 `json:"profit_weight"`
}

// Validate проверяет корректность данных в структуре CriteriaWeights
func (w *CriteriaWeights) Validate() error {
	if w.RevenueWeight < 0 {
		return fmt.Errorf("revenue weight cannot be negative, got %f", w.RevenueWeight)
	}

	if w.QuantityWeight < 0 {
		return fmt.Errorf("quantity weight cannot be negative, got %f", w.QuantityWeight)
	}

	if w.ProfitWeight < 0 {
		return fmt.Errorf("profit weight cannot be negative, got %f", w.ProfitWeight)
	}

	// Проверяем, что хотя бы один вес положительный
	if w.RevenueWeight == 0 && w.QuantityWeight == 0 && w.ProfitWeight == 0 {
		return errors.New("at least one weight must be positive")
	}

	// Проверяем, что сумма весов равна 1 (с небольшой погрешностью)
	sum := w.RevenueWeight + w.QuantityWeight + w.ProfitWeight
	if math.Abs(sum-1.0) > 0.001 {
		return fmt.Errorf("sum of weights must be 1.0, got %f", sum)
	}

	return nil
}
