// internal/domain/entities/thresholds.go
package entities

import (
	"fmt"
)

// Thresholds содержит пороговые значения для определения сегментов A, B, C
type Thresholds struct {
	AThreshold float64 `json:"a_threshold"`
	BThreshold float64 `json:"b_threshold"`
}

// Validate проверяет корректность данных в структуре Thresholds
func (t *Thresholds) Validate() error {
	if t.AThreshold <= 0 || t.AThreshold >= 100 {
		return fmt.Errorf("a threshold must be between 0 and 100, got %f", t.AThreshold)
	}

	if t.BThreshold <= t.AThreshold || t.BThreshold >= 100 {
		return fmt.Errorf("b threshold must be between A threshold (%f) and 100, got %f",
			t.AThreshold, t.BThreshold)
	}

	return nil
}
