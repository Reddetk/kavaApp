// internal/domain/enteties/abc_analysis_criteria.go
package entities

import (
	"errors"
	"fmt"
	"time"
)

// ABCAnalysisCriteria содержит критерии для проведения ABC-анализа
type ABCAnalysisCriteria struct {
	StartDate          time.Time       `json:"start_date"`
	EndDate            time.Time       `json:"end_date"`
	ThresholdsRevenue  Thresholds      `json:"thresholds_revenue"`
	ThresholdsQuantity Thresholds      `json:"thresholds_quantity"`
	ThresholdsProfit   Thresholds      `json:"thresholds_profit"`
	Weights            CriteriaWeights `json:"weights"`
}

// Validate проверяет корректность данных в структуре ABCAnalysisCriteria
func (c *ABCAnalysisCriteria) Validate() error {
	if c.StartDate.IsZero() {
		return errors.New("start date is required")
	}

	if c.EndDate.IsZero() {
		return errors.New("end date is required")
	}

	if c.StartDate.After(c.EndDate) {
		return fmt.Errorf("start date (%s) cannot be after end date (%s)",
			c.StartDate.Format(time.RFC3339), c.EndDate.Format(time.RFC3339))
	}

	if err := c.ThresholdsRevenue.Validate(); err != nil {
		return fmt.Errorf("invalid revenue thresholds: %w", err)
	}

	if err := c.ThresholdsQuantity.Validate(); err != nil {
		return fmt.Errorf("invalid quantity thresholds: %w", err)
	}

	if err := c.ThresholdsProfit.Validate(); err != nil {
		return fmt.Errorf("invalid profit thresholds: %w", err)
	}

	if err := c.Weights.Validate(); err != nil {
		return fmt.Errorf("invalid weights: %w", err)
	}

	return nil
}
