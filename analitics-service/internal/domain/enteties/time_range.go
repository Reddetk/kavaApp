// internal/domain/entities/time_range.go
package entities

// TimeRange представляет временные периоды для отчетов
type TimeRange string

const (
	Daily   TimeRange = "daily"
	Weekly  TimeRange = "weekly"
	Monthly TimeRange = "monthly"
)
