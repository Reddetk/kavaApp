// internal/domain/services/state_transition_service.go
package services

import (
	"user-service/internal/domain/entities"

	"github.com/google/uuid"
)

type StateTransitionService interface {
	BuildTransitionMatrix(segment entities.Segment, transactions []entities.Transaction) (map[string]map[string]float64, error)
	PredictNextState(userID uuid.UUID, currentState string) (string, float64, error)
	UpdateUserState(userID uuid.UUID, churnProbability float64) error
}