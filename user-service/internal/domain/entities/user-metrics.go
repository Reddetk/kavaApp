// internal/domain/entities/user_metrics.go
package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserMetrics struct {
	UserID             uuid.UUID
	Recency            int
	Frequency          int
	Monetary           float64
	TBP                float64 // Time Between Purchases
	AvgCheck           float64
	LastSegmentID      uuid.UUID
	CLV                float64
	LastCLVUpdate      time.Time
	Age                int
	AvgSessionDuration float64
	SessionCount       float64
	Churned            bool
}

func (m *UserMetrics) IsValid() bool {
	return m.AvgCheck > 0 &&
		m.LastCLVUpdate.Add(7*24*time.Hour).After(time.Now())
}
