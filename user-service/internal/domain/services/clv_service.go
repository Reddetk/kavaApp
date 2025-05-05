// internal/domain/services/clv_service.go
package services

import (
	"user-service/internal/domain/entities"

	"github.com/google/uuid"
)

type CLVService interface {
	CalculateCLV(userID uuid.UUID, retentionRate float64, avgCheck float64) (float64, error)
	UpdateCLV(userID uuid.UUID) (float64, error)
	EstimateCLV(userID uuid.UUID, scenario string) (float64, error)
	GetHistoricalCLV(userID uuid.UUID, periodMonths int) ([]entities.CLVDataPoint, error)
}
