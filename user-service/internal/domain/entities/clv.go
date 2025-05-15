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

// CLV представляет пожизненную ценность клиента
type CLV struct {
	UserID       uuid.UUID `json:"user_id"`
	Value        float64   `json:"value"`
	Currency     string    `json:"currency"`
	CalculatedAt time.Time `json:"calculated_at"`
	Forecast     float64   `json:"forecast"`
	Confidence   float64   `json:"confidence"`
	Scenario     string    `json:"scenario"`
}
