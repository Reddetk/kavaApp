// internal/domain/entities/analysis_metadata.go
package entities

import (
	"errors"
	"fmt"
	"time"
)

// AnalysisMetadata содержит метаданные для анализа
type AnalysisMetadata struct {
	AnalysisDate time.Time `json:"analysis_date"`
	PeriodStart  time.Time `json:"period_start"`
	PeriodEnd    time.Time `json:"period_end"`
}

// Validate проверяет корректность данных в структуре AnalysisMetadata
func (m *AnalysisMetadata) Validate() error {
	if m.AnalysisDate.IsZero() {
		return errors.New("analysis date is required")
	}

	if m.PeriodStart.IsZero() {
		return errors.New("period start date is required")
	}

	if m.PeriodEnd.IsZero() {
		return errors.New("period end date is required")
	}

	if m.PeriodStart.After(m.PeriodEnd) {
		return fmt.Errorf("period start date (%s) cannot be after period end date (%s)",
			m.PeriodStart.Format(time.RFC3339), m.PeriodEnd.Format(time.RFC3339))
	}

	return nil
}
