// internal/domain/services/clv_service.go
package services

import "github.com/google/uuid"

type CLVService interface {
	CalculateCLV(userID uuid.UUID, retentionRate float64, avgCheck float64) (float64, error)
	UpdateCLV(userID uuid.UUID) (float64, error)
}
