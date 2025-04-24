// internal/domain/entities/user_metrics.go
package entities

import "github.com/google/uuid"

type UserMetrics struct {
	UserID        uuid.UUID
	Recency       int
	Frequency     int
	Monetary      float64
	TBP           float64 // Time Between Purchases
	AvgCheck      float64
	LastSegmentID uuid.UUID
}
