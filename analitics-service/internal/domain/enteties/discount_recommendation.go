// internal/domain/entities/discount_recommendation.go
package entities

import (
	"errors"
	"fmt"
)

// DiscountRecommendation представляет рекомендацию по оптимальной скидке
type DiscountRecommendation struct {
	AnalysisMetadata
	ProductID        string  `json:"product_id,omitempty"`
	Category         string  `json:"category,omitempty"`
	OptimalDiscount  float64 `json:"optimal_discount"`
	LiftFactor       float64 `json:"lift_factor"`
	ABCCategory      Segment `json:"abc_category"`
	Confidence       float64 `json:"confidence"`
	AdjustmentReason string  `json:"adjustment_reason,omitempty"`
}

// Validate проверяет корректность данных в структуре DiscountRecommendation
func (dr *DiscountRecommendation) Validate() error {
	if dr.ProductID == "" && dr.Category == "" {
		return errors.New("either product ID or category must be specified")
	}

	if dr.OptimalDiscount < 0 || dr.OptimalDiscount > 100 {
		return fmt.Errorf("optimal discount must be between 0 and 100, got %f", dr.OptimalDiscount)
	}

	if dr.LiftFactor <= 0 {
		return fmt.Errorf("lift factor must be positive, got %f", dr.LiftFactor)
	}

	if !isValidSegment(dr.ABCCategory) {
		return fmt.Errorf("invalid ABC category: %s", dr.ABCCategory)
	}

	if dr.Confidence < 0 || dr.Confidence > 1 {
		return fmt.Errorf("confidence must be between 0 and 1, got %f", dr.Confidence)
	}

	return nil
}

// isValidSegment проверяет, является ли сегмент допустимым
func isValidSegment(s Segment) bool {
	return s == SegmentA || s == SegmentB || s == SegmentC
}
