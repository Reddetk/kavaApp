// internal/domain/entities/clv.go
package entities

import (
	"time"

	"github.com/google/uuid"
)

// CLVDataPoint представляет историческую точку данных CLV
type CLVDataPoint struct {
	UserID    uuid.UUID
	Date      time.Time
	Value     float64
	Scenario  string // Сценарий расчета ("default", "optimistic", и т.д.)
	Algorithm string // Использованный алгоритм расчета
}
